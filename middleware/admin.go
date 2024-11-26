package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.Next() // Let preflight requests pass
			return
		}

		isAdmin, exists := ctx.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			fmt.Println("Admin access denied")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			ctx.Abort()
			return
		}
		fmt.Println("Admin access granted")
		ctx.Next()
	}
}
