package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context)
// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
// c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow requests from your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
