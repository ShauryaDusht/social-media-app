package models

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	PostID    uint           `json:"post_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
	Post Post `json:"post" gorm:"foreignKey:PostID"`
}

// Ensure unique constraint on user_id and post_id combination
func (Like) TableName() string {
	return "likes"
}

// LikeRequest represents a like/unlike request
type LikeRequest struct {
	PostID uint `json:"post_id" binding:"required"`
}
