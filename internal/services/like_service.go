package services

import (
	"errors"
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
)

type LikeService struct {
	likeRepo repository.LikeRepository
}

func NewLikeService(likeRepo repository.LikeRepository) *LikeService {
	return &LikeService{
		likeRepo: likeRepo,
	}
}

func (s *LikeService) LikePost(userID, postID uint) error {
	exists, err := s.likeRepo.Exists(userID, postID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("post already liked by this user")
	}

	like := &models.Like{
		UserID: userID,
		PostID: postID,
	}

	return s.likeRepo.Create(like)
}

func (s *LikeService) UnlikePost(userID, postID uint) error {
	_, err := s.likeRepo.Exists(userID, postID)
	if err != nil {
		return err
	}

	return s.likeRepo.Delete(userID, postID)
}
