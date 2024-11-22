package database

import (
	"context"
	"fmt"
	"os"
	"spa_media_review/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

var MongoClient *mongo.Client
var DB *mongo.Database
var BookCollection *mongo.Collection
var ReviewCollection *mongo.Collection
var UserCollection *mongo.Collection

func Connect_to_mongodb() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	MongoClient = client
	DB = client.Database("media_review_app")
	BookCollection = DB.Collection("books")
	ReviewCollection = DB.Collection("reviews")
	UserCollection = DB.Collection("users")

	fmt.Println("Connected to MongoDB.")
	return nil
}

func DisconnectDB() {
	if err := MongoClient.Disconnect(context.Background()); err != nil {
		fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
	}
}

func SetupAdminUser(db *mongo.Database) error {
	collection := db.Collection("users")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)

	admin := models.User{
		ID:        primitive.NewObjectID(),
		Username:  "admin",
		Email:     "admin@admin.com",
		Password:  string(hashedPassword),
		IsAdmin:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"email": admin.Email},
		bson.M{"$set": admin},
		options.Update().SetUpsert(true),
	)

	return err
}
