package routes

import (
	"dashboard-server/controllers"
	"dashboard-server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORS())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Dashboard routes
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("", controllers.GetDashboards)
			dashboards.POST("", controllers.CreateDashboard)
			dashboards.GET("/:id", controllers.GetDashboard)
			dashboards.PUT("/:id", controllers.UpdateDashboard)
			dashboards.DELETE("/:id", controllers.DeleteDashboard)

			// Widget routes nested under specific dashboard
			dashboards.GET("/:id/widgets", controllers.GetWidgets)
			dashboards.POST("/:id/widgets", controllers.CreateWidget)
		}

		// Widget routes (for direct access)
		widgets := v1.Group("/widgets")
		{
			widgets.GET("/:id", controllers.GetWidget)
			widgets.PUT("/:id", controllers.UpdateWidget)
			widgets.PUT("/:id/state", controllers.UpdateWidgetState)
			widgets.DELETE("/:id", controllers.DeleteWidget)
		}

		// AdGuard proxy routes
		v1.GET("/adguard/:widget_id", controllers.ProxyAdGuardStats)
		v1.POST("/adguard/test", controllers.TestAdGuardConnection)

		// Sonarr proxy routes
		v1.GET("/sonarr/:widget_id", controllers.ProxySonarrStats)
		v1.POST("/sonarr/test", controllers.TestSonarrConnection)

		// Radarr proxy routes
		v1.GET("/radarr/:widget_id", controllers.ProxyRadarrStats)
		v1.POST("/radarr/test", controllers.TestRadarrConnection)

		// Lidarr proxy routes
		v1.GET("/lidarr/:widget_id", controllers.ProxyLidarrStats)
		v1.POST("/lidarr/test", controllers.TestLidarrConnection)

		// Transmission proxy routes
		v1.GET("/transmission/:widget_id", controllers.ProxyTransmissionStats)
		v1.POST("/transmission/test", controllers.TestTransmissionConnection)

		// System stats route
		v1.GET("/system/stats", controllers.GetSystemStats)
	}

	return r
}
