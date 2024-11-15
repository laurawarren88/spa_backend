package controllers

import (
	"context"
	"net/http"
	"spa_media_review/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HomeController struct {
	bookCollection *mongo.Collection
}

func NewHomeController(collection *mongo.Collection) *HomeController {
	return &HomeController{bookCollection: collection}
}

func (hc *HomeController) GetHome(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, gin.H{"books": books})
}
