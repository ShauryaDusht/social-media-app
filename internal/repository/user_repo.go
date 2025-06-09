package repository

import (
	"social-media-app/internal/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID gets user by ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername gets user by username
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates user
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes user
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// EmailExists checks if email already exists
func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// UsernameExists checks if username already exists
func (r *userRepository) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// SearchUsers searches for users by name or username
func (r *userRepository) SearchUsers(query string, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("username LIKE ? OR first_name LIKE ? OR last_name LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// GetFollowers gets all users who follow the specified user
func (r *userRepository) GetFollowers(userID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.Table("users").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Find(&users).Error
	return users, err
}
