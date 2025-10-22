package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database holds the database connection pool
type Database struct {
	Pool *pgxpool.Pool
}

// Connect initializes the PostgreSQL database connection pool
func Connect(ctx context.Context) (*Database, error) {
	// Get database connection string from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Build connection string from individual components
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5432")
		user := getEnvOrDefault("DB_USER", "postgres")
		password := getEnvOrDefault("DB_PASSWORD", "postgres")
		dbName := getEnvOrDefault("DB_NAME", "albums")
		sslMode := getEnvOrDefault("DB_SSLMODE", "disable")

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbName, sslMode)
	}

	// Create connection pool
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection established")

	return &Database{Pool: pool}, nil
}

// migration represents a single database migration
type migration struct {
	version     int
	description string
	up          string
}

// migrations contains all database migrations in order
//
// To add a new migration:
// 1. Add a new migration struct to this array
// 2. Increment the version number (must be sequential)
// 3. Provide a clear description
// 4. Write the SQL in the 'up' field
//
// Example:
//   {
//       version:     2,
//       description: "Add genre column to albums",
//       up: `ALTER TABLE albums ADD COLUMN genre VARCHAR(50);`,
//   },
//
// Important rules:
// - NEVER modify existing migrations (breaks version tracking)
// - ALWAYS increment version sequentially (1, 2, 3, not 1, 3, 5)
// - Each migration runs in a transaction (auto-rollback on error)
// - Applied migrations are tracked in 'schema_migrations' table
var migrations = []migration{
	{
		version:     1,
		description: "Create albums table",
		up: `
			CREATE TABLE IF NOT EXISTS albums (
				id SERIAL PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				artist VARCHAR(255) NOT NULL,
				price DECIMAL(10, 2) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP
			);

			CREATE INDEX IF NOT EXISTS idx_albums_deleted_at ON albums(deleted_at);
		`,
	},
	// Add future migrations here:
	// {
	//     version:     2,
	//     description: "Add genre column to albums",
	//     up: `ALTER TABLE albums ADD COLUMN genre VARCHAR(50);`,
	// },
}

// Migrate runs database migrations
func (d *Database) Migrate(ctx context.Context) error {
	// Create migrations tracking table if it doesn't exist
	createTrackingTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := d.Pool.Exec(ctx, createTrackingTable); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get current migration version
	var currentVersion int
	err := d.Pool.QueryRow(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	log.Printf("Current database version: %d", currentVersion)

	// Run pending migrations
	migrationsRan := 0
	for _, m := range migrations {
		if m.version <= currentVersion {
			continue // Skip already applied migrations
		}

		log.Printf("Running migration %d: %s", m.version, m.description)

		// Start transaction for this migration
		tx, err := d.Pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.version, err)
		}

		// Run the migration
		if _, err := tx.Exec(ctx, m.up); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("migration %d failed: %w", m.version, err)
		}

		// Record that this migration was applied
		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, description) VALUES ($1, $2)", m.version, m.description); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration %d: %w", m.version, err)
		}

		// Commit transaction
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.version, err)
		}

		log.Printf("✓ Migration %d completed: %s", m.version, m.description)
		migrationsRan++
	}

	if migrationsRan == 0 {
		log.Println("Database schema is up to date")
	} else {
		log.Printf("Successfully ran %d migration(s)", migrationsRan)
	}

	return nil
}

// Close closes the database connection pool
func (d *Database) Close() {
	if d.Pool != nil {
		log.Println("Closing database connection...")
		d.Pool.Close()
	}
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
