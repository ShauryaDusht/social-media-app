package services

import (
	"errors"
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
	"time"
)

type PostService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) CreatePost(userID uint, content string, imageURL string) error {
	post := &models.Post{
		Content:   content,
		ImageURL:  imageURL,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return s.postRepo.Create(post)
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetByUserID(userID uint, limit, offset int) ([]models.Post, error) {
	posts, err := s.postRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetAll(limit, offset int) ([]models.Post, error) {
	posts, err := s.postRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) Update(post *models.Post) error {
	exists, err := s.Exists(post.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("post not found")
	}
	return s.postRepo.Update(post)
}

func (s *PostService) Delete(id uint) error {
	exists, err := s.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("post not found")
	}
	return s.postRepo.Delete(id)
}

func (s *PostService) Exists(id uint) (bool, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return false, err
	}
	return post != nil, nil
}

func (s *PostService) GetTimeline(userID uint, limit, offset int) ([]models.Post, error) {
	posts, err := s.postRepo.GetTimeline(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
