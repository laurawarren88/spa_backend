package controllers

import (
	"context"
	"fmt"
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
	userCollection   *mongo.Collection
}

func NewReviewController(reviewCollection, bookCollection, userCollection *mongo.Collection) *ReviewController {
	return &ReviewController{
		reviewCollection: reviewCollection,
		bookCollection:   bookCollection,
		userCollection:   userCollection,
	}
}

func (rc *ReviewController) GetReviews(ctx *gin.Context) {
	var reviews []models.Review
	cursor, err := rc.reviewCollection.Aggregate(context.TODO(), []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user_info",
			},
		},
		{
			"$addFields": bson.M{
				"username": bson.M{"$arrayElemAt": []interface{}{"$user_info.username", 0}},
			},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	if err := cursor.All(context.TODO(), &reviews); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reviews"})
		return
	}

	for i, review := range reviews {
		var user models.User
		err := rc.userCollection.FindOne(context.TODO(), bson.M{"_id": review.UserID}).Decode(&user)
		if err == nil {
			reviews[i].Username = user.Username
		} else {
			reviews[i].Username = "Unknown User"
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (rc *ReviewController) NewReview(ctx *gin.Context) {
	bookID := ctx.Param("bookId")

	objID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := rc.bookCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&book); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))
	var user models.User
	err = rc.userCollection.FindOne(context.TODO(), bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"book": book,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"username": user.Username,
		},
	})
}

func (rc *ReviewController) CreateReview(ctx *gin.Context) {
	var input struct {
		BookID string `json:"book_id" binding:"required"`
		Review string `json:"review"`
		Rating int    `json:"rating" binding:"required,min=1,max=5"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

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
	if err := rc.userCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	newReview := models.Review{
		ID:        primitive.NewObjectID(),
		UserID:    objectID,
		Username:  user.Username,
		Review:    input.Review,
		Rating:    input.Rating,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Book:      book,
		User:      user,
	}

	_, err = rc.reviewCollection.InsertOne(context.TODO(), newReview)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		log.Println("Failed to create review:", err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Review created",
		"review":  newReview,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"username": user.Username,
		},
	})
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
	log.Printf("Querying reviews with bookId: %s", objID.Hex())

	pipeline := []bson.M{
		{"$match": bson.M{"book._id": objID}},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user_info",
			},
		},
		{
			"$addFields": bson.M{
				"username": bson.M{"$arrayElemAt": []interface{}{"$user_info.username", 0}},
			},
		},
	}

	cursor, err := rc.reviewCollection.Aggregate(context.TODO(), pipeline)
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
		log.Printf("Update failed for review ID %s: %v", id, result.Err())
		ctx.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	updatedReview := models.Review{}
	if err := result.Decode(&updatedReview); err != nil {
		log.Printf("Failed to decode updated review: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

func (rc *ReviewController) DeleteReviewConfirmation(ctx *gin.Context) {
	fmt.Printf("Received DELETE confirmation request for ID: %s\n", ctx.Param("id"))
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
	ctx.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

func (rc *ReviewController) DeleteReview(ctx *gin.Context) {
	fmt.Printf("Received DELETE request for ID: %s\n", ctx.Param("id"))
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	fmt.Printf("Attempting to delete review with ID: %s\n", id)
	log.Printf("Received ID: %s", id)

	result, err := rc.reviewCollection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		fmt.Println("Error during deletion:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	log.Printf("Delete result: %+v", result)
	if result.DeletedCount == 0 {
		fmt.Println("Review not found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	fmt.Printf("Delete result: %+v\n", result)
	fmt.Printf("Error: %v\n", err)
	ctx.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
