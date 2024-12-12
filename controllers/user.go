package controllers

import (
	"log"
	"net/http"
	"os"
	"spa_media_review/middleware"
	"spa_media_review/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userCollection *mongo.Collection
}

func NewUserController(collection *mongo.Collection) *UserController {
	return &UserController{userCollection: collection}
}

func (uc *UserController) GetSignupForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Registration form"})
}

func (uc *UserController) SignupUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	count, err := uc.userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hash)

	_, err = uc.userCollection.InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"_id":      user.ID.Hex(),
			"email":    user.Email,
			"username": user.Username,
			"isAdmin":  user.IsAdmin,
		},
	})
}

func (uc *UserController) GetLoginForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Login form", "user": nil})
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var user models.User
	err := uc.userCollection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	domain, secure, httpOnly, err := middleware.GetCookieSettings()
	if err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"access_token",
		accessToken,
		3600*24,
		"/",
		domain,
		secure,
		httpOnly,
	)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		3600*24,
		"/",
		domain,
		secure,
		httpOnly,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"_id":      user.ID.Hex(),
			"email":    user.Email,
			"username": user.Username,
			"isAdmin":  user.IsAdmin,
		},
	})
}

func (uc *UserController) ForgotPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Reset password form"})
}

func (uc *UserController) ResetPassword(ctx *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	env := os.Getenv("ENV")

	var domain string
	if env == "development" {
		domain = strings.Split(os.Getenv("DEV_ALLOWED_ORIGIN"), "//")[1]
	} else {
		domain = strings.Split(os.Getenv("PROD_ALLOWED_ORIGIN"), "//")[1]
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		domain,
		false,
		false,
	)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		domain,
		false,
		false,
	)

	result := uc.userCollection.FindOneAndUpdate(
		ctx,
		bson.M{"email": input.Email},
		bson.M{"$set": bson.M{"passwordResetToken": "token"}},
	)

	if result.Err() != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent"})
}

func (uc *UserController) LogoutUser(ctx *gin.Context) {
	log.Println("LogoutUser endpoint hit")

	env := os.Getenv("ENV")

	var domain string
	if env == "development" {
		domain = strings.Split(os.Getenv("DEV_ALLOWED_ORIGIN"), "//")[1]
	} else {
		domain = strings.Split(os.Getenv("PROD_ALLOWED_ORIGIN"), "//")[1]
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		domain,
		false,
		false,
	)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		domain,
		false,
		false,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
