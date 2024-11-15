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

func (rc *ReviewController) GetReviews(c *gin.Context) {
	var reviews []models.Review
	cursor, err := rc.reviewCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// Get reviews by book ID
func (rc *ReviewController) GetReviewsByBookID(c *gin.Context) {
	bookID := c.Param("bookId")
	var reviews []models.Review

	filter := bson.M{"bookId": bookID}
	cursor, err := rc.reviewCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// Get a review by its ID
func (rc *ReviewController) GetReviewByID(c *gin.Context) {
	reviewID := c.Param("id")
	var review models.Review

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	err = rc.reviewCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review})
}

// Create a new review
func (rc *ReviewController) CreateReview(c *gin.Context) {
	var newReview models.Review

	if err := c.BindJSON(&newReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newReview.ID = primitive.NewObjectID()
	newReview.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := rc.reviewCollection.InsertOne(context.TODO(), newReview)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review created", "review": newReview})
}

// Handle "new review" logic
func (rc *ReviewController) NewReview(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "New Review"})
}
