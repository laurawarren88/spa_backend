package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// log.Println("AuthMiddleware invoked")

// 		tokenString, err := ctx.Cookie("token")
// 		if err != nil || tokenString == "" {
// 			fmt.Println("No token found in cookies")
// 			tokenString = strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
// 		}
// 		authHeader := ctx.GetHeader("Authorization")
// 		fmt.Printf("Token from cookie: %s\n", tokenString)
// 		fmt.Printf("Authorization Header: %s\n", authHeader)

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return []byte(os.Getenv("SECRET_KEY")), nil
// 		})

// 		if err != nil || !token.Valid {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			ctx.Abort()
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			userID, ok := claims["sub"].(string)
// 			if !ok || userID == "" {
// 				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user claims"})
// 				ctx.Abort()
// 				return
// 			}

// 			isAdmin := claims["isAdmin"].(bool)
// 			ctx.Set("userID", userID)
// 			ctx.Set("isAdmin", isAdmin)

// 			// ** Uncomment for Debugging **
// 			// fmt.Printf("Token parsed: %v\n", claims)
// 			// fmt.Printf("Token from cookie or header: %s\n", tokenString)
// 			// fmt.Printf("Fetching user with ID: %v\n", userID)
// 			// fmt.Printf("Error fetching user: %v\n", err)

// 			ctx.Next()
// 		} else {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
// 			ctx.Abort()
// 		}
// 	}
// }

// Define custom claims for better type safety
type Claims struct {
	UserID  string `json:"sub"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenString string

		// First, check for the token in cookies
		if cookieToken, err := ctx.Cookie("token"); err == nil {
			tokenString = cookieToken
		} else {
			// Fall back to checking the Authorization header
			authHeader := ctx.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			ctx.Abort()
			return
		}

		// Parse and validate the JWT
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Extract claims and attach to context
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx.Set("userID", claims.UserID)
			ctx.Set("isAdmin", claims.IsAdmin)
			ctx.Next()
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
		}
	}
}
