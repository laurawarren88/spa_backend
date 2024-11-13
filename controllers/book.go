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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookController struct {
	bookCollection *mongo.Collection
}

func NewBookController(collection *mongo.Collection) *BookController {
	return &BookController{bookCollection: collection}
}

// func (bc *BookController) GetBooks(ctx *gin.Context) {
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
// 	sortBy := ctx.DefaultQuery("sort", "title")
// 	order := ctx.DefaultQuery("order", "asc")

// 	skip := (page - 1) * limit

// 	opts := options.Find().
// 		SetSort(bson.D{{Key: sortBy, Value: getSortOrder(order)}}).
// 		SetSkip(int64(skip)).
// 		SetLimit(int64(limit))

// 	cursor, err := bc.bookCollection.Find(context.TODO(), bson.D{}, opts)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
// 		return
// 	}
// 	defer cursor.Close(context.TODO())

// 	var books []models.Book
// 	if err = cursor.All(context.TODO(), &books); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"books": books,
// 		"page":  page,
// 		"limit": limit,
// 		"total": len(books),
// 	})
// }

func (bc *BookController) GetBooks(c *gin.Context) {
	cursor, err := bc.bookCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var books []models.Book
	if err := cursor.All(context.TODO(), &books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}

func (bc *BookController) NewBook(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "New Book"})
}

func (bc *BookController) GetBookByID(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var book models.Book
	if err := bc.bookCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&book); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func (bc *BookController) CreateBook(ctx *gin.Context) {
	var book models.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = primitive.NewObjectID()
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	if errors := book.Validate(); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	_, err := bc.bookCollection.InsertOne(context.TODO(), book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (bc *BookController) UpdateBook(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "Book updated successfully"})
}

func (bc *BookController) EditedBook(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updateBook models.Book
	if err := ctx.ShouldBindJSON(&updateBook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":       updateBook.Title,
			"author":      updateBook.Author,
			"category":    updateBook.Category,
			"description": updateBook.Description,
			"updated_at":  time.Now(),
		},
	}

	result := bc.bookCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": objectId},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

func (bc *BookController) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := bc.bookCollection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (bc *BookController) SearchBooks(ctx *gin.Context) {
	query := ctx.Query("q")
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"author": bson.M{"$regex": query, "$options": "i"}},
			{"category": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := bc.bookCollection.Find(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search books"})
		return
	}
	defer cursor.Close(context.TODO())

	var books []models.Book
	if err = cursor.All(context.TODO(), &books); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"books": books})
}

func getSortOrder(order string) int {
	if order == "desc" {
		return -1
	}
	return 1
}
