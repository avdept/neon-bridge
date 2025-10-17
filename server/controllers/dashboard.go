package controllers

import (
	"net/http"

	"dashboard-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDB(database *gorm.DB) {
	DB = database
}

// GetDashboards retrieves all dashboards
func GetDashboards(c *gin.Context) {
	var dashboards []models.Dashboard

	result := DB.Preload("Widgets").Find(&dashboards)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Convert to response format with filtered configs
	var dashboardResponses []models.DashboardResponse
	for _, dashboard := range dashboards {
		dashboardResponses = append(dashboardResponses, dashboard.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboardResponses})
}

// GetDashboard retrieves a single dashboard by ID
func GetDashboard(c *gin.Context) {
	id := c.Param("id")

	var dashboard models.Dashboard
	result := DB.Preload("Widgets").First(&dashboard, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboard.ToResponse()})
}

// CreateDashboard creates a new dashboard
func CreateDashboard(c *gin.Context) {
	var dashboard models.Dashboard

	if err := c.ShouldBindJSON(&dashboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&dashboard)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": dashboard.ToResponse()})
}

// UpdateDashboard updates an existing dashboard
func UpdateDashboard(c *gin.Context) {
	id := c.Param("id")

	var dashboard models.Dashboard
	if err := DB.First(&dashboard, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
		return
	}

	if err := c.ShouldBindJSON(&dashboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Save(&dashboard)
	c.JSON(http.StatusOK, gin.H{"data": dashboard.ToResponse()})
}

// DeleteDashboard deletes a dashboard
func DeleteDashboard(c *gin.Context) {
	id := c.Param("id")

	var dashboard models.Dashboard
	if err := DB.First(&dashboard, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
		return
	}

	// Delete all widgets belonging to this dashboard first
	DB.Where("dashboard_id = ?", id).Delete(&models.Widget{})

	// Delete the dashboard
	DB.Delete(&dashboard)

	c.JSON(http.StatusOK, gin.H{"message": "Dashboard deleted successfully"})
}
