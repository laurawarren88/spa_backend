package routes

import (
	"spa_media_review/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {

	userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("/login", uc.AddUser) // Add a user
		userRoutes.POST("/login", uc.LoginUser)
		userRoutes.GET("/register", uc.RegisterUser)
		userRoutes.POST("/register", uc.SignupUser)
		userRoutes.GET("/forgot-password", uc.ForgotPassword)
		userRoutes.GET("/logout", uc.LogoutUser)
	}
}
