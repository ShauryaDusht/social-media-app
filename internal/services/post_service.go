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
	userRepo  repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository, likeRepo repository.LikeRepository, cacheRepo repository.CacheRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo:  postRepo,
		likeRepo:  likeRepo,
		cacheRepo: cacheRepo,
		userRepo:  userRepo,
	}
}

func (s *PostService) CreatePost(userID uint, content, imageURL string) error {
	post := &models.Post{
		UserID:   userID,
		Content:  content,
		ImageURL: imageURL,
	}

	err := s.postRepo.Create(post)
	if err != nil {
		return err
	}

	// Invalidate timeline cache for all followers
	s.invalidateFollowersTimeline(userID)

	return nil
}

func (s *PostService) GetByID(postID uint) (*models.PostResponse, error) {
	// Try to get from cache first
	if cached, err := s.cacheRepo.GetPostCache(postID); err == nil {
		return cached, nil
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	response := post.ToResponse()
	// Cache the post for 10 minutes
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
	// Invalidate post cache
	s.cacheRepo.DeletePostCache(post.ID)

	// Invalidate timeline cache for all followers
	s.invalidateFollowersTimeline(post.UserID)

	return s.postRepo.Update(post)
}

func (s *PostService) Delete(postID uint) error {
	// Get post to find user ID before deletion
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	// Invalidate post cache
	s.cacheRepo.DeletePostCache(postID)

	// Invalidate timeline cache for all followers
	s.invalidateFollowersTimeline(post.UserID)

	return s.postRepo.Delete(postID)
}

func (s *PostService) GetTimeline(userID uint, limit, offset int) ([]models.PostResponse, error) {
	// Try to get from cache first
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

	// Cache the timeline for 5 minutes
	if len(responses) > 0 {
		s.cacheRepo.SetTimeline(userID, responses, 5*time.Minute)
	}

	return responses, nil
}

// Helper function to invalidate timeline cache for all followers of a user
func (s *PostService) invalidateFollowersTimeline(userID uint) {
	// Invalidate the user's own timeline
	s.cacheRepo.DeleteTimeline(userID)

	// Get all followers of the user
	followers, err := s.userRepo.GetFollowers(userID)
	if err != nil {
		return
	}

	// Invalidate timeline cache for each follower
	for _, follower := range followers {
		s.cacheRepo.DeleteTimeline(follower.ID)
	}
}
