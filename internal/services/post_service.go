package services

import (
	"time"

	"social-media-app/internal/models"
	"social-media-app/internal/repository"
)

type PostService struct {
	postRepo  repository.PostRepository
	likeRepo  repository.LikeRepository
	cacheRepo repository.CacheRepository
}

func NewPostService(postRepo repository.PostRepository, likeRepo repository.LikeRepository, cacheRepo repository.CacheRepository) *PostService {
	return &PostService{
		postRepo:  postRepo,
		likeRepo:  likeRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *PostService) CreatePost(userID uint, content, imageURL string) error {
	post := &models.Post{
		UserID:   userID,
		Content:  content,
		ImageURL: imageURL,
	}

	return s.postRepo.Create(post)
}

func (s *PostService) GetByID(postID uint) (*models.PostResponse, error) {
	if cached, err := s.cacheRepo.GetPostCache(postID); err == nil {
		return cached, nil
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	response := post.ToResponse()
	s.cacheRepo.SetPostCache(postID, response, 10*time.Minute)

	return &response, nil
}

func (s *PostService) GetAll(limit, offset int) ([]models.PostResponse, error) {
	posts, err := s.postRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}

	return responses, nil
}

func (s *PostService) GetByUserID(userID uint, limit, offset int) ([]models.PostResponse, error) {
	posts, err := s.postRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}

	return responses, nil
}

func (s *PostService) Update(post *models.Post) error {
	s.cacheRepo.DeletePostCache(post.ID)
	return s.postRepo.Update(post)
}

func (s *PostService) Delete(postID uint) error {
	s.cacheRepo.DeletePostCache(postID)
	return s.postRepo.Delete(postID)
}

func (s *PostService) GetTimeline(userID uint, limit, offset int) ([]models.PostResponse, error) {
	if cached, err := s.cacheRepo.GetTimeline(userID); err == nil && len(cached) > 0 {
		start := offset
		end := offset + limit
		if start >= len(cached) {
			return []models.PostResponse{}, nil
		}
		if end > len(cached) {
			end = len(cached)
		}
		return cached[start:end], nil
	}

	posts, err := s.postRepo.GetTimeline(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		response := post.ToResponse()

		isLiked, _ := s.likeRepo.Exists(userID, post.ID)
		response.IsLiked = isLiked

		responses = append(responses, response)
	}

	if len(responses) > 0 {
		s.cacheRepo.SetTimeline(userID, responses, 5*time.Minute)
	}

	return responses, nil
}
