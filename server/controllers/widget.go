package controllers

import (
	"net/http"
	"strconv"

	"dashboard-server/models"
	"dashboard-server/database"

	"github.com/gin-gonic/gin"
)
func GetWidgets(c *gin.Context) {
	dashboardID := c.Param("id")

	var widgets []models.Widget
	result := database.DB.Where("dashboard_id = ?", dashboardID).Find(&widgets)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var widgetResponses []models.WidgetResponse
	for _, widget := range widgets {
		widgetResponses = append(widgetResponses, widget.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"data": widgetResponses})
}

func GetWidget(c *gin.Context) {
	id := c.Param("id")

	var widget models.Widget
	result := database.DB.Preload("Dashboard").First(&widget, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": widget.ToResponse()})
}

func CreateWidget(c *gin.Context) {
	dashboardID := c.Param("id")

	var widget models.Widget
	if err := c.ShouldBindJSON(&widget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(dashboardID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dashboard ID"})
		return
	}
	widget.DashboardID = uint(id)

	var dashboard models.Dashboard
	if err := database.DB.First(&dashboard, widget.DashboardID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dashboard not found"})
		return
	}

	result := database.DB.Create(&widget)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": widget.ToResponse()})
}

func UpdateWidget(c *gin.Context) {
	id := c.Param("id")

	var widget models.Widget
	if err := database.DB.First(&widget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
		return
	}

	if err := c.ShouldBindJSON(&widget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&widget)
	c.JSON(http.StatusOK, gin.H{"data": widget.ToResponse()})
}

func UpdateWidgetState(c *gin.Context) {
	id := c.Param("id")

	var widget models.Widget
	if err := database.DB.First(&widget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
		return
	}

	var stateData struct {
		LastState models.JSON `json:"last_state"`
	}

	if err := c.ShouldBindJSON(&stateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	widget.LastState = stateData.LastState
	database.DB.Save(&widget)

	c.JSON(http.StatusOK, gin.H{"data": widget.ToResponse()})
}

func DeleteWidget(c *gin.Context) {
	id := c.Param("id")

	var widget models.Widget
	if err := database.DB.First(&widget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
		return
	}

	database.DB.Delete(&widget)
	c.JSON(http.StatusOK, gin.H{"message": "Widget deleted successfully"})
}
