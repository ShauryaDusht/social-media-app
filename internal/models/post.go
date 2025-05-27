package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Content   string         `json:"content" gorm:"not null;size:1000"`
	ImageURL  string         `json:"image_url" gorm:"size:255"`
	LikeCount int            `json:"like_count" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User  User   `json:"user" gorm:"foreignKey:UserID"`
	Likes []Like `json:"likes,omitempty" gorm:"foreignKey:PostID"`
}

// PostResponse represents a post with user info for API responses
type PostResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	ImageURL  string       `json:"image_url"`
	LikeCount int          `json:"like_count"`
	CreatedAt time.Time    `json:"created_at"`
	User      UserResponse `json:"user"`
	IsLiked   bool         `json:"is_liked"` // Whether current user liked this post
}

// CreatePostRequest represents the request to create a new post
type CreatePostRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ImageURL string `json:"image_url"`
}

// UpdatePostRequest represents the request to update a post
type UpdatePostRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ImageURL string `json:"image_url"`
}

// ToResponse converts Post to PostResponse
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:        p.ID,
		Content:   p.Content,
		ImageURL:  p.ImageURL,
		LikeCount: p.LikeCount,
		CreatedAt: p.CreatedAt,
		User:      p.User.ToResponse(),
		IsLiked:   false, // This will be set by the service layer
	}
}
