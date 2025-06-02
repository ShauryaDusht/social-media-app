package services

import (
	"errors"

	"social-media-app/internal/models"
	"social-media-app/internal/repository"
)

type FollowService struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
	cacheRepo  repository.CacheRepository
}

func NewFollowService(followRepo repository.FollowRepository, userRepo repository.UserRepository, cacheRepo repository.CacheRepository) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
	}
}

func (s *FollowService) FollowUser(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	_, err := s.userRepo.GetByID(followingID)
	if err != nil {
		return errors.New("user not found")
	}

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

	err = s.followRepo.Create(follow)
	if err != nil {
		return err
	}

	s.cacheRepo.DeleteTimeline(followerID)
	return nil
}

func (s *FollowService) UnfollowUser(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("cannot unfollow yourself")
	}

	exists, err := s.followRepo.Exists(followerID, followingID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("not following this user")
	}

	err = s.followRepo.Delete(followerID, followingID)
	if err != nil {
		return err
	}

	s.cacheRepo.DeleteTimeline(followerID)
	return nil
}

func (s *FollowService) GetFollowers(userID uint) ([]models.FollowResponse, error) {
	follows, err := s.followRepo.GetFollowers(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.FollowResponse
	for _, follow := range follows {
		responses = append(responses, models.FollowResponse{
			User:      follow.Follower.ToResponse(),
			CreatedAt: follow.CreatedAt,
		})
	}

	return responses, nil
}

func (s *FollowService) GetFollowing(userID uint) ([]models.FollowResponse, error) {
	follows, err := s.followRepo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.FollowResponse
	for _, follow := range follows {
		responses = append(responses, models.FollowResponse{
			User:      follow.Following.ToResponse(),
			CreatedAt: follow.CreatedAt,
		})
	}

	return responses, nil
}
