package routes

import (
	"spa_media_review/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(router *gin.Engine, rc *controllers.ReviewController) {

	reviewRoutes := router.Group("/api/reviews")
	{
		reviewRoutes.GET("/", rc.GetReviews)
		reviewRoutes.POST("/new", rc.CreateReview)
		// reviewRoutes.GET("/new", rc.NewReview)
		reviewRoutes.GET("book/:book_id", rc.NewReview)
		reviewRoutes.GET("/:id", rc.GetReviewByBookID)
	}
}
