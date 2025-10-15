package main

import (
	"example/web-service-gin/database"
	"example/web-service-gin/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	// Initialize database connection
	database.Connect()

	// Run migrations
	database.Migrate()

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Get server port from environment variable, default to "8080"
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting server on port %s...\n", port)
	router.Run(fmt.Sprintf(":%s", port))
}
