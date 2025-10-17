package album

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test database connection
func setupTestDB(t *testing.T) *pgxpool.Pool {
	// Use a test database URL
	// In real scenarios, you'd want to use a separate test database
	dbURL := "postgres://postgres:postgres@localhost:5432/albums_test?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skip("Skipping test: PostgreSQL not available")
		return nil
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		t.Skip("Skipping test: Cannot connect to PostgreSQL")
		return nil
	}

	return pool
}

// Helper function to clean up test database
func cleanupTestDB(t *testing.T, pool *pgxpool.Pool) {
	if pool != nil {
		_, err := pool.Exec(context.Background(), "DELETE FROM albums")
		if err != nil {
			t.Logf("Warning: Failed to cleanup test data: %v", err)
		}
		pool.Close()
	}
}

// Helper function to create test albums table
func createTestTable(t *testing.T, pool *pgxpool.Pool) {
	query := `
	CREATE TABLE IF NOT EXISTS albums (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		artist VARCHAR(255) NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP
	);
	`
	_, err := pool.Exec(context.Background(), query)
	require.NoError(t, err)
}

func TestAlbumRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	album := &Album{
		Title:  "The Wall",
		Artist: "Pink Floyd",
		Price:  24.99,
	}

	err := repo.Create(context.Background(), album)
	require.NoError(t, err)

	// Verify album was created
	assert.NotZero(t, album.ID, "Album ID should be set after creation")
	assert.NotZero(t, album.CreatedAt, "CreatedAt should be set")
	assert.NotZero(t, album.UpdatedAt, "UpdatedAt should be set")
	assert.Equal(t, "The Wall", album.Title)
	assert.Equal(t, "Pink Floyd", album.Artist)
	assert.Equal(t, 24.99, album.Price)
}

func TestAlbumRepository_FindAll(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Create test albums
	album1 := &Album{Title: "The Wall", Artist: "Pink Floyd", Price: 24.99}
	album2 := &Album{Title: "Dark Side of the Moon", Artist: "Pink Floyd", Price: 22.99}

	require.NoError(t, repo.Create(context.Background(), album1))
	require.NoError(t, repo.Create(context.Background(), album2))

	// Find all albums
	albums, err := repo.FindAll(context.Background())
	require.NoError(t, err)

	assert.Len(t, albums, 2)
	assert.Equal(t, "The Wall", albums[0].Title)
	assert.Equal(t, "Dark Side of the Moon", albums[1].Title)
}

func TestAlbumRepository_FindByID(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Create a test album
	album := &Album{Title: "The Wall", Artist: "Pink Floyd", Price: 24.99}
	require.NoError(t, repo.Create(context.Background(), album))

	// Find by ID
	found, err := repo.FindByID(context.Background(), album.ID)
	require.NoError(t, err)

	assert.Equal(t, album.ID, found.ID)
	assert.Equal(t, "The Wall", found.Title)
	assert.Equal(t, "Pink Floyd", found.Artist)
	assert.Equal(t, 24.99, found.Price)
}

func TestAlbumRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Try to find non-existent album
	_, err := repo.FindByID(context.Background(), 99999)
	assert.Equal(t, ErrNotFound, err)
}

func TestAlbumRepository_Update(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Create a test album
	album := &Album{Title: "The Wall", Artist: "Pink Floyd", Price: 24.99}
	require.NoError(t, repo.Create(context.Background(), album))

	// Update the album
	album.Price = 29.99
	err := repo.Update(context.Background(), album)
	require.NoError(t, err)

	// Verify update
	updated, err := repo.FindByID(context.Background(), album.ID)
	require.NoError(t, err)
	assert.Equal(t, 29.99, updated.Price)
	assert.True(t, updated.UpdatedAt.After(updated.CreatedAt), "UpdatedAt should be after CreatedAt")
}

func TestAlbumRepository_Update_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Try to update non-existent album
	album := &Album{ID: 99999, Title: "Test", Artist: "Test", Price: 9.99}
	err := repo.Update(context.Background(), album)
	assert.Equal(t, ErrNotFound, err)
}

func TestAlbumRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Create a test album
	album := &Album{Title: "The Wall", Artist: "Pink Floyd", Price: 24.99}
	require.NoError(t, repo.Create(context.Background(), album))

	// Delete the album (soft delete)
	err := repo.Delete(context.Background(), album.ID)
	require.NoError(t, err)

	// Verify album is not found (soft deleted)
	_, err = repo.FindByID(context.Background(), album.ID)
	assert.Equal(t, ErrNotFound, err)

	// Verify FindAll doesn't return deleted albums
	albums, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, albums, 0)
}

func TestAlbumRepository_Delete_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Try to delete non-existent album
	err := repo.Delete(context.Background(), 99999)
	assert.Equal(t, ErrNotFound, err)
}

func TestAlbumRepository_ContextCancellation(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Create a context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Try to create album with cancelled context
	album := &Album{Title: "Test", Artist: "Test", Price: 9.99}
	err := repo.Create(ctx, album)
	assert.Error(t, err, "Should return error for cancelled context")
}

func TestAlbumRepository_ConcurrentCreates(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		return
	}
	defer cleanupTestDB(t, pool)

	createTestTable(t, pool)
	repo := NewRepository(pool)

	// Create multiple albums concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(num int) {
			album := &Album{
				Title:  "Test Album",
				Artist: "Test Artist",
				Price:  9.99,
			}
			err := repo.Create(context.Background(), album)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all albums were created
	albums, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, albums, 10)
}
