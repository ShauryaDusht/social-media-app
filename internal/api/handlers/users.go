package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"social-media-app/internal/models"
	"social-media-app/internal/services"
	"social-media-app/internal/utils"
)

var userService *services.UserService

// InitUserHandler initializes the user handler with the provided service
func InitUserHandler(service *services.UserService) {
	userService = service
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	profile, err := userService.GetProfile(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", profile)
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	profile, err := userService.UpdateProfile(userID.(uint), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", profile)
}

func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	profile, err := userService.GetUserByID(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", profile)
}

func SearchUsers(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Search query is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := userService.SearchUsers(query, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search users")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
}
