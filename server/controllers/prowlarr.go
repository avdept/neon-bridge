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

type ProwlarrTestConfig struct {
	ServerURL string `json:"serverUrl" binding:"required"`
	ApiKey    string `json:"apiKey" binding:"required"`
}

type ProwlarrIndexerStats struct {
	IndexerID                 int    `json:"indexerId"`
	IndexerName               string `json:"indexerName"`
	AverageResponseTime       int    `json:"averageResponseTime"`
	AverageGrabResponseTime   int    `json:"averageGrabResponseTime"`
	NumberOfQueries           int    `json:"numberOfQueries"`
	NumberOfGrabs             int    `json:"numberOfGrabs"`
	NumberOfRssQueries        int    `json:"numberOfRssQueries"`
	NumberOfAuthQueries       int    `json:"numberOfAuthQueries"`
	NumberOfFailedQueries     int    `json:"numberOfFailedQueries"`
	NumberOfFailedGrabs       int    `json:"numberOfFailedGrabs"`
	NumberOfFailedRssQueries  int    `json:"numberOfFailedRssQueries"`
	NumberOfFailedAuthQueries int    `json:"numberOfFailedAuthQueries"`
}

type ProwlarrIndexerStatsResponse struct {
	ID         int                      `json:"id"`
	Indexers   []ProwlarrIndexerStats   `json:"indexers"`
	UserAgents []map[string]interface{} `json:"userAgents"`
	Hosts      []map[string]interface{} `json:"hosts"`
}

type ProwlarrHealthCheck struct {
	Message string `json:"message"`
}

type ProwlarrStats struct {
	TotalQueries       int     `json:"totalQueries"`
	TotalGrabs         int     `json:"totalGrabs"`
	TotalFailedQueries int     `json:"totalFailedQueries"`
	ActiveIndexers     int     `json:"activeIndexers"`
	Alerts             []Alert `json:"alerts"`
}

func TestProwlarrConnection(c *gin.Context) {
	var config ProwlarrTestConfig
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

	stats, err := fetchProwlarrStats(config.ServerURL, config.ApiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func ProxyProwlarrStats(c *gin.Context) {
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

	if widget.Type != "prowlarr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not a Prowlarr widget"})
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

	stats, err := fetchProwlarrStats(serverURL, apiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func fetchProwlarrStats(serverURL, apiKey string) (*ProwlarrStats, error) {
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

	stats := &ProwlarrStats{}
	indexerStats, err := fetchProwlarrIndexerStats(client, serverURL, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indexer stats: %v", err)
	}

	totalQueries := 0
	totalGrabs := 0
	totalFailedQueries := 0
	activeIndexers := 0

	for _, indexer := range indexerStats.Indexers {
		totalQueries += indexer.NumberOfQueries
		totalGrabs += indexer.NumberOfGrabs
		totalFailedQueries += indexer.NumberOfFailedQueries
		activeIndexers++
	}

	stats.TotalQueries = totalQueries
	stats.TotalGrabs = totalGrabs
	stats.TotalFailedQueries = totalFailedQueries
	stats.ActiveIndexers = activeIndexers

	var alerts = make([]Alert, 0)

	if healthChecks, err := fetchProwlarrHealth(client, serverURL, apiKey); err == nil {
		for _, health := range healthChecks {
			alerts = append(alerts, Alert{
				Message: health.Message,
				Level:   "warning",
			})

		}
	}
	if totalQueries > 0 {
		failureRate := float64(totalFailedQueries) / float64(totalQueries) * 100
		if failureRate > 20 { // More than 20% failure rate
			alerts = append(alerts, Alert{
				Message: fmt.Sprintf("High failure rate: %.1f%% of queries are failing", failureRate),
				Level:   "warning",
			})
		}
	}

	if activeIndexers == 0 {
		alerts = append(alerts, Alert{
			Message: "No active indexers found",
			Level:   "error",
		})
	}

	stats.Alerts = alerts

	return stats, nil
}

func fetchProwlarrIndexerStats(client *http.Client, serverURL, apiKey string) (*ProwlarrIndexerStatsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/indexerstats?apikey=%s", serverURL, apiKey)
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

	var stats ProwlarrIndexerStatsResponse
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func fetchProwlarrHealth(client *http.Client, serverURL, apiKey string) ([]ProwlarrHealthCheck, error) {
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

	var healthChecks []ProwlarrHealthCheck
	if err := json.Unmarshal(body, &healthChecks); err != nil {
		return nil, err
	}

	return healthChecks, nil
}
