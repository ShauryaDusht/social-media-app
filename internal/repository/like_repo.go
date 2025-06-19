package repository

import (
	"social-media-app/internal/models"

	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

// Delete removes a like based on who liked it (liked_by) and post_id
func (r *likeRepository) Delete(likedBy, postID uint) error {
	return r.db.Where("liked_by = ? AND post_id = ?", likedBy, postID).Delete(&models.Like{}).Error
}

// Exists checks if a user has already liked a post
func (r *likeRepository) Exists(likedBy, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("liked_by = ? AND post_id = ?", likedBy, postID).Count(&count).Error
	return count > 0, err
}

func (r *likeRepository) GetByPostID(postID uint) ([]models.Like, error) {
	var likes []models.Like
	err := r.db.Preload("User").Where("post_id = ?", postID).Find(&likes).Error
	return likes, err
}

// GetLikeCount returns the total number of likes for a post
func (r *likeRepository) GetLikeCount(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// GetLikedUserIDs returns all user IDs who liked a specific post
func (r *likeRepository) GetLikedUserIDs(postID uint) ([]uint, error) {
	var userIDs []uint
	err := r.db.Model(&models.Like{}).Where("post_id = ?", postID).Pluck("liked_by", &userIDs).Error
	return userIDs, err
}
