package routes

import (
	"spa_media_review/controllers"
	"spa_media_review/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterHomeRoute(router *gin.Engine, hc *controllers.HomeController) {

	homeRoutes := router.Group("/api")
	{
		homeRoutes.GET("/", hc.GetHome)
	}

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile/:userId", hc.GetProfile)
	}
}
