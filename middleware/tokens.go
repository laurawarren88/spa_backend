package middleware

import (
	"fmt"
	"os"
	"spa_media_review/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	accessSecret := os.Getenv("ACCESS_SECRET_KEY")
	if accessSecret == "" {
		return "", fmt.Errorf("access secret key not set in environment")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessSecret))
}

func GenerateRefreshToken(user models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}

	refreshSecret := os.Getenv("REFRESH_SECRET_KEY")
	if refreshSecret == "" {
		return "", fmt.Errorf("refresh secret key not set in environment")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecret))
}
