package config

import (
	"log"
	"os"
	"spa_media_review/controllers"
	"spa_media_review/middleware"
	"spa_media_review/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func SetGinMode() {
	gin.SetMode(gin.ReleaseMode)
}

func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func SetupServer() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	return router
}

func SetupHandlers(router *gin.Engine, bookCollection *mongo.Collection, reviewCollection *mongo.Collection, userCollection *mongo.Collection) {
	homeController := controllers.NewHomeController(bookCollection, userCollection)
	bookController := controllers.NewBookController(bookCollection, reviewCollection)
	reviewController := controllers.NewReviewController(reviewCollection, bookCollection, userCollection)
	userController := controllers.NewUserController(userCollection)

	routes.RegisterHomeRoute(router, homeController)
	routes.RegisterBookRoutes(router, bookController)
	routes.RegisterReviewRoutes(router, reviewController)
	routes.RegisterUserRoutes(router, userController)
}
