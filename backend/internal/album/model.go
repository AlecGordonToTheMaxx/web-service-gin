package album

import (
	"time"
)

// Album represents an album record in the database
type Album struct {
	ID        int        `json:"id"`
	Title     string     `json:"title" binding:"required"`
	Artist    string     `json:"artist" binding:"required"`
	Price     float64    `json:"price" binding:"required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
