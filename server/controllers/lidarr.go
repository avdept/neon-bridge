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

type LidarrTestConfig struct {
	ServerURL string `json:"serverUrl" binding:"required"`
	ApiKey    string `json:"apiKey" binding:"required"`
}

type LidarrSystemStatus struct {
	Version string `json:"version"`
}

type LidarrArtist struct {
	Monitored  bool                   `json:"monitored"`
	Statistics LidarrArtistStatistics `json:"statistics"`
}

type LidarrArtistStatistics struct {
	AlbumCount      int `json:"albumCount"`
	TrackFileCount  int `json:"trackFileCount"`
	TotalTrackCount int `json:"totalTrackCount"`
}

type LidarrQueue struct {
	TotalRecords int               `json:"totalRecords"`
	Records      []LidarrQueueItem `json:"records"`
}

type LidarrQueueItem struct {
	Size     float64 `json:"size"`
	SizeLeft float64 `json:"sizeleft"`
}

type LidarrDiskSpace struct {
	Path       string `json:"path"`
	Label      string `json:"label"`
	FreeSpace  int64  `json:"freeSpace"`
	TotalSpace int64  `json:"totalSpace"`
}

type LidarrHealthCheck struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl"`
}

type LidarrStats struct {
	QueuedItems      int                 `json:"queuedItems"`
	DownloadProgress float64             `json:"downloadProgress"`
	MissingAlbums    int                 `json:"missingAlbums"`
	MonitoredArtists int                 `json:"monitoredArtists"`
	TotalAlbums      int                 `json:"totalAlbums"`
	TotalTracks      int                 `json:"totalTracks"`
	TracksWithFiles  int                 `json:"tracksWithFiles"`
	FreeStorage      int64               `json:"freeStorage"`
	TotalStorage     int64               `json:"totalStorage"`
	HealthAlerts     []LidarrHealthCheck `json:"healthAlerts"`
}

func TestLidarrConnection(c *gin.Context) {
	var config LidarrTestConfig

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration: " + err.Error()})
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", config.ServerURL+"/api/v1/system/status", nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server URL"})
		return
	}

	req.Header.Set("X-Api-Key", config.ApiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Cannot connect to Lidarr server"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("Lidarr returned status %d", resp.StatusCode)})
		return
	}

	var systemStatus LidarrSystemStatus
	if err := json.NewDecoder(resp.Body).Decode(&systemStatus); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Invalid response from Lidarr"})
		return
	}

	stats, err := fetchLidarrStats(client, config.ServerURL, config.ApiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch Lidarr statistics: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func ProxyLidarrStats(c *gin.Context) {
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

	if widget.Type != "lidarr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not a Lidarr widget"})
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

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	stats, err := fetchLidarrStats(client, serverURL, apiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func fetchLidarrStats(client *http.Client, serverURL, apiKey string) (*LidarrStats, error) {
	// TODO: Refactor this to use waitgroups for concurrent requests it should increase performance
	stats := &LidarrStats{}

	queueReq, err := http.NewRequest("GET", serverURL+"/api/v1/queue?pageSize=100", nil)
	if err != nil {
		return stats, fmt.Errorf("failed to create queue request: %v", err)
	}
	queueReq.Header.Set("X-Api-Key", apiKey)
	queueReq.Header.Set("Accept", "application/json")

	queueResp, err := client.Do(queueReq)
	if err != nil {
		return stats, fmt.Errorf("failed to fetch queue: %v", err)
	}
	defer queueResp.Body.Close()

	if queueResp.StatusCode == 200 {
		var queue LidarrQueue
		queueBody, _ := io.ReadAll(queueResp.Body)
		if err := json.Unmarshal(queueBody, &queue); err == nil {
			stats.QueuedItems = queue.TotalRecords

			var totalSize, completedSize float64
			for _, item := range queue.Records {
				totalSize += item.Size
				completedSize += (item.Size - item.SizeLeft)
			}
			if totalSize > 0 {
				stats.DownloadProgress = (completedSize / totalSize) * 100
			}
		}
	}

	artistReq, err := http.NewRequest("GET", serverURL+"/api/v1/artist", nil)
	if err != nil {
		return stats, fmt.Errorf("failed to create artist request: %v", err)
	}
	artistReq.Header.Set("X-Api-Key", apiKey)
	artistReq.Header.Set("Accept", "application/json")

	artistResp, err := client.Do(artistReq)
	if err != nil {
		return stats, fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer artistResp.Body.Close()

	if artistResp.StatusCode == 200 {
		var artists []LidarrArtist
		artistBody, _ := io.ReadAll(artistResp.Body)
		if err := json.Unmarshal(artistBody, &artists); err == nil {
			for _, artist := range artists {
				if artist.Monitored {
					stats.MonitoredArtists++
				}
				stats.TotalAlbums += artist.Statistics.AlbumCount
				stats.TotalTracks += artist.Statistics.TotalTrackCount
				stats.TracksWithFiles += artist.Statistics.TrackFileCount
			}
		}
	}

	stats.MissingAlbums = 0

	diskReq, err := http.NewRequest("GET", serverURL+"/api/v1/diskspace", nil)
	if err != nil {
		return stats, fmt.Errorf("failed to create diskspace request: %v", err)
	}
	diskReq.Header.Set("X-Api-Key", apiKey)
	diskReq.Header.Set("Accept", "application/json")

	diskResp, err := client.Do(diskReq)
	if err != nil {
		return stats, fmt.Errorf("failed to fetch disk space: %v", err)
	}
	defer diskResp.Body.Close()

	if diskResp.StatusCode == 200 {
		var diskSpaces []LidarrDiskSpace
		diskBody, _ := io.ReadAll(diskResp.Body)
		if err := json.Unmarshal(diskBody, &diskSpaces); err == nil && len(diskSpaces) > 0 {
			stats.FreeStorage = diskSpaces[3].FreeSpace
			stats.TotalStorage = diskSpaces[3].TotalSpace
		}
	}

	if healthChecks, err := fetchLidarrHealth(client, serverURL, apiKey); err == nil {
		stats.HealthAlerts = healthChecks
	} else {
		stats.HealthAlerts = []LidarrHealthCheck{}
	}

	return stats, nil
}

func fetchLidarrHealth(client *http.Client, serverURL, apiKey string) ([]LidarrHealthCheck, error) {
	url := fmt.Sprintf("%s/api/v1/health?apikey=%s", serverURL, apiKey)
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

	var healthChecks []LidarrHealthCheck
	if err := json.Unmarshal(body, &healthChecks); err != nil {
		return nil, err
	}

	return healthChecks, nil
}
