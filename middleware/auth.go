package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// log.Println("AuthMiddleware invoked")

		tokenString, err := ctx.Cookie("token")
		if err != nil || tokenString == "" {
			tokenString = strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user claims"})
				ctx.Abort()
				return
			}

			isAdmin := claims["isAdmin"].(bool)
			ctx.Set("userID", userID)
			ctx.Set("isAdmin", isAdmin)

			// ** Uncomment for Debugging **
			// fmt.Printf("Token parsed: %v\n", claims)
			// fmt.Printf("Token from cookie or header: %s\n", tokenString)
			// fmt.Printf("Fetching user with ID: %v\n", userID)
			// fmt.Printf("Error fetching user: %v\n", err)

			ctx.Next()
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
		}
	}
}