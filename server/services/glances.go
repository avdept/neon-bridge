package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"dashboard-server/models"

	"gorm.io/gorm"
)

type GlancesService struct {
	db *gorm.DB
}

func NewGlancesService(db *gorm.DB) *GlancesService {
	return &GlancesService{
		db: db,
	}
}

type GlancesConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
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
	LoadAverage float64 `json:"load_average"`
	Processes   int     `json:"processes"`
}

func (s *GlancesService) GetGlancesConfigFromFirstDashboard() (*GlancesConfig, error) {
	var dashboard models.Dashboard
	if err := s.db.First(&dashboard).Error; err != nil {
		return nil, fmt.Errorf("no dashboard found: %w", err)
	}

	if dashboard.GlancesConfig == "" {
		return nil, fmt.Errorf("no Glances configuration found")
	}

	var config GlancesConfig
	if err := json.Unmarshal([]byte(dashboard.GlancesConfig), &config); err != nil {
		return nil, fmt.Errorf("invalid Glances configuration: %w", err)
	}

	if config.URL == "" {
		return nil, fmt.Errorf("Glances URL not configured")
	}

	return &config, nil
}

func (s *GlancesService) FetchGlancesStats(config *GlancesConfig) (*GlancesStats, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/4/all", config.URL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(req)
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

	if uptimeData, ok := data["uptime"].(string); ok {
		stats.Uptime.Display = uptimeData
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

	if sensorsData, ok := data["sensors"].([]interface{}); ok {
		for _, sensor := range sensorsData {
			if sensorMap, ok := sensor.(map[string]interface{}); ok {
				if label, labelOk := sensorMap["label"].(string); labelOk {
					if label == "Package id 0" || label == "Core 0" {
						if value, valueOk := sensorMap["value"].(float64); valueOk {
							stats.CPU.Temperature = value
							break
						}
					}
				}
			}
		}
	}

	return stats
}

func (s *GlancesService) TestGlancesConnection(config *GlancesConfig) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/4/all", config.URL), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status %d", resp.StatusCode)
	}

	return nil
}
