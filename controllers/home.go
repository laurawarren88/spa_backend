package controllers

import (
	"context"
	"net/http"
	"spa_media_review/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HomeController struct {
	bookCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewHomeController(bookCollection, userCollection *mongo.Collection) *HomeController {
	return &HomeController{
		bookCollection: bookCollection,
		userCollection: userCollection,
	}
}

func (hc *HomeController) GetHome(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	var username string
	if exists {
		var user models.User
		err := hc.userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
		if err == nil {
			username = user.Username
		}
	}

	query := ctx.Query("q")
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := hc.bookCollection.Find(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer cursor.Close(context.TODO())

	var books []models.Book
	if err := cursor.All(context.TODO(), &books); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"books":    books,
		"username": username,
	})
}

func (hc *HomeController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	err = hc.userCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"_id":      user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
		"isAdmin":  user.IsAdmin,
	})
}
