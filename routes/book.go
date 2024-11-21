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

	protected := router.Group("/api/books")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/", bc.CreateBook)
		protected.GET("/new", bc.NewBook)
		protected.GET("/edit/:id", bc.UpdateBook)
		protected.PUT("/edit/:id", bc.EditedBook)
		protected.DELETE("/delete/:id", bc.DeleteBook)
	}
}
