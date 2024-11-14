package controllers

import (
	"context"
	"encoding/base64"
	"io"
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

type BookController struct {
	bookCollection *mongo.Collection
}

func NewBookController(collection *mongo.Collection) *BookController {
	return &BookController{bookCollection: collection}
}

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

// func (bc *BookController) CreateBook(ctx *gin.Context) {
// 	fmt.Println("Received create book request")
// 	var book models.Book
// 	if err := ctx.ShouldBindJSON(&book); err != nil {
// 		fmt.Println("Binding error:", err)
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	book.ID = primitive.NewObjectID()
// 	book.CreatedAt = time.Now()
// 	book.UpdatedAt = time.Now()

// 	if errors := book.Validate(); len(errors) > 0 {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
// 		return
// 	}

// 	_, err := bc.bookCollection.InsertOne(context.TODO(), book)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, book)
// }

func (bc *BookController) CreateBookWithImage(c *gin.Context) {
	// Set a reasonable max size for the multipart form
	maxSize := int64(10 << 20) // 10MB
	c.Request.ParseMultipartForm(maxSize)

	// Get form values
	book := models.Book{
		ID:          primitive.NewObjectID(),
		Title:       c.PostForm("title"),
		Author:      c.PostForm("author"),
		Category:    c.PostForm("category"),
		Description: c.PostForm("description"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Handle image file
	file, err := c.FormFile("image")
	if err == nil {
		openFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer openFile.Close()

		// Read and encode file contents
		fileBytes, err := io.ReadAll(openFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
			return
		}
		// Encode to Base64
		book.Image = base64.StdEncoding.EncodeToString(fileBytes)
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file"})
		return
	}

	// Validate and save
	if errors := book.Validate(); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	_, err = bc.bookCollection.InsertOne(context.TODO(), book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		log.Println("Failed to create book:", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"book": book})
}

func (bc *BookController) UpdateBook(ctx *gin.Context) {
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
