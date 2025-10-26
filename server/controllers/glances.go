package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"dashboard-server/database"
	"dashboard-server/models"

	"github.com/gin-gonic/gin"
)

var glancesHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

type GlancesConfig struct {
	URL      string `json:"url"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type GlancesStats struct {
	CPU struct {
		Usage       float64 `json:"usage"`
		Temperature float64 `json:"temperature"`
	} `json:"cpu"`
	Memory struct {
		Used       float64 `json:"used"`
		Total      float64 `json:"total"`
		Percentage float64 `json:"percentage"`
	} `json:"memory"`
	Uptime struct {
		Days    int    `json:"days"`
		Display string `json:"display"`
	} `json:"uptime"`
	LoadAverage float64 `json:"loadAverage"`
	Processes   int     `json:"processes"`
}

func GetGlancesStats(c *gin.Context) {
	dashboardID := c.Param("id")

	var dashboard models.Dashboard
	if err := database.DB.First(&dashboard, dashboardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
		return
	}

	if dashboard.GlancesConfig == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Glances configuration found for this dashboard"})
		return
	}

	var config GlancesConfig
	if err := json.Unmarshal([]byte(dashboard.GlancesConfig), &config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Glances configuration"})
		return
	}

	stats, err := fetchGlancesData(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch Glances data: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func fetchGlancesData(config GlancesConfig) (*GlancesStats, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/4/all", config.URL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := glancesHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("glances API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	var glancesResp map[string]interface{}
	if err := json.Unmarshal(body, &glancesResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	stats := transformGlancesData(glancesResp)
	return stats, nil
}

func transformGlancesData(data map[string]interface{}) *GlancesStats {
	stats := &GlancesStats{}

	if cpuData, ok := data["cpu"].(map[string]interface{}); ok {
		if usage, ok := cpuData["total"].(float64); ok {
			stats.CPU.Usage = usage
		}
	}

	if memData, ok := data["mem"].(map[string]interface{}); ok {
		if used, ok := memData["used"].(float64); ok {
			stats.Memory.Used = used / (1024 * 1024 * 1024)
		}
		if total, ok := memData["total"].(float64); ok {
			stats.Memory.Total = total / (1024 * 1024 * 1024)
		}
		if percent, ok := memData["percent"].(float64); ok {
			stats.Memory.Percentage = percent
		}
	}

	if uptimeData, ok := data["uptime"].(float64); ok {
		days := int(uptimeData / (24 * 3600))
		stats.Uptime.Days = days
		if days > 0 {
			stats.Uptime.Display = fmt.Sprintf("%dd", days)
		} else {
			hours := int(uptimeData / 3600)
			stats.Uptime.Display = fmt.Sprintf("%dh", hours)
		}
	}
	if loadData, ok := data["load"].(map[string]interface{}); ok {
		if load1, ok := loadData["min1"].(float64); ok {
			stats.LoadAverage = load1
		}
	}
	if processes, ok := data["processcount"].(map[string]interface{}); ok {
		if running, ok := processes["running"].(float64); ok {
			stats.Processes = int(running)
		}
	}

	stats.CPU.Temperature = 0

	return stats
}

func TestGlancesConnection(c *gin.Context) {
	var config GlancesConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/3/all", config.URL), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL", "success": false})
		return
	}

	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "success": false})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Connection successful"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": fmt.Sprintf("HTTP %d", resp.StatusCode)})
	}
}
