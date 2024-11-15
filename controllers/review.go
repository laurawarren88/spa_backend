package controllers

import (
	"context"
	"net/http"
	"spa_media_review/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewController struct {
	reviewCollection *mongo.Collection
}

func NewReviewController(collection *mongo.Collection) *ReviewController {
	return &ReviewController{reviewCollection: collection}
}

func (rc *ReviewController) GetReviews(ctx *gin.Context) {
	cursor, err := rc.reviewCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var reviews []models.Review
	if err := cursor.All(context.TODO(), &reviews); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode reviews"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}

func (rc *ReviewController) NewReview(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "New Review"})
}

func (rc *ReviewController) GetReviewByID(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var review models.Review
	if err := rc.reviewCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&review); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	ctx.JSON(http.StatusOK, review)
}

func (rc *ReviewController) CreateReview(ctx *gin.Context) {
	var review models.Review
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	review.ID = primitive.NewObjectID()

	result, err := rc.reviewCollection.InsertOne(context.TODO(), review)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}
