package services

import (
	"errors"
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
)

type LikeService struct {
	likeRepo repository.LikeRepository
	postRepo repository.PostRepository
}

func NewLikeService(likeRepo repository.LikeRepository, postRepo repository.PostRepository) *LikeService {
	return &LikeService{
		likeRepo: likeRepo,
		postRepo: postRepo,
	}
}

func (s *LikeService) LikePost(likedBy, postID uint) error {
	// Check if post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	// Check if user has already liked this post
	exists, err := s.likeRepo.Exists(likedBy, postID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("post already liked by this user")
	}

	// Create new like record
	like := &models.Like{
		UserID:  likedBy, // Keep for backward compatibility
		PostID:  postID,
		LikedBy: likedBy, // New field to track who liked
	}

	return s.likeRepo.Create(like)
}

func (s *LikeService) UnlikePost(likedBy, postID uint) error {
	// Check if post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	// Check if user has liked this post
	exists, err := s.likeRepo.Exists(likedBy, postID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("post not liked by this user")
	}

	// Remove the like
	return s.likeRepo.Delete(likedBy, postID)
}

func (s *LikeService) GetPostLikes(postID uint) ([]models.Like, error) {
	return s.likeRepo.GetByPostID(postID)
}

func (s *LikeService) GetLikeCount(postID uint) (int64, error) {
	return s.likeRepo.GetLikeCount(postID)
}

func (s *LikeService) IsPostLikedByUser(postID, userID uint) (bool, error) {
	return s.likeRepo.Exists(userID, postID)
}

func (s *LikeService) GetLikedUserIDs(postID uint) ([]uint, error) {
	return s.likeRepo.GetLikedUserIDs(postID)
}
