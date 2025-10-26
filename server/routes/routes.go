package routes

import (
	"dashboard-server/controllers"
	"dashboard-server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("", controllers.GetDashboards)
			dashboards.POST("", controllers.CreateDashboard)
			dashboards.GET("/:id", controllers.GetDashboard)
			dashboards.PUT("/:id", controllers.UpdateDashboard)
			dashboards.DELETE("/:id", controllers.DeleteDashboard)

			dashboards.GET("/:id/widgets", controllers.GetWidgets)
			dashboards.POST("/:id/widgets", controllers.CreateWidget)
		}

		widgets := v1.Group("/widgets")
		{
			widgets.GET("/:id", controllers.GetWidget)
			widgets.PUT("/:id", controllers.UpdateWidget)
			widgets.PUT("/:id/state", controllers.UpdateWidgetState)
			widgets.DELETE("/:id", controllers.DeleteWidget)
		}

		v1.GET("/adguard/:widget_id", controllers.ProxyAdGuardStats)
		v1.POST("/adguard/test", controllers.TestAdGuardConnection)

		v1.GET("/sonarr/:widget_id", controllers.ProxySonarrStats)
		v1.POST("/sonarr/test", controllers.TestSonarrConnection)

		v1.GET("/radarr/:widget_id", controllers.ProxyRadarrStats)
		v1.POST("/radarr/test", controllers.TestRadarrConnection)

		v1.GET("/lidarr/:widget_id", controllers.ProxyLidarrStats)
		v1.POST("/lidarr/test", controllers.TestLidarrConnection)

		v1.GET("/transmission/:widget_id", controllers.ProxyTransmissionStats)
		v1.POST("/transmission/test", controllers.TestTransmissionConnection)

		v1.GET("/system/stats", controllers.GetSystemStats)
	}

	return r
}
