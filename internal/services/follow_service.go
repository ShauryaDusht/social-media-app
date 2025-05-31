package services

import (
	"errors"
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
)

type FollowService struct {
	followRepo repository.FollowRepository
}

func NewFollowService(followRepo repository.FollowRepository) *FollowService {
	return &FollowService{
		followRepo: followRepo,
	}
}

func (s *FollowService) FollowUser(followerID, followingID uint) error {
	exists, err := s.followRepo.Exists(followerID, followingID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already following this user")
	}

	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return s.followRepo.Create(follow)
}

func (s *FollowService) UnfollowUser(followerID, followingID uint) error {
	return s.followRepo.Delete(followerID, followingID)
}

func (s *FollowService) GetFollowers(userID uint) ([]models.Follow, error) {
	return s.followRepo.GetFollowers(userID)
}

func (s *FollowService) GetFollowing(userID uint) ([]models.Follow, error) {
	return s.followRepo.GetFollowing(userID)
}
