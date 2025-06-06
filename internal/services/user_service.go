package services

import (
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
	"social-media-app/internal/utils"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) GetUserByID(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) UpdateProfile(userID uint, req models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

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
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// SearchUsers searches for users by name or username
func (s *UserService) SearchUsers(query string, limit, offset int) ([]models.UserResponse, error) {
	users, err := s.userRepo.SearchUsers(query, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, nil
}
