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

func (r *likeRepository) Delete(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Like{}).Error
}

func (r *likeRepository) Exists(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

func (r *likeRepository) GetByPostID(postID uint) ([]models.Like, error) {
	var likes []models.Like
	err := r.db.Preload("User").Where("post_id = ?", postID).Find(&likes).Error
	return likes, err
}
