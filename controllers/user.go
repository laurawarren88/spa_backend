package controllers

import (
	"net/http"
	"os"
	"spa_media_review/models"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)

	ctx.SetCookie(
		"token",     // Name of the token cookie
		tokenString, // Token string
		3600*24,     // Expiration time (1 day)
		"/",         // Path
		"localhost", // Domain
		false,       // Secure
		false,       // HTTP only
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    gin.H{"email": user.Email, "name": user.Username}})
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

	ctx.SetCookie(
		"token",     // Name of the token cookie
		"",          // Empty the value
		-1,          // Expires immediately
		"/",         // path
		"localhost", // Domain
		false,       // Secure
		false,       // HTTP only
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
	ctx.SetCookie(
		"token",     // Name of the token cookie
		"",          // Empty the value
		-1,          // Expires immediately
		"/",         // path
		"localhost", // Domain
		false,       // HTTPS only
		false,       // HTTP only
	)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
