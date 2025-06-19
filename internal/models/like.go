package models

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	PostID    uint           `json:"post_id" gorm:"not null;index"`
	LikedBy   uint           `json:"liked_by" gorm:"not null;index"` // New field to track who liked
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
	Post Post `json:"post" gorm:"foreignKey:PostID"`
}

// Ensure unique constraint on liked_by and post_id combination
func (Like) TableName() string {
	return "likes"
}

// LikeRequest represents a like/unlike request
type LikeRequest struct {
	PostID uint `json:"post_id" binding:"required"`
}

// GetLikedUserIDs returns a slice of user IDs who liked a specific post
func GetLikedUserIDs(db *gorm.DB, postID uint) ([]uint, error) {
	var userIDs []uint
	err := db.Model(&Like{}).Where("post_id = ?", postID).Pluck("liked_by", &userIDs).Error
	return userIDs, err
}

// IsPostLikedByUser checks if a specific user has liked a specific post
func IsPostLikedByUser(db *gorm.DB, postID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&Like{}).Where("post_id = ? AND liked_by = ?", postID, userID).Count(&count).Error
	return count > 0, err
}
