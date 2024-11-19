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
		reviewRoutes.GET("/edit/:id", rc.UpdateReview)           // Update a Review
		reviewRoutes.PUT("/edit/:id", rc.EditedReview)           // Update a Review
		reviewRoutes.DELETE("/delete/:id", rc.DeleteReview)      // Delete a Review
		reviewRoutes.GET("/:id", rc.GetReviewByID)               // Fetch a single Review
	}
}
