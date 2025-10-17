package controllers

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

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
	stats := &SystemStats{}

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		stats.CPU.Usage = cpuPercent[0]
	}

	stats.CPU.Temperature = getCPUTemperature()

	memInfo, err := mem.VirtualMemory()
	if err == nil {
		stats.Memory.Used = float64(memInfo.Used) / (1024 * 1024 * 1024)   // Convert to GB
		stats.Memory.Total = float64(memInfo.Total) / (1024 * 1024 * 1024) // Convert to GB
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

// TODO: This is example solution I found on SO - so keeping it for now. Not sure if it will work in docker, but works on host apple ARM machine.
func getCPUTemperature() float64 {
	// Try to get temperature from gopsutil sensors
	temps, err := sensors.TemperaturesWithContext(context.TODO())
	if err != nil {
		return 0
	}

	for _, temp := range temps {
		sensorKey := temp.SensorKey
		if sensorKey == "coretemp_core_0_input" || // Intel Linux
			sensorKey == "cpu_thermal" || // Generic
			sensorKey == "Package id 0" || // Intel
			sensorKey == "k10temp_tctl" || // AMD
			sensorKey == "acpi_0" || // ACPI thermal
			sensorKey == "thermal_zone0" || // Generic thermal zone
			len(temps) == 1 { // If only one sensor, assume it's CPU
			if temp.Temperature > 0 && temp.Temperature < 150 { // Reasonable temperature range (0-150°C)
				return temp.Temperature
			}
		}
	}

	for _, temp := range temps {
		if temp.Temperature > 0 && temp.Temperature < 150 {
			return temp.Temperature
		}
	}

	// Default fallback temperature if no sensors work
	return 0
}
