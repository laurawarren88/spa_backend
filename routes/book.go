package routes

import (
	"spa_media_review/controllers"
	"spa_media_review/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(router *gin.Engine, bc *controllers.BookController) {

	bookRoutes := router.Group("/api/books")
	{
		bookRoutes.GET("/", bc.GetBooks)
		bookRoutes.GET("/search", bc.SearchBooks)
		bookRoutes.GET("/:id", bc.GetBookByID)
	}

	adminRoutes := router.Group("/api/books")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireAdmin())
	{
		adminRoutes.POST("/", bc.CreateBook)
		adminRoutes.GET("/new", bc.NewBook)
		adminRoutes.GET("/edit/:id", bc.UpdateBook)
		adminRoutes.PUT("/edit/:id", bc.EditedBook)
		adminRoutes.DELETE("/delete/:id", bc.DeleteBook)
	}
}
