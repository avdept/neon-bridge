package controllers

import (
	"encoding/base64"
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

type AdGuardTestConfig struct {
	ServerURL string `json:"serverUrl" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type AdGuardVersionResponse struct {
	NewVersion      string `json:"new_version"`
	Announcement    string `json:"announcement"`
	AnnouncementURL string `json:"announcement_url"`
	CanAutoUpdate   bool   `json:"can_autoupdate"`
	Disabled        bool   `json:"disabled"`
}

type AdGuardStatsResponse struct {
	NumDNSQueries       int     `json:"num_dns_queries"`
	NumBlockedFiltering int     `json:"num_blocked_filtering"`
	AvgProcessingTime   float64 `json:"avg_processing_time"`
	TimeUnits           string  `json:"time_units"`
	Health              string  `json:"health,omitempty"`
	TotalQueries        int     `json:"totalQueries"`
	BlockedQueries      int     `json:"blockedQueries"`
	BlockingPercentage  float64 `json:"blockingPercentage"`
	TimeUnit            string  `json:"timeUnit"`
}

func TestAdGuardConnection(c *gin.Context) {
	var config AdGuardTestConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration: " + err.Error()})
		return
	}

	if config.ServerURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serverUrl is required"})
		return
	}

	if config.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	if config.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	statsResponse, statusCode, err := proxyAdGuardRequest(config.ServerURL, config.Username, config.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.Data(statusCode, "application/json", statsResponse)
}

func ProxyAdGuardStats(c *gin.Context) {
	widgetIDStr := c.Param("widget_id")
	widgetID, err := strconv.ParseUint(widgetIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid widget ID"})
		return
	}

	var widget models.Widget
	if err := database.DB.First(&widget, uint(widgetID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
		return
	}

	if widget.Type != "adguard-home" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not an AdGuard Home widget"})
		return
	}

	config := widget.Config
	serverURL, ok := config["serverUrl"].(string)
	if !ok || serverURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serverUrl not found in widget configuration"})
		return
	}

	username, ok := config["username"].(string)
	if !ok || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username not found in widget configuration"})
		return
	}

	password, ok := config["password"].(string)
	if !ok || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password not found in widget configuration"})
		return
	}

	statsResponse, statusCode, err := proxyAdGuardRequest(serverURL, username, password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.Data(statusCode, "application/json", statsResponse)
}

func proxyAdGuardRequest(serverURL, username, password string) ([]byte, int, error) {
	if serverURL[len(serverURL)-1] == '/' {
		serverURL = serverURL[:len(serverURL)-1]
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	statsURL := fmt.Sprintf("%s/control/stats", serverURL)
	statsData, statusCode, err := makeAdGuardRequest(client, statsURL, username, password)
	if err != nil || statusCode != http.StatusOK {
		return nil, statusCode, err
	}

	stats := &AdGuardStatsResponse{}

	fmt.Printf("DEBUG AdGuard: Raw response: %s\n", string(statsData))

	if err := json.Unmarshal(statsData, &stats); err != nil {
		fmt.Printf("DEBUG AdGuard: Unmarshal error: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to parse stats response: %v", err)
	}

	stats.TotalQueries = stats.NumDNSQueries
	stats.BlockedQueries = stats.NumBlockedFiltering
	stats.TimeUnit = stats.TimeUnits
	if stats.TotalQueries > 0 {
		stats.BlockingPercentage = float64(stats.BlockedQueries) / float64(stats.TotalQueries) * 100
	} else {
		stats.BlockingPercentage = 0
	}

	versionURL := fmt.Sprintf("%s/control/version.json", serverURL)
	versionData, _, versionErr := makeAdGuardRequest(client, versionURL, username, password)

	if versionErr == nil {
		var version AdGuardVersionResponse
		if json.Unmarshal(versionData, &version) == nil {
			stats.Health = version.Announcement
		}
	}

	responseData, err := json.Marshal(stats)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to marshal response: %v", err)
	}

	return responseData, statusCode, nil
}

func makeAdGuardRequest(client *http.Client, url, username, password string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %v", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Homepage-Dashboard/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %v", err)
	}

	return body, resp.StatusCode, nil
}
