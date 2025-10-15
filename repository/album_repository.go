package repository

import (
	"errors"
	"example/web-service-gin/models"

	"gorm.io/gorm"
)

// AlbumRepository handles album data access
type AlbumRepository interface {
	FindAll() ([]models.Album, error)
	FindByID(id uint) (*models.Album, error)
	Create(album *models.Album) error
	Update(album *models.Album) error
	Delete(id uint) error
}

// albumRepository implements AlbumRepository
type albumRepository struct {
	db *gorm.DB
}

// NewAlbumRepository creates a new album repository
func NewAlbumRepository(db *gorm.DB) AlbumRepository {
	return &albumRepository{db: db}
}

// FindAll retrieves all albums
func (r *albumRepository) FindAll() ([]models.Album, error) {
	var albums []models.Album
	result := r.db.Find(&albums)
	return albums, result.Error
}

// FindByID retrieves a single album by ID
func (r *albumRepository) FindByID(id uint) (*models.Album, error) {
	var album models.Album
	result := r.db.First(&album, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &album, nil
}

// Create creates a new album
func (r *albumRepository) Create(album *models.Album) error {
	return r.db.Create(album).Error
}

// Update updates an existing album
func (r *albumRepository) Update(album *models.Album) error {
	return r.db.Save(album).Error
}

// Delete deletes an album by ID (soft delete)
func (r *albumRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Album{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
