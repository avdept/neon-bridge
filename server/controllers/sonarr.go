package controllers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"dashboard-server/models"

	"github.com/gin-gonic/gin"
)

type SonarrTestConfig struct {
	ServerURL string `json:"serverUrl" binding:"required"`
	ApiKey    string `json:"apiKey" binding:"required"`
}

type SonarrSystemStatus struct {
	Version string `json:"version"`
}

type SonarrSeries struct {
	Statistics SonarrSeriesStatistics `json:"statistics"`
}

type SonarrSeriesStatistics struct {
	EpisodeFileCount  int `json:"episodeFileCount"`
	TotalEpisodeCount int `json:"totalEpisodeCount"`
}

type SonarrQueue struct {
	TotalRecords int                 `json:"totalRecords"`
	Records      []SonarrQueueRecord `json:"records"`
}

type SonarrQueueRecord struct {
	Size     int64 `json:"size"`
	Sizeleft int64 `json:"sizeleft"`
}

type SonarrDiskSpace struct {
	Path       string `json:"path"`
	Label      string `json:"label"`
	FreeSpace  int64  `json:"freeSpace"`
	TotalSpace int64  `json:"totalSpace"`
}

type SonarrHealthCheck struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl"`
}

type SonarrStats struct {
	TotalSeries      int                 `json:"totalSeries"`
	TotalEpisodes    int                 `json:"totalEpisodes"`
	MissingEpisodes  int                 `json:"missingEpisodes"`
	QueuedItems      int                 `json:"queuedItems"`
	DownloadProgress float64             `json:"downloadProgress"`
	TotalStorage     int64               `json:"totalStorage"`
	FreeStorage      int64               `json:"freeStorage"`
	Version          string              `json:"version"`
	HealthAlerts     []SonarrHealthCheck `json:"healthAlerts"`
}

func TestSonarrConnection(c *gin.Context) {
	var config SonarrTestConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration: " + err.Error()})
		return
	}

	if config.ServerURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serverUrl is required"})
		return
	}

	if config.ApiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "apiKey is required"})
		return
	}

	stats, err := fetchSonarrStats(config.ServerURL, config.ApiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func ProxySonarrStats(c *gin.Context) {
	widgetIDStr := c.Param("widget_id")
	widgetID, err := strconv.ParseUint(widgetIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid widget ID"})
		return
	}

	var widget models.Widget
	if err := DB.First(&widget, uint(widgetID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
		return
	}

	if widget.Type != "sonarr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not a Sonarr widget"})
		return
	}

	config := widget.Config
	serverURL, ok := config["serverUrl"].(string)
	if !ok || serverURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serverUrl not found in widget configuration"})
		return
	}

	apiKey, ok := config["apiKey"].(string)
	if !ok || apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "apiKey not found in widget configuration"})
		return
	}

	stats, err := fetchSonarrStats(serverURL, apiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func fetchSonarrStats(serverURL, apiKey string) (*SonarrStats, error) {
	// TODO: Refactor this to use waitgroups for concurrent requests it should increase performance
	if serverURL[len(serverURL)-1] == '/' {
		serverURL = serverURL[:len(serverURL)-1]
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	stats := &SonarrStats{}

	if systemStatus, err := fetchSonarrSystemStatus(client, serverURL, apiKey); err == nil {
		stats.Version = systemStatus.Version
	}

	series, err := fetchSonarrSeries(client, serverURL, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series data: %v", err)
	}

	stats.TotalSeries = len(series)
	totalEpisodes := 0
	missingEpisodes := 0

	for _, s := range series {
		totalEpisodes += s.Statistics.TotalEpisodeCount
		missing := s.Statistics.TotalEpisodeCount - s.Statistics.EpisodeFileCount
		if missing > 0 {
			missingEpisodes += missing
		}
	}

	stats.TotalEpisodes = totalEpisodes
	stats.MissingEpisodes = missingEpisodes

	if queue, err := fetchSonarrQueue(client, serverURL, apiKey); err == nil {
		stats.QueuedItems = queue.TotalRecords

		if len(queue.Records) > 0 {
			var totalSize, totalCompleted int64
			for _, record := range queue.Records {
				if record.Size > 0 {
					totalSize += record.Size
					totalCompleted += (record.Size - record.Sizeleft)
				}
			}
			if totalSize > 0 {
				stats.DownloadProgress = float64(totalCompleted) / float64(totalSize) * 100
			}
		}
	}

	if diskSpaces, err := fetchSonarrDiskSpace(client, serverURL, apiKey); err == nil && len(diskSpaces) > 0 {
		// TODO: Temporary set to 4th disk which for me is /tv folder. Need to test it with other folks setups
		stats.TotalStorage = diskSpaces[3].TotalSpace
		stats.FreeStorage = diskSpaces[3].FreeSpace
	}

	if healthChecks, err := fetchSonarrHealth(client, serverURL, apiKey); err == nil {
		stats.HealthAlerts = healthChecks
	} else {
		stats.HealthAlerts = []SonarrHealthCheck{}
	}

	return stats, nil
}

func fetchSonarrSystemStatus(client *http.Client, serverURL, apiKey string) (*SonarrSystemStatus, error) {
	url := fmt.Sprintf("%s/api/v3/system/status?apikey=%s", serverURL, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Homepage-Dashboard/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var status SonarrSystemStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func fetchSonarrSeries(client *http.Client, serverURL, apiKey string) ([]SonarrSeries, error) {
	url := fmt.Sprintf("%s/api/v3/series?apikey=%s", serverURL, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Homepage-Dashboard/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var series []SonarrSeries
	if err := json.Unmarshal(body, &series); err != nil {
		return nil, err
	}

	return series, nil
}

func fetchSonarrQueue(client *http.Client, serverURL, apiKey string) (*SonarrQueue, error) {
	url := fmt.Sprintf("%s/api/v3/queue?apikey=%s", serverURL, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Homepage-Dashboard/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var queue SonarrQueue
	if err := json.Unmarshal(body, &queue); err != nil {
		return nil, err
	}

	return &queue, nil
}

func fetchSonarrDiskSpace(client *http.Client, serverURL, apiKey string) ([]SonarrDiskSpace, error) {
	url := fmt.Sprintf("%s/api/v3/diskspace?apikey=%s", serverURL, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Homepage-Dashboard/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var diskSpaces []SonarrDiskSpace
	if err := json.Unmarshal(body, &diskSpaces); err != nil {
		return nil, err
	}

	return diskSpaces, nil
}

func fetchSonarrHealth(client *http.Client, serverURL, apiKey string) ([]SonarrHealthCheck, error) {
	url := fmt.Sprintf("%s/api/v3/health?apikey=%s", serverURL, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Homepage-Dashboard/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var healthChecks []SonarrHealthCheck
	if err := json.Unmarshal(body, &healthChecks); err != nil {
		return nil, err
	}

	return healthChecks, nil
}
