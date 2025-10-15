package controllers

import (
	"errors"
	"example/web-service-gin/models"
	"example/web-service-gin/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AlbumController handles album HTTP requests
type AlbumController struct {
	repo repository.AlbumRepository
}

// NewAlbumController creates a new album controller
func NewAlbumController(repo repository.AlbumRepository) *AlbumController {
	return &AlbumController{repo: repo}
}

// GetAlbums retrieves all albums
func (ctrl *AlbumController) GetAlbums(c *gin.Context) {
	albums, err := ctrl.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve albums"})
		return
	}

	c.JSON(http.StatusOK, albums)
}

// GetAlbum retrieves a single album by ID
func (ctrl *AlbumController) GetAlbum(c *gin.Context) {
	// Validate and parse ID
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	album, err := ctrl.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// CreateAlbum creates a new album
func (ctrl *AlbumController) CreateAlbum(c *gin.Context) {
	var album models.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.repo.Create(c.Request.Context(), &album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}

	c.JSON(http.StatusCreated, album)
}

// UpdateAlbum updates an existing album
func (ctrl *AlbumController) UpdateAlbum(c *gin.Context) {
	// Validate and parse ID
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Check if album exists
	album, err := ctrl.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album"})
		return
	}

	// Bind the updated data
	var updatedAlbum models.Album
	if err := c.ShouldBindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the album fields
	album.Title = updatedAlbum.Title
	album.Artist = updatedAlbum.Artist
	album.Price = updatedAlbum.Price

	if err := ctrl.repo.Update(c.Request.Context(), album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// DeleteAlbum deletes an album by ID
func (ctrl *AlbumController) DeleteAlbum(c *gin.Context) {
	// Validate and parse ID
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Delete the album
	if err := ctrl.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}
