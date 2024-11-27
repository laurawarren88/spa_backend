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
		userRoutes.GET("/forgot_password", uc.ForgotPassword)
		userRoutes.POST("/forgot_password", uc.ResetPassword)
	}

	protected := router.Group("/api/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", uc.LogoutUser)
	}
}
