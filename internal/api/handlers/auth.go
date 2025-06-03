package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"social-media-app/internal/models"
	"social-media-app/internal/services"
	"social-media-app/internal/utils"
)

var authService *services.AuthService

// InitAuthHandler initializes the auth handler with the provided service
func InitAuthHandler(service *services.AuthService) {
	authService = service
}
func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	response, err := authService.Register(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}
func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := authService.Login(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User logged in successfully", response)
}

func Logout(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "User logged out successfully", gin.H{
		"message": "Please remove the token from client storage",
	})
}
