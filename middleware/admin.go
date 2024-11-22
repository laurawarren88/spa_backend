package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func RequireAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user, exists := c.Get("user")
// 		if !exists {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			c.Abort()
// 			return
// 		}

// 		if !user.(models.User).IsAdmin {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdmin, exists := ctx.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
