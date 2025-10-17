package album

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles album HTTP requests
type Handler struct {
	repo Repository
}

// NewHandler creates a new album handler
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// RegisterRoutes registers album routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetAlbums)
	router.GET("/:id", h.GetAlbum)
	router.POST("", h.CreateAlbum)
	router.PUT("/:id", h.UpdateAlbum)
	router.DELETE("/:id", h.DeleteAlbum)
}

// GetAlbums retrieves all albums
func (h *Handler) GetAlbums(c *gin.Context) {
	albums, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve albums"})
		return
	}

	c.JSON(http.StatusOK, albums)
}

// GetAlbum retrieves a single album by ID
func (h *Handler) GetAlbum(c *gin.Context) {
	// Validate and parse ID
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	album, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// CreateAlbum creates a new album
func (h *Handler) CreateAlbum(c *gin.Context) {
	var album Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), &album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}

	c.JSON(http.StatusCreated, album)
}

// UpdateAlbum updates an existing album
func (h *Handler) UpdateAlbum(c *gin.Context) {
	// Validate and parse ID
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Check if album exists
	album, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album"})
		return
	}

	// Bind the updated data
	var updatedAlbum Album
	if err := c.ShouldBindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the album fields
	album.Title = updatedAlbum.Title
	album.Artist = updatedAlbum.Artist
	album.Price = updatedAlbum.Price

	if err := h.repo.Update(c.Request.Context(), album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// DeleteAlbum deletes an album by ID
func (h *Handler) DeleteAlbum(c *gin.Context) {
	// Validate and parse ID
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Delete the album
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}
