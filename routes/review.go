package routes

import (
	"spa_media_review/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(router *gin.Engine, rc *controllers.ReviewController) {

	reviewRoutes := router.Group("/api/reviews")
	{
		reviewRoutes.GET("/", rc.GetReviews)
		reviewRoutes.POST("/", rc.CreateReview)
		// reviewRoutes.GET("/new", rc.NewReview)
		reviewRoutes.GET("book/:bookId", rc.NewReview)
		reviewRoutes.GET("/:id", rc.GetReviewByID)
	}
}
