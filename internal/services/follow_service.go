package services

import (
	"errors"
	"strings"

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

	// Check if target user exists
	_, err := s.userRepo.GetByID(followingID)
	if err != nil {
		return errors.New("user not found")
	}

	// Try to create the follow relationship
	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	err = s.followRepo.Create(follow)
	if err != nil {
		// Handle unique constraint violations gracefully
		errorStr := strings.ToLower(err.Error())
		if strings.Contains(errorStr, "unique") ||
			strings.Contains(errorStr, "duplicate") ||
			strings.Contains(errorStr, "unique_follower_following") ||
			strings.Contains(errorStr, "constraint") {
			// If it's a unique constraint violation, just return success
			// This makes the operation idempotent
			return nil
		}
		return err
	}

	// Clear cache after successful follow
	s.cacheRepo.DeleteTimeline(followerID)
	return nil
}

func (s *FollowService) UnfollowUser(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("cannot unfollow yourself")
	}

	// Just attempt to delete - if it doesn't exist, that's fine
	err := s.followRepo.Delete(followerID, followingID)
	if err != nil {
		return err
	}

	// Clear cache after operation
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
