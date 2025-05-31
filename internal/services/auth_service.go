package services

import (
	"errors"
	"social-media-app/internal/config"
	"social-media-app/internal/models"
	"social-media-app/internal/repository"
	"social-media-app/internal/utils"
	"strings"
)

type AuthService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register registers a new user
func (s *AuthService) Register(req models.RegisterRequest) (*models.LoginResponse, error) {
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	emailExists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, errors.New("failed to check email existence")
	}
	if emailExists {
		return nil, errors.New("email already exists")
	}

	usernameExists, err := s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, errors.New("failed to check username existence")
	}
	if usernameExists {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
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
		return nil, errors.New("failed to create user")
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		s.config.JWT.Secret,
		s.config.JWT.Expiry,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	if err := s.validateLoginRequest(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		s.config.JWT.Secret,
		s.config.JWT.Expiry,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return utils.ValidateToken(tokenString, s.config.JWT.Secret)
}

// GetUserByID gets user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) validateRegisterRequest(req models.RegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if strings.TrimSpace(req.FirstName) == "" {
		return errors.New("first name is required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return errors.New("last name is required")
	}
	return nil
}

func (s *AuthService) validateLoginRequest(req models.LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}
	return nil
}
