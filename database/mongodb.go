package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var DB *mongo.Database
var BookCollection *mongo.Collection
var ReviewCollection *mongo.Collection

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

	fmt.Println("Connected to MongoDB.")
	return nil
}

func DisconnectDB() {
	if err := MongoClient.Disconnect(context.Background()); err != nil {
		fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
	}
}
