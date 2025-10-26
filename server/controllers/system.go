package controllers

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"dashboard-server/database"
	"dashboard-server/services"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

type SystemStats struct {
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

func GetSystemStats(c *gin.Context) {
	glancesService := services.NewGlancesService(database.DB)

	if config, err := glancesService.GetGlancesConfigFromFirstDashboard(); err == nil {
		if glancesStats, err := glancesService.FetchGlancesStats(config); err == nil {
			stats := &SystemStats{}

			stats.CPU.Usage = glancesStats.CPU.Usage
			stats.CPU.Temperature = glancesStats.CPU.Temperature

			stats.Memory.Used = glancesStats.Memory.Used
			stats.Memory.Total = glancesStats.Memory.Total
			stats.Memory.Percentage = glancesStats.Memory.Percentage

			stats.Uptime.Days = glancesStats.Uptime.Days
			stats.Uptime.Display = glancesStats.Uptime.Display

			stats.LoadAverage = glancesStats.LoadAverage
			stats.Processes = glancesStats.Processes

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    stats,
				"source":  "glances",
			})
			return
		}
		fmt.Printf("Failed to fetch stats from Glances: %v\n", err)
	}

	stats := &SystemStats{}

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		stats.CPU.Usage = cpuPercent[0]
	}

	stats.CPU.Temperature = getCPUTemperature()

	memInfo, err := mem.VirtualMemory()
	if err == nil {
		stats.Memory.Used = float64(memInfo.Used) / (1024 * 1024 * 1024)
		stats.Memory.Total = float64(memInfo.Total) / (1024 * 1024 * 1024)
		stats.Memory.Percentage = memInfo.UsedPercent
	}
	hostInfo, err := host.Info()
	if err == nil {
		uptimeDays := int(hostInfo.Uptime / (24 * 3600))
		stats.Uptime.Days = uptimeDays
		stats.Uptime.Display = formatUptime(hostInfo.Uptime)
	}

	loadInfo, err := load.Avg()
	if err == nil {
		stats.LoadAverage = loadInfo.Load1
	}
	stats.Processes = runtime.NumGoroutine()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"source":  "local",
	})
}

func formatUptime(uptimeSeconds uint64) string {
	days := uptimeSeconds / (24 * 3600)
	hours := (uptimeSeconds % (24 * 3600)) / 3600
	minutes := (uptimeSeconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

func getCPUTemperature() float64 {
	temps, err := sensors.TemperaturesWithContext(context.TODO())
	if err != nil {
		return 0
	}

	for _, temp := range temps {
		sensorKey := temp.SensorKey
		if sensorKey == "coretemp_core_0_input" ||
			sensorKey == "cpu_thermal" ||
			sensorKey == "Package id 0" ||
			sensorKey == "k10temp_tctl" ||
			sensorKey == "acpi_0" ||
			sensorKey == "thermal_zone0" ||
			len(temps) == 1 {
			if temp.Temperature > 0 && temp.Temperature < 150 {
				return temp.Temperature
			}
		}
	}

	for _, temp := range temps {
		if temp.Temperature > 0 && temp.Temperature < 150 {
			return temp.Temperature
		}
	}

	return 0
}
