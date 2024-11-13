package routes

import (
	"spa_media_review/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(router *gin.Engine, bc *controllers.BookController) {

	bookRoutes := router.Group("/api/books")
	{
		bookRoutes.GET("/", bc.GetBooks)
		bookRoutes.GET("/search", bc.SearchBooks)
		bookRoutes.GET("/new", bc.NewBook)
		bookRoutes.POST("/", bc.CreateBook)
		bookRoutes.GET("/:id", bc.GetBookByID)
		bookRoutes.GET("/edit/:id", bc.UpdateBook)
		bookRoutes.PUT("/edit/:id", bc.EditedBook)
		bookRoutes.DELETE("/delete/:id", bc.DeleteBook)
	}
}
