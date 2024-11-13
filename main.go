package main

import (
	"fmt"
	"log"

	"spa_media_review/config"
	"spa_media_review/database"
)

func init() {
	config.LoadEnvVariables()

	if err := database.Connect_to_mongodb(); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	config.SetGinMode()
}

func main() {
	defer database.DisconnectDB()

	router := config.SetupServer()

	config.SetupHandlers(router, database.BookCollection)

	fmt.Printf("Starting the server\n")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
