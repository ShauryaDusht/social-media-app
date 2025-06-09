package repository

import "social-media-app/internal/models"

// UserRepository defines user database operations
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
	SearchUsers(query string, limit, offset int) ([]models.User, error)
	GetFollowers(userID uint) ([]models.User, error)
}

// PostRepository defines post database operations
type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	GetByUserID(userID uint, limit, offset int) ([]models.Post, error)
	GetAll(limit, offset int) ([]models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
	GetTimeline(userID uint, limit, offset int) ([]models.Post, error)
}

// LikeRepository defines like database operations
type LikeRepository interface {
	Create(like *models.Like) error
	Delete(userID, postID uint) error
	Exists(userID, postID uint) (bool, error)
	GetByPostID(postID uint) ([]models.Like, error)
}

// FollowRepository defines follow database operations
type FollowRepository interface {
	Create(follow *models.Follow) error
	Delete(followerID, followingID uint) error
	Exists(followerID, followingID uint) (bool, error)
	GetFollowers(userID uint) ([]models.Follow, error)
	GetFollowing(userID uint) ([]models.Follow, error)
}
