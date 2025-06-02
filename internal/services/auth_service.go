package services

import (
	"errors"

	"social-media-app/internal/config"
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
	"social-media-app/internal/utils"
)

type AuthService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.LoginResponse, error) {
	emailExists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("email already exists")
	}

	usernameExists, err := s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, s.config.JWT.Secret, s.config.JWT.Expiry)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, s.config.JWT.Secret, s.config.JWT.Expiry)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}
