package routes

import (
	"spa_media_review/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterHomeRoute(router *gin.Engine, hc *controllers.HomeController) {
	router.GET("/api/", hc.GetHome)
}
