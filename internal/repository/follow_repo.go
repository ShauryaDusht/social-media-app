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
	// Use a transaction to ensure atomicity
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(follow).Error
	})
}

func (r *followRepository) Delete(followerID, followingID uint) error {
	// Use transaction for consistency
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Use Unscoped to permanently delete the record instead of soft delete
		result := tx.Unscoped().Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{})
		if result.Error != nil {
			return result.Error
		}
		// Even if no rows were affected, don't return an error (idempotent operation)
		return nil
	})
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
