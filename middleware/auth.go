package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"spa_media_review/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserID   string `json:"sub"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		if cookieToken, err := ctx.Cookie("access_token"); err == nil {
			accessToken = cookieToken
		} else {
			authHeader := ctx.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				accessToken = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if accessToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not provided"})
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			if err.Error() == "Token is expired" {
				refreshToken, err := ctx.Cookie("refresh_token")
				if err != nil {
					ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not provided"})
					ctx.Abort()
					return
				}

				_, err = jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(os.Getenv("SECRET_KEY")), nil
				})
				if err != nil {
					ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
					ctx.Abort()
					return
				}

				var user models.User
				accessToken, err := GenerateToken(user)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
					ctx.Abort()
					return
				}

				domain, secure, httpOnly, err := GetCookieSettings()
				if err != nil {
					log.Fatalf("Failed to parse environment variables: %v", err)
				}

				ctx.SetCookie("access_token", accessToken, 3600*24, "/", domain, secure, httpOnly)

				ctx.Set("userID", user.ID.Hex())
				ctx.Set("isAdmin", user.IsAdmin)
				ctx.Next()
				return
			}

			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

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
