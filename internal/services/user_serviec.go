package services

import (
	"errors"
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
	"social-media-app/internal/utils"
)

type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile gets user profile by ID
func (s *UserService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uint, req models.UpdateProfileRequest) (*models.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// Update password if provided
	if req.Password != "" {
		if len(req.Password) < 6 {
			return nil, errors.New("password must be at least 6 characters")
		}
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = hashedPassword
	}

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByID gets user by ID (public profile)
func (s *UserService) GetUserByID(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}
