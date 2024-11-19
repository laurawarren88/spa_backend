package controllers

import (
	"context"
	"log"
	"net/http"
	"spa_media_review/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewController struct {
	reviewCollection *mongo.Collection
	bookCollection   *mongo.Collection
}

func NewReviewController(reviewCollection, bookCollection *mongo.Collection) *ReviewController {
	return &ReviewController{
		reviewCollection: reviewCollection,
		bookCollection:   bookCollection,
	}
}

func (rc *ReviewController) GetReviews(ctx *gin.Context) {
	var reviews []models.Review
	cursor, err := rc.reviewCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &reviews); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reviews"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (rc *ReviewController) NewReview(ctx *gin.Context) {
	bookID := ctx.Param("bookId")

	// Convert bookID to ObjectID
	objID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Fetch the book from the database
	var book models.Book
	if err := rc.bookCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&book); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Return the book data
	ctx.JSON(http.StatusOK, gin.H{"book": book})
}

// Create a new review
func (rc *ReviewController) CreateReview(ctx *gin.Context) {
	var input struct {
		BookID string `json:"book_id" binding:"required"`
		Review string `json:"review" binding:"required"`
		Rating int    `json:"rating" binding:"required,min=1,max=5"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the book data using the provided book ID
	bookID, err := primitive.ObjectIDFromHex(input.BookID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := rc.bookCollection.FindOne(context.TODO(), bson.M{"_id": bookID}).Decode(&book); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Create the review object
	newReview := models.Review{
		ID:        primitive.NewObjectID(),
		Review:    input.Review,
		Rating:    input.Rating,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Book:      book,
	}

	// Insert the review into the database
	_, err = rc.reviewCollection.InsertOne(context.TODO(), newReview)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		log.Println("Failed to create review:", err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Review created", "review": newReview})
}

func (rc *ReviewController) GetReviewsByBookID(c *gin.Context) {
	bookID := c.Param("bookId")
	log.Printf("Received request for book ID: %s", bookID)

	objID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		log.Printf("Invalid book ID conversion: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	err = rc.bookCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		log.Printf("Book find error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	log.Printf("Found book: %s", book.Title)

	// Log the query we're about to make
	log.Printf("Querying reviews with bookId: %s", objID.Hex())

	// cursor, err := rc.reviewCollection.Find(context.TODO(), bson.M{"book_id": objID})
	// Try this query instead
	cursor, err := rc.reviewCollection.Find(context.TODO(), bson.M{"book._id": bson.M{"$eq": objID}})

	// Add a debug query to see what we're matching against
	log.Printf("Query filter: %+v", bson.M{"book._id": bson.M{"$eq": objID}})

	if err != nil {
		log.Printf("Review query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	var reviews []models.Review
	if err := cursor.All(context.TODO(), &reviews); err != nil {
		log.Printf("Cursor parsing error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reviews"})
		return
	}
	log.Printf("Found %d reviews", len(reviews))
	log.Printf("Reviews to be sent: %+v", reviews)

	c.JSON(http.StatusOK, gin.H{
		"bookTitle": book.Title,
		"reviews":   reviews,
	})
}

// Get a review by its ID
func (rc *ReviewController) GetReviewByID(ctx *gin.Context) {
	reviewID := ctx.Param("id")

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var review models.Review
	if err := rc.reviewCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&review); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	ctx.JSON(http.StatusOK, review)
}

func (rc *ReviewController) UpdateReview(ctx *gin.Context) {
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

func (rc *ReviewController) EditedReview(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updateReview models.Review
	if err := ctx.ShouldBindJSON(&updateReview); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"rating":     updateReview.Rating,
			"review":     updateReview.Review,
			"updated_at": time.Now(),
		},
	}

	result := rc.reviewCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": objectId},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

func (rc *ReviewController) DeleteReview(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := rc.reviewCollection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
