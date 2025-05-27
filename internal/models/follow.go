package models

import (
	"time"

	"gorm.io/gorm"
)

type Follow struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FollowerID  uint           `json:"follower_id" gorm:"not null;index"`  // User who follows
	FollowingID uint           `json:"following_id" gorm:"not null;index"` // User being followed
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Follower  User `json:"follower" gorm:"foreignKey:FollowerID"`
	Following User `json:"following" gorm:"foreignKey:FollowingID"`
}

// Ensure unique constraint on follower_id and following_id combination
func (Follow) TableName() string {
	return "follows"
}

// FollowRequest represents a follow/unfollow request
type FollowRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// FollowResponse represents follow relationship info
type FollowResponse struct {
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}
