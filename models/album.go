package models

import "gorm.io/gorm"

// Album represents an album record in the database
type Album struct {
	gorm.Model
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}
