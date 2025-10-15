package repository

import (
	"context"
	"errors"
	"example/web-service-gin/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("record not found")
)

// AlbumRepository handles album data access
type AlbumRepository interface {
	FindAll(ctx context.Context) ([]models.Album, error)
	FindByID(ctx context.Context, id int) (*models.Album, error)
	Create(ctx context.Context, album *models.Album) error
	Update(ctx context.Context, album *models.Album) error
	Delete(ctx context.Context, id int) error
}

// albumRepository implements AlbumRepository
type albumRepository struct {
	pool *pgxpool.Pool
}

// NewAlbumRepository creates a new album repository
func NewAlbumRepository(pool *pgxpool.Pool) AlbumRepository {
	return &albumRepository{pool: pool}
}

// FindAll retrieves all albums (excluding soft-deleted)
func (r *albumRepository) FindAll(ctx context.Context) ([]models.Album, error) {
	query := `
		SELECT id, title, artist, price, created_at, updated_at, deleted_at
		FROM albums
		WHERE deleted_at IS NULL
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var album models.Album
		err := rows.Scan(
			&album.ID,
			&album.Title,
			&album.Artist,
			&album.Price,
			&album.CreatedAt,
			&album.UpdatedAt,
			&album.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

// FindByID retrieves a single album by ID
func (r *albumRepository) FindByID(ctx context.Context, id int) (*models.Album, error) {
	query := `
		SELECT id, title, artist, price, created_at, updated_at, deleted_at
		FROM albums
		WHERE id = $1 AND deleted_at IS NULL
	`

	var album models.Album
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&album.ID,
		&album.Title,
		&album.Artist,
		&album.Price,
		&album.CreatedAt,
		&album.UpdatedAt,
		&album.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &album, nil
}

// Create creates a new album
func (r *albumRepository) Create(ctx context.Context, album *models.Album) error {
	query := `
		INSERT INTO albums (title, artist, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.pool.QueryRow(
		ctx,
		query,
		album.Title,
		album.Artist,
		album.Price,
		now,
		now,
	).Scan(&album.ID, &album.CreatedAt, &album.UpdatedAt)

	return err
}

// Update updates an existing album
func (r *albumRepository) Update(ctx context.Context, album *models.Album) error {
	query := `
		UPDATE albums
		SET title = $1, artist = $2, price = $3, updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING updated_at
	`

	now := time.Now()
	err := r.pool.QueryRow(
		ctx,
		query,
		album.Title,
		album.Artist,
		album.Price,
		now,
		album.ID,
	).Scan(&album.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// Delete deletes an album by ID (soft delete)
func (r *albumRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE albums
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
