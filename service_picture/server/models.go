package server

import (
	"time"
)

// PictureResponse represents the response for a picture
type PictureResponse struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	ViewCount int       `json:"view_count,omitempty"`
}
