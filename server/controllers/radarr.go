package controllers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"dashboard-server/database"
	"dashboard-server/models"

	"github.com/gin-gonic/gin"
)

type RadarrTestConfig struct {
	ServerURL string `json:"serverUrl" binding:"required"`
	ApiKey    string `json:"apiKey" binding:"required"`
}

type RadarrSystemStatus struct {
	Version string `json:"version"`
}

type RadarrMovie struct {
	HasFile    bool `json:"hasFile"`
	Downloaded bool `json:"downloaded"`
	Monitored  bool `json:"monitored"`
}

type RadarrQueue struct {
	TotalRecords int                 `json:"totalRecords"`
	Records      []RadarrQueueRecord `json:"records"`
}

type RadarrQueueRecord struct {
	Size     int64 `json:"size"`
	Sizeleft int64 `json:"sizeleft"`
}

type RadarrDiskSpace struct {
	Path       string `json:"path"`
	Label      string `json:"label"`
	FreeSpace  int64  `json:"freeSpace"`
	TotalSpace int64  `json:"totalSpace"`
}

type RadarrHealthCheck struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl"`
}

type RadarrStats struct {
	TotalMovies      int                 `json:"totalMovies"`
	DownloadedMovies int                 `json:"downloadedMovies"`
	MissingMovies    int                 `json:"missingMovies"`
	QueuedItems      int                 `json:"queuedItems"`
	DownloadProgress float64             `json:"downloadProgress"`
	TotalStorage     int64               `json:"totalStorage"`
	FreeStorage      int64               `json:"freeStorage"`
	Version          string              `json:"version"`
	HealthAlerts     []RadarrHealthCheck `json:"healthAlerts"`
}

func TestRadarrConnection(c *gin.Context) {
	var config RadarrTestConfig
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

	stats, err := fetchRadarrStats(config.ServerURL, config.ApiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func ProxyRadarrStats(c *gin.Context) {
	widgetIDStr := c.Param("widget_id")
	widgetID, err := strconv.ParseUint(widgetIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid widget ID"})
		return
	}

	widget := &models.Widget{}
	if err := database.DB.First(widget, uint(widgetID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
		return
	}

	if widget.Type != "radarr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not a Radarr widget"})
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

	stats, err := fetchRadarrStats(serverURL, apiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func fetchRadarrStats(serverURL, apiKey string) (*RadarrStats, error) {
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

	stats := &RadarrStats{}

	if systemStatus, err := fetchRadarrSystemStatus(client, serverURL, apiKey); err == nil {
		stats.Version = systemStatus.Version
	}

	movies, err := fetchRadarrMovies(client, serverURL, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies data: %v", err)
	}

	stats.TotalMovies = len(movies)
	downloadedMovies := 0
	missingMovies := 0

	for _, movie := range movies {
		if movie.HasFile || movie.Downloaded {
			downloadedMovies++
		} else if movie.Monitored {
			missingMovies++
		}
	}

	stats.DownloadedMovies = downloadedMovies
	stats.MissingMovies = missingMovies

	if queue, err := fetchRadarrQueue(client, serverURL, apiKey); err == nil {
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

	if diskSpaces, err := fetchRadarrDiskSpace(client, serverURL, apiKey); err == nil && len(diskSpaces) > 0 {
		if len(diskSpaces) > 0 {
			diskIndex := 0
			if len(diskSpaces) > 1 {
				for i, disk := range diskSpaces {
					if disk.Path == "/movies" || disk.Label == "movies" {
						diskIndex = i
						break
					}
				}
			}
			stats.TotalStorage = diskSpaces[diskIndex].TotalSpace
			stats.FreeStorage = diskSpaces[diskIndex].FreeSpace
		}
	}

	if healthChecks, err := fetchRadarrHealth(client, serverURL, apiKey); err == nil {
		stats.HealthAlerts = healthChecks
	} else {
		stats.HealthAlerts = []RadarrHealthCheck{}
	}

	return stats, nil
}

func fetchRadarrSystemStatus(client *http.Client, serverURL, apiKey string) (*RadarrSystemStatus, error) {
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

	var status RadarrSystemStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func fetchRadarrMovies(client *http.Client, serverURL, apiKey string) ([]RadarrMovie, error) {
	url := fmt.Sprintf("%s/api/v3/movie?apikey=%s", serverURL, apiKey)
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

	var movies []RadarrMovie
	if err := json.Unmarshal(body, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func fetchRadarrQueue(client *http.Client, serverURL, apiKey string) (*RadarrQueue, error) {
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

	var queue RadarrQueue
	if err := json.Unmarshal(body, &queue); err != nil {
		return nil, err
	}

	return &queue, nil
}

func fetchRadarrDiskSpace(client *http.Client, serverURL, apiKey string) ([]RadarrDiskSpace, error) {
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

	var diskSpaces []RadarrDiskSpace
	if err := json.Unmarshal(body, &diskSpaces); err != nil {
		return nil, err
	}

	return diskSpaces, nil
}

func fetchRadarrHealth(client *http.Client, serverURL, apiKey string) ([]RadarrHealthCheck, error) {
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

	var healthChecks []RadarrHealthCheck
	if err := json.Unmarshal(body, &healthChecks); err != nil {
		return nil, err
	}

	return healthChecks, nil
}
