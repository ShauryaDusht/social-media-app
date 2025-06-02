package repository

import (
	"social-media-app/internal/models"

	"gorm.io/gorm"
)

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(follow *models.Follow) error {
	return r.db.Create(follow).Error
}

func (r *followRepository) Delete(followerID, followingID uint) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{}).Error
}

func (r *followRepository) Exists(followerID, followingID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}

func (r *followRepository) GetFollowers(userID uint) ([]models.Follow, error) {
	var follows []models.Follow
	err := r.db.Preload("Follower").Where("following_id = ?", userID).Find(&follows).Error
	return follows, err
}

func (r *followRepository) GetFollowing(userID uint) ([]models.Follow, error) {
	var follows []models.Follow
	err := r.db.Preload("Following").Where("follower_id = ?", userID).Find(&follows).Error
	return follows, err
}
