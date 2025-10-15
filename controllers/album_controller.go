package controllers

import (
	"example/web-service-gin/database"
	"example/web-service-gin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAlbums retrieves all albums
func GetAlbums(c *gin.Context) {
	var albums []models.Album

	result := database.DB.Find(&albums)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, albums)
}

// GetAlbum retrieves a single album by ID
func GetAlbum(c *gin.Context) {
	id := c.Param("id")
	var album models.Album

	result := database.DB.First(&album, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// CreateAlbum creates a new album
func CreateAlbum(c *gin.Context) {
	var album models.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&album)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, album)
}

// UpdateAlbum updates an existing album
func UpdateAlbum(c *gin.Context) {
	id := c.Param("id")
	var album models.Album

	// Check if album exists
	if result := database.DB.First(&album, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// Bind the updated data
	var updatedAlbum models.Album
	if err := c.ShouldBindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the album
	album.Title = updatedAlbum.Title
	album.Artist = updatedAlbum.Artist
	album.Price = updatedAlbum.Price

	if result := database.DB.Save(&album); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

// DeleteAlbum deletes an album by ID
func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	var album models.Album

	// Check if album exists
	if result := database.DB.First(&album, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// Delete the album
	if result := database.DB.Delete(&album); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}
