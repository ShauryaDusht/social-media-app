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

func (s *LikeService) LikePost(userID, postID uint) error {
	exists, err := s.likeRepo.Exists(userID, postID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("post already liked")
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	like := &models.Like{
		UserID: userID,
		PostID: postID,
	}

	if err := s.likeRepo.Create(like); err != nil {
		return err
	}

	post.LikeCount++
	return s.postRepo.Update(post)
}

func (s *LikeService) UnlikePost(userID, postID uint) error {
	exists, err := s.likeRepo.Exists(userID, postID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("post not liked")
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	if err := s.likeRepo.Delete(userID, postID); err != nil {
		return err
	}

	if post.LikeCount > 0 {
		post.LikeCount--
	}
	return s.postRepo.Update(post)
}
