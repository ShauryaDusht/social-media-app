package models

import (
	"time"
)

// UpdateProfileRequest represents the request to update a user's profile
type UpdateProfileRequest struct {
	FirstName string    `json:"first_name" binding:"omitempty,min=1,max=50"`
	LastName  string    `json:"last_name" binding:"omitempty,min=1,max=50"`
	Bio       string    `json:"bio" binding:"omitempty,min=1,max=500"`
	Avatar    string    `json:"avatar" binding:"omitempty,url,max=255"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime"`
	Password  string    `json:"password" binding:"omitempty,min=8"` // optional
}
