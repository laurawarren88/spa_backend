package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// if ctx.Request.Method == "OPTIONS" {
		// 	ctx.Next()
		// 	return
		// }
		isAdmin, exists := ctx.Get("isAdmin")
		if !exists {
			// The "isAdmin" key is not set in the context, indicating an issue with authentication middleware
			fmt.Println("isAdmin context key missing")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			ctx.Abort()
			return
		}

		// Ensure the value of isAdmin is a boolean
		isAdminBool, ok := isAdmin.(bool)
		if !ok || !isAdminBool {
			fmt.Println("Admin access denied")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			ctx.Abort()
			return
		}

		fmt.Println("Admin access granted")
		ctx.Next()
	}
}
