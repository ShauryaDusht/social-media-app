package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"size:50"`
	LastName  string         `json:"last_name" gorm:"size:50"`
	Bio       string         `json:"bio" gorm:"size:500"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Posts     []Post   `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	Likes     []Like   `json:"likes,omitempty" gorm:"foreignKey:UserID"`
	Followers []Follow `json:"followers,omitempty" gorm:"foreignKey:FollowingID"`
	Following []Follow `json:"following,omitempty" gorm:"foreignKey:FollowerID"`
}

// UserResponse is used for API responses (excludes sensitive data)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Bio:       u.Bio,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
	}
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
}

// LoginRequest represents user login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
