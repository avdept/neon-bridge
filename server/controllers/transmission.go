package controllers

import (
	"bytes"
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

type TransmissionConfig struct {
	ServerURL        string `json:"serverUrl"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	RPCPath          string `json:"rpcPath"`
	MaxDownloadSpeed int    `json:"maxDownloadSpeed"`
	MaxUploadSpeed   int    `json:"maxUploadSpeed"`
}

type TransmissionStats struct {
	DownloadingTorrents int     `json:"downloadingTorrents"`
	SeedingTorrents     int     `json:"seedingTorrents"`
	ErrorTorrents       int     `json:"errorTorrents"`
	TotalTorrents       int     `json:"totalTorrents"`
	DownloadSpeed       float64 `json:"downloadSpeed"` // bytes per second
	UploadSpeed         float64 `json:"uploadSpeed"`   // bytes per second
}

type TransmissionRPCRequest struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments"`
	Tag       int         `json:"tag,omitempty"`
}

type TransmissionRPCResponse struct {
	Result    string      `json:"result"`
	Arguments interface{} `json:"arguments"`
	Tag       int         `json:"tag,omitempty"`
}

type TorrentGetArguments struct {
	Fields []string `json:"fields"`
}

type TorrentGetResponse struct {
	Torrents []Torrent `json:"torrents"`
}

type Torrent struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Status       int     `json:"status"`
	PercentDone  float64 `json:"percentDone"`
	RateDownload float64 `json:"rateDownload"` // bytes per second
	RateUpload   float64 `json:"rateUpload"`   // bytes per second
	Error        int     `json:"error"`        // error status
}

type SessionStats struct {
	DownloadSpeed float64 `json:"downloadSpeed"`
	UploadSpeed   float64 `json:"uploadSpeed"`
	TorrentCount  int     `json:"torrentCount"`
}

func ProxyTransmissionStats(c *gin.Context) {
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

	if widget.Type != "transmission" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Widget is not a Transmission widget"})
		return
	}

	config := widget.Config
	serverURL, ok := config["serverUrl"].(string)
	if !ok || serverURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serverUrl not found in widget configuration"})
		return
	}

	transmissionConfig := TransmissionConfig{
		ServerURL: serverURL,
	}

	if username, ok := config["username"].(string); ok {
		transmissionConfig.Username = username
	}
	if password, ok := config["password"].(string); ok {
		transmissionConfig.Password = password
	}
	if rpcPath, ok := config["rpcPath"].(string); ok {
		transmissionConfig.RPCPath = rpcPath
	}
	if maxDownload, ok := config["maxDownloadSpeed"].(float64); ok {
		transmissionConfig.MaxDownloadSpeed = int(maxDownload)
	}
	if maxUpload, ok := config["maxUploadSpeed"].(float64); ok {
		transmissionConfig.MaxUploadSpeed = int(maxUpload)
	}

	stats, err := fetchTransmissionStats(transmissionConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to fetch Transmission stats: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

func TestTransmissionConnection(c *gin.Context) {
	var config TransmissionConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration", "details": err.Error()})
		return
	}

	_, err := fetchTransmissionStats(config)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Connection test failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully connected to Transmission",
	})
}

func fetchTransmissionStats(config TransmissionConfig) (*TransmissionStats, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	baseURL := config.ServerURL
	if config.RPCPath == "" {
		config.RPCPath = "/transmission/rpc"
	}
	url := baseURL + config.RPCPath

	sessionID, err := getTransmissionSessionID(client, url, config.Username, config.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get session ID: %v", err)
	}

	torrents, err := getTransmissionTorrents(client, url, config.Username, config.Password, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get torrents: %v", err)
	}

	stats := &TransmissionStats{}

	for _, torrent := range torrents {
		stats.TotalTorrents++
		stats.DownloadSpeed += torrent.RateDownload
		stats.UploadSpeed += torrent.RateUpload

		if torrent.Error != 0 {
			stats.ErrorTorrents++
		}

		// status values: 0=stopped, 1=check-wait, 2=check, 3=download-wait, 4=download, 5=seed-wait, 6=seed
		switch torrent.Status {
		case 4: // downloading
			stats.DownloadingTorrents++
		case 6: // seeding
			stats.SeedingTorrents++
		}
	}

	return stats, nil
}
func getTransmissionSessionID(client *http.Client, url, username, password string) (string, error) {
	// Create a dummy request to get the session ID
	reqBody := TransmissionRPCRequest{
		Method: "session-get",
		Arguments: map[string]interface{}{
			"fields": []string{"version"},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 409 {
		sessionID := resp.Header.Get("X-Transmission-Session-Id")
		if sessionID != "" {
			return sessionID, nil
		}
	}

	if resp.StatusCode == 200 {
		return "", nil
	}

	return "", fmt.Errorf("unexpected response status: %d", resp.StatusCode)
}

func getTransmissionTorrents(client *http.Client, url, username, password, sessionID string) ([]Torrent, error) {
	reqBody := TransmissionRPCRequest{
		Method: "torrent-get",
		Arguments: TorrentGetArguments{
			Fields: []string{"id", "name", "status", "percentDone", "rateDownload", "rateUpload", "error"},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if sessionID != "" {
		req.Header.Set("X-Transmission-Session-Id", sessionID)
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
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

	var rpcResp TransmissionRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Result != "success" {
		return nil, fmt.Errorf("RPC error: %s", rpcResp.Result)
	}

	argsJSON, err := json.Marshal(rpcResp.Arguments)
	if err != nil {
		return nil, err
	}

	var torrentResp TorrentGetResponse
	if err := json.Unmarshal(argsJSON, &torrentResp); err != nil {
		return nil, err
	}

	return torrentResp.Torrents, nil
}
