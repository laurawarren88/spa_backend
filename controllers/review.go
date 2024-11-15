package controllers

import (
	"context"
	"net/http"
	"spa_media_review/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewController struct {
	reviewCollection *mongo.Collection
}

func NewReviewController(collection *mongo.Collection) *ReviewController {
	return &ReviewController{reviewCollection: collection}
}

func (rc *ReviewController) GetReviews(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	skip := (page - 1) * limit

	options := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := rc.reviewCollection.Find(context.TODO(), bson.M{}, options)
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

	ctx.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (rc *ReviewController) NewReview(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "New Review"})
}

func (rc *ReviewController) GetReviewsByBookID(ctx *gin.Context) {
	bookID := ctx.Param("bookId")
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID format"})
		return
	}

	cursor, err := rc.reviewCollection.Find(context.TODO(), bson.M{"book_id": objectID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer cursor.Close(context.TODO())

	var reviews []models.Review
	if err := cursor.All(context.TODO(), &reviews); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode reviews"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (rc *ReviewController) CreateReview(ctx *gin.Context) {
	var review models.Review
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate book ID
	if _, err := rc.reviewCollection.FindOne(context.TODO(), bson.M{"_id": review.BookID}).DecodeBytes(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID"})
		return
	}

	review.ID = primitive.NewObjectID()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	result, err := rc.reviewCollection.InsertOne(context.TODO(), review)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}
