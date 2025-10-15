package database

import (
	"database/sql"
	"example/web-service-gin/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // Pure Go SQLite driver (no CGO required)
)

var DB *gorm.DB

// Connect initializes the database connection
func Connect() {
	var err error

	// Get database name from environment variable, default to "albums.db"
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "albums.db"
	}

	// Open database with pure Go SQLite driver
	sqlDB, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Use GORM with the existing connection
	DB, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Printf("Database connection established: %s\n", dbName)
}

// Migrate runs database migrations
func Migrate() {
	err := DB.AutoMigrate(&models.Album{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}
