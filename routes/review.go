package routes

import (
	"spa_media_review/controllers"
	"spa_media_review/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(router *gin.Engine, rc *controllers.ReviewController) {
	reviewRoutes := router.Group("/api/reviews")
	{
		reviewRoutes.GET("/", rc.GetReviews)
		reviewRoutes.GET("/book/:bookId", rc.GetReviewsByBookID)
		reviewRoutes.GET("/:id", rc.GetReviewByID)
	}
	protected := router.Group("/api/reviews")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/", rc.CreateReview)
		protected.GET("/new/:bookId", rc.NewReview)
	}
	adminRoutes := router.Group("/api/reviews")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireAdmin())
	{
		adminRoutes.GET("/edit/:id", rc.UpdateReview)
		adminRoutes.PUT("/edit/:id", rc.EditedReview)
		adminRoutes.GET("/delete/:id", rc.DeleteReviewConfirmation)
		adminRoutes.DELETE("/delete/:id", rc.DeleteReview)
	}
}
