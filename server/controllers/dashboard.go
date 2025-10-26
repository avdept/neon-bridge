package controllers

import (
	"net/http"

	"dashboard-server/models"
	"dashboard-server/database"
	"github.com/gin-gonic/gin"

)


func GetDashboards(c *gin.Context) {
	var dashboards []models.Dashboard

	result := database.DB.Preload("Widgets").Find(&dashboards)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var dashboardResponses []models.DashboardResponse
	for _, dashboard := range dashboards {
		dashboardResponses = append(dashboardResponses, dashboard.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboardResponses})
}

func GetDashboard(c *gin.Context) {
	id := c.Param("id")

	var dashboard models.Dashboard
	result := database.DB.Preload("Widgets").First(&dashboard, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboard.ToResponse()})
}

func CreateDashboard(c *gin.Context) {
	var dashboard models.Dashboard

	if err := c.ShouldBindJSON(&dashboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&dashboard)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": dashboard.ToResponse()})
}

func UpdateDashboard(c *gin.Context) {
	id := c.Param("id")

	var dashboard models.Dashboard
	if err := database.DB.First(&dashboard, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
		return
	}

	if err := c.ShouldBindJSON(&dashboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&dashboard)
	c.JSON(http.StatusOK, gin.H{"data": dashboard.ToResponse()})
}

func DeleteDashboard(c *gin.Context) {
	id := c.Param("id")

	var dashboard models.Dashboard
	if err := database.DB.First(&dashboard, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
		return
	}

	database.DB.Where("dashboard_id = ?", id).Delete(&models.Widget{})

	database.DB.Delete(&dashboard)

	c.JSON(http.StatusOK, gin.H{"message": "Dashboard deleted successfully"})
}
