package main

import (
	"fmt"
	"log"

	"spa_media_review/config"
	"spa_media_review/database"
)

func init() {
	config.LoadEnv()

	if err := database.Connect_to_mongodb(); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	if err := database.SetupAdminUser(database.DB); err != nil {
		log.Printf("Admin user setup: %v", err)
	}

	config.SetGinMode()
}

func main() {
	defer database.DisconnectDB()

	router := config.SetupServer()

	config.SetupHandlers(router, database.BookCollection, database.ReviewCollection, database.UserCollection)

	fmt.Printf("Starting the server\n")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
