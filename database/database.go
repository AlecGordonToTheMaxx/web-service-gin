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

// Database holds the database connection and underlying SQL DB
type Database struct {
	DB    *gorm.DB
	sqlDB *sql.DB
}

// Connect initializes the database connection
func Connect() (*Database, error) {
	// Get database name from environment variable, default to "albums.db"
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "albums.db"
	}

	// Open database with pure Go SQLite driver
	sqlDB, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
	}

	// Use GORM with the existing connection
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		sqlDB.Close()
		return nil, err
	}

	log.Printf("Database connection established: %s\n", dbName)

	return &Database{
		DB:    gormDB,
		sqlDB: sqlDB,
	}, nil
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	err := d.DB.AutoMigrate(&models.Album{})
	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.sqlDB != nil {
		log.Println("Closing database connection...")
		return d.sqlDB.Close()
	}
	return nil
}
