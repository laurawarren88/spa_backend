package routes

import (
	"spa_media_review/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(router *gin.Engine, rc *controllers.ReviewController) {

	reviewRoutes := router.Group("/api/reviews")
	{
		reviewRoutes.GET("/", rc.GetReviews) // Fetch All Reviews
		reviewRoutes.POST("/", rc.CreateReview)
		reviewRoutes.GET("/new/:bookId", rc.NewReview)           // Create a new Review
		reviewRoutes.GET("/book/:bookId", rc.GetReviewsByBookID) // Fetch Reviews by Book ID
		reviewRoutes.GET("/:id", rc.GetReviewByID)               // Fetch a single Review
	}
}
