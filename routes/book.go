package routes

import (
	"spa_media_review/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(router *gin.Engine, bc *controllers.BookController) {

	bookRoutes := router.Group("/api/books")
	{
		bookRoutes.GET("/", bc.GetBooks)
		// bookRoutes.POST("/", bc.CreateBook)
		bookRoutes.POST("/", bc.CreateBookWithImage)
		bookRoutes.GET("/new", bc.NewBook)
		bookRoutes.GET("/search", bc.SearchBooks)
		bookRoutes.GET("/edit/:id", bc.UpdateBook)
		bookRoutes.PUT("/edit/:id", bc.EditedBook)
		bookRoutes.DELETE("/delete/:id", bc.DeleteBook)
		bookRoutes.GET("/:id", bc.GetBookByID)
	}
}
