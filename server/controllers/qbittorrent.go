package controllers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dashboard-server/database"
	"dashboard-server/models"

	"github.com/gin-gonic/gin"
)

type QBittorrentConfig struct {
	ServerURL        string `json:"serverUrl"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	MaxDownloadSpeed int    `json:"maxDownloadSpeed"`
	MaxUploadSpeed   int    `json:"maxUploadSpeed"`
}

type QBittorrentStats struct {
	DownloadingTorrents int     `json:"downloadingTorrents"`
	SeedingTorrents     int     `json:"seedingTorrents"`
	ErrorTorrents       int     `json:"errorTorrents"`
	TotalTorrents       int     `json:"totalTorrents"`
	DownloadSpeed       float64 `json:"downloadSpeed"`
	UploadSpeed         float64 `json:"uploadSpeed"`
}

type QBittorrentTorrent struct {
	Hash     string  `json:"hash"`
	Name     string  `json:"name"`
	State    string  `json:"state"`
	Progress float64 `json:"progress"`
	Dlspeed  int64   `json:"dlspeed"`
	Upspeed  int64   `json:"upspeed"`
	Priority int     `json:"priority"`
	Size     int64   `json:"size"`
}

type QBittorrentGlobalStats struct {
	DlInfoSpeed int64 `json:"dl_info_speed"`
	UpInfoSpeed int64 `json:"up_info_speed"`
}

func ProxyQBittorrentStats(c *gin.Context) {
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

	if widget.Type != "qbittorrent" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not a qBittorrent widget"})
		return
	}

	config := widget.Config
	serverURL, ok := config["serverUrl"].(string)
	if !ok || serverURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serverUrl not found in widget configuration"})
		return
	}

	qbitConfig := QBittorrentConfig{
		ServerURL: serverURL,
	}

	if username, ok := config["username"].(string); ok {
		qbitConfig.Username = username
	}
	if password, ok := config["password"].(string); ok {
		qbitConfig.Password = password
	}
	if maxDownload, ok := config["maxDownloadSpeed"].(float64); ok {
		qbitConfig.MaxDownloadSpeed = int(maxDownload)
	}
	if maxUpload, ok := config["maxUploadSpeed"].(float64); ok {
		qbitConfig.MaxUploadSpeed = int(maxUpload)
	}

	stats, err := fetchQBittorrentStats(qbitConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to fetch qBittorrent stats: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

func TestQBittorrentConnection(c *gin.Context) {
	var config QBittorrentConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration", "details": err.Error()})
		return
	}

	_, err := fetchQBittorrentStats(config)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Connection test failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully connected to qBittorrent",
	})
}

func fetchQBittorrentStats(config QBittorrentConfig) (*QBittorrentStats, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	baseURL := strings.TrimSuffix(config.ServerURL, "/")

	var cookie string
	if config.Username != "" && config.Password != "" {
		loginCookie, err := qbittorrentLogin(client, baseURL, config.Username, config.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to login: %v", err)
		}
		cookie = loginCookie
	}

	torrents, err := getQBittorrentTorrents(client, baseURL, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to get torrents: %v", err)
	}

	globalStats, err := getQBittorrentGlobalStats(client, baseURL, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to get global stats: %v", err)
	}

	stats := &QBittorrentStats{
		DownloadSpeed: float64(globalStats.DlInfoSpeed),
		UploadSpeed:   float64(globalStats.UpInfoSpeed),
	}

	for _, torrent := range torrents {
		stats.TotalTorrents++

		switch torrent.State {
		case "downloading", "queuedDL", "allocating", "metaDL", "pausedDL", "forcedDL":
			stats.DownloadingTorrents++
		case "uploading", "queuedUP", "forcedUP":
			stats.SeedingTorrents++
		case "error", "missingFiles", "stalledDL", "stalledUP":
			stats.ErrorTorrents++
		}
	}

	return stats, nil
}

func qbittorrentLogin(client *http.Client, baseURL, username, password string) (string, error) {
	loginURL := baseURL + "/api/v2/auth/login"

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	if string(body) != "Ok." {
		return "", fmt.Errorf("login failed: %s", string(body))
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			return fmt.Sprintf("%s=%s", cookie.Name, cookie.Value), nil
		}
	}

	return "", fmt.Errorf("no session cookie found")
}

func getQBittorrentTorrents(client *http.Client, baseURL, cookie string) ([]QBittorrentTorrent, error) {
	torrentURL := baseURL + "/api/v2/torrents/info"

	req, err := http.NewRequest("GET", torrentURL, nil)
	if err != nil {
		return nil, err
	}

	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var torrents []QBittorrentTorrent
	if err := json.NewDecoder(resp.Body).Decode(&torrents); err != nil {
		return nil, err
	}

	return torrents, nil
}

func getQBittorrentGlobalStats(client *http.Client, baseURL, cookie string) (*QBittorrentGlobalStats, error) {
	statsURL := baseURL + "/api/v2/transfer/info"

	req, err := http.NewRequest("GET", statsURL, nil)
	if err != nil {
		return nil, err
	}

	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var stats QBittorrentGlobalStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
