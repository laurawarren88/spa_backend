package config

import (
	"log"
	"spa_media_review/controllers"
	"spa_media_review/middleware"
	"spa_media_review/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func SetGinMode() {
	gin.SetMode(gin.ReleaseMode)
}

func SetupServer() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	return router
}

func SetupHandlers(router *gin.Engine, bookCollection *mongo.Collection, reviewCollection *mongo.Collection) {
	homeController := controllers.NewHomeController(bookCollection)
	bookController := controllers.NewBookController(bookCollection)
	reviewController := controllers.NewReviewController(reviewCollection)

	routes.RegisterHomeRoute(router, homeController)
	routes.RegisterBookRoutes(router, bookController)
	routes.RegisterReviewRoutes(router, reviewController)
}
