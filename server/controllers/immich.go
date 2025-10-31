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

type ImmichTestConfig struct {
	ServerURL string `json:"serverUrl" binding:"required"`
	ApiKey    string `json:"apiKey" binding:"required"`
}

type ImmichAboutInfo struct {
	Version string `json:"version"`
}

type ImmichVersionCheck struct {
	ReleaseVersion string `json:"releaseVersion"`
}

type ImmichServerStatistics struct {
	Photos      int64               `json:"photos"`
	Videos      int64               `json:"videos"`
	Usage       int64               `json:"usage"`
	UsageByUser []ImmichUsageByUser `json:"usageByUser"`
}

type ImmichUsageByUser struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Photos   int64  `json:"photos"`
	Videos   int64  `json:"videos"`
	Usage    int64  `json:"usage"`
}

type ImmichStorage struct {
	DiskAvailable      string  `json:"diskAvailable"`
	DiskSize           string  `json:"diskSize"`
	DiskUsedPercentage float32 `json:"diskUsagePercentage"`
	DiskUsed           string  `json:"diskUse"`
}

type ImmichNotifications struct {
	Total int64 `json:"total"`
}

type ImmichStats struct {
	ServerStats   ImmichServerStatistics `json:"serverStats"`
	Storage       ImmichStorage          `json:"storage"`
	Users         int64                  `json:"users"`
	Alerts       []Alert                `json:"alerts"`
}

func TestImmichConnection(c *gin.Context) {
	var config ImmichTestConfig
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

	stats, err := fetchImmichStats(config.ServerURL, config.ApiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func ProxyImmichStats(c *gin.Context) {
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

	if widget.Type != "immich" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not an Immich widget"})
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

	stats, err := fetchImmichStats(serverURL, apiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func fetchImmichStats(serverURL, apiKey string) (*ImmichStats, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	stats := &ImmichStats{}

	serverStats, err := fetchImmichServerStatistics(client, serverURL, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server statistics: %v", err)
	}

	stats.Users = int64(len(serverStats.UsageByUser))
	stats.ServerStats = *serverStats

	storage, err := fetchImmichStorage(client, serverURL, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch storage: %v", err)
	}
	stats.Storage = *storage

	about, err := fetchImmichAbout(client, serverURL, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch about info: %v", err)
	}

	var alerts = make([]Alert, 0)

	notificationCount, err := fetchImmichNotifications(client, serverURL, apiKey)

	if err != nil {
		fmt.Printf("Warning: failed to fetch notifications: %v\n", err)
	} else if notificationCount > 0 {
		alerts = append(alerts, Alert{
			Message: fmt.Sprintf("You have %d unread notifications", notificationCount),
			Level:    "warning",
		})
	}

	versionCheck, err := fetchImmichVersionCheck(client, serverURL, apiKey)
	if err != nil {
		fmt.Printf("Warning: failed to fetch version check: %v\n", err)
	} else {
		if versionCheck.ReleaseVersion != about.Version {
			message := fmt.Sprintf("A new Immich version %s is available! You are running version %s.", versionCheck.ReleaseVersion, about.Version)
			alerts = append(alerts, Alert{
				Message: message,
				Level:    "warning",
			})
		}
	}

	stats.Alerts = alerts

	return stats, nil
}

func fetchImmichServerStatistics(client *http.Client, serverURL, apiKey string) (*ImmichServerStatistics, error) {
	url := fmt.Sprintf("%s/api/server/statistics", serverURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var stats ImmichServerStatistics
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func fetchImmichStorage(client *http.Client, serverURL, apiKey string) (*ImmichStorage, error) {
	url := fmt.Sprintf("%s/api/server/storage", serverURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var storage ImmichStorage
	if err := json.NewDecoder(resp.Body).Decode(&storage); err != nil {
		return nil, err
	}

	return &storage, nil
}

func fetchImmichAbout(client *http.Client, serverURL, apiKey string) (*ImmichAboutInfo, error) {
	url := fmt.Sprintf("%s/api/server/about", serverURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var about ImmichAboutInfo
	if err := json.NewDecoder(resp.Body).Decode(&about); err != nil {
		return nil, err
	}

	return &about, nil
}

func fetchImmichNotifications(client *http.Client, serverURL, apiKey string) (int64, error) {
	url := fmt.Sprintf("%s/api/notifications", serverURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var notificationsArray []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&notificationsArray); err != nil {
		return 0, err
	}

	return int64(len(notificationsArray)), nil
}

func fetchImmichVersionCheck(client *http.Client, serverURL, apiKey string) (*ImmichVersionCheck, error) {
	url := fmt.Sprintf("%s/api/server/version-check", serverURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var versionCheck ImmichVersionCheck
	if err := json.NewDecoder(resp.Body).Decode(&versionCheck); err != nil {
		return nil, err
	}

	return &versionCheck, nil
}
