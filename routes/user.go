package routes

import (
	"spa_media_review/controllers"
	"spa_media_review/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {
	userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("/register", uc.GetSignupForm)
		userRoutes.POST("/register", uc.SignupUser)
		userRoutes.GET("/login", uc.GetLoginForm)
		userRoutes.POST("/login", uc.LoginUser)
	}

	protected := router.Group("/api/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/forgot-password", uc.ForgotPassword)
		protected.POST("/reset-password", uc.ResetPassword)
		protected.POST("/logout", uc.LogoutUser)
	}
}
