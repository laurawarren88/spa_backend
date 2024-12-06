package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
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
	bookCollection   *mongo.Collection
	reviewCollection *mongo.Collection
}

func NewBookController(bookCollection, reviewCollection *mongo.Collection) *BookController {
	return &BookController{
		bookCollection:   bookCollection,
		reviewCollection: reviewCollection,
	}
}

func (bc *BookController) GetBooks(ctx *gin.Context) {
	query := ctx.Query("q")
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := bc.bookCollection.Find(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var books []models.Book
	if err := cursor.All(context.TODO(), &books); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
	maxSize := int64(20 << 20)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)

	if err := ctx.Request.ParseMultipartForm(maxSize); err != nil {
		if err.Error() == "http: request body too large" {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds 10MB limit"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	book := models.Book{
		ID:          primitive.NewObjectID(),
		Title:       ctx.PostForm("title"),
		Author:      ctx.PostForm("author"),
		Category:    ctx.PostForm("category"),
		Description: ctx.PostForm("description"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	file, err := ctx.FormFile("image")
	if err == nil {
		openFile, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer openFile.Close()

		fileBytes, err := io.ReadAll(openFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
			return
		}
		book.Image = base64.StdEncoding.EncodeToString(fileBytes)
	} else if err != http.ErrMissingFile {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file"})
		return
	}

	if errors := book.Validate(); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	_, err = bc.bookCollection.InsertOne(context.TODO(), book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		log.Println("Failed to create book:", err)
		return
	}

	log.Printf("Inserting book into collection: %s", bc.bookCollection.Name())
	ctx.JSON(http.StatusCreated, gin.H{"book": book})
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
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	fmt.Printf("Authenticated user ID: %v\n", userID)

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

func (bc *BookController) DeleteBookConfirmation(ctx *gin.Context) {
	fmt.Printf("Received DELETE confirmation request for ID: %s\n", ctx.Param("id"))
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
	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (bc *BookController) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	log.Printf("Received ID: %s", id)

	_, err = bc.reviewCollection.DeleteMany(context.TODO(), bson.M{"book._id": objectId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated reviews"})
		return
	}

	result, err := bc.bookCollection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	log.Printf("Delete result: %+v", result)
	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Book and associated reviews deleted successfully"})
}

func (bc *BookController) SearchBooks(ctx *gin.Context) {
	title := ctx.Query("title")
	author := ctx.Query("author")
	category := ctx.Query("category")

	filter := bson.M{}

	if title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if author != "" {
		filter["author"] = bson.M{"$regex": author, "$options": "i"}
	}
	if category != "" {
		filter["category"] = bson.M{"$regex": category, "$options": "i"}
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
