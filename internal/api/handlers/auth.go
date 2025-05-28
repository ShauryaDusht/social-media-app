package handlers

import (
	"net/http"
	"time"

	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"social-media-app/internal/config"
	"social-media-app/internal/database"
	"social-media-app/internal/models"
)

// Register handles user registration
func Register(c *gin.Context) {
	var req models.RegisterRequest

	// JSON to RegisterRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Make new user model
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Save user
	if err := database.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate token upon successful registration
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, config.Load().JWT.Secret, time.Hour*24)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Send success response with user data and token
	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	})

}

// User Login
func Login(c *gin.Context) {
	req := models.LoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Find user with email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token upon successful login
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, config.Load().JWT.Secret, time.Hour*24)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Success response with user data and token
	utils.SuccessResponse(c, http.StatusOK, "User logged in successfully", models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	})
}
