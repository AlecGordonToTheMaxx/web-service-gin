package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web-service-gin/backend/internal/album"
	"web-service-gin/backend/internal/chat"
	"web-service-gin/backend/internal/middleware"
	"web-service-gin/backend/internal/platform/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	// Create context for database initialization
	ctx := context.Background()

	// Initialize database connection
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(ctx); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize album domain
	albumRepo := album.NewRepository(db.Pool)
	albumHandler := album.NewHandler(albumRepo)

	// Initialize chat domain
	chatService := chat.NewService(albumRepo)
	chatHandler := chat.NewHandler(chatService)

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORS())

	// Setup routes
	albumGroup := router.Group("/albums")
	albumHandler.RegisterRoutes(albumGroup)

	chatGroup := router.Group("/chat")
	chatHandler.RegisterRoutes(chatGroup)

	// Get server port from environment variable, default to "8080"
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
