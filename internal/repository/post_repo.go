package repository

import (
	"social-media-app/internal/models"

	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetByUserID(userID uint, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("User").Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) GetAll(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) Update(post *models.Post) error {
	// Use Updates instead of Save to only update specific fields
	// This prevents updating timestamps automatically
	return r.db.Model(post).Updates(map[string]interface{}{
		"content":   post.Content,
		"image_url": post.ImageURL,
	}).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) GetTimeline(userID uint, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("User").
		Where("user_id IN (SELECT following_id FROM follows WHERE follower_id = ?) OR user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error
	return posts, err
}
