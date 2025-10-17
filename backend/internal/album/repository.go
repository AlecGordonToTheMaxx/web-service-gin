package album

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrNotFound is returned when a record is not found
	ErrNotFound = errors.New("record not found")
)

// Repository handles album data access
type Repository interface {
	FindAll(ctx context.Context) ([]Album, error)
	FindByID(ctx context.Context, id int) (*Album, error)
	Create(ctx context.Context, album *Album) error
	Update(ctx context.Context, album *Album) error
	Delete(ctx context.Context, id int) error
}

// repository implements Repository
type repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new album repository
func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

// FindAll retrieves all albums (excluding soft-deleted)
func (r *repository) FindAll(ctx context.Context) ([]Album, error) {
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

	var albums []Album
	for rows.Next() {
		var album Album
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
func (r *repository) FindByID(ctx context.Context, id int) (*Album, error) {
	query := `
		SELECT id, title, artist, price, created_at, updated_at, deleted_at
		FROM albums
		WHERE id = $1 AND deleted_at IS NULL
	`

	var album Album
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
func (r *repository) Create(ctx context.Context, album *Album) error {
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
func (r *repository) Update(ctx context.Context, album *Album) error {
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
func (r *repository) Delete(ctx context.Context, id int) error {
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
