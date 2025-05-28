package handlers

import (
	"net/http"

	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
)

// Gets current user's profile
func GetProfile(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get profile endpoint not implemented yet")
}

// Updates current user's profile
func UpdateProfile(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Update profile endpoint not implemented yet")
}

// Gets user by ID
func GetUserByID(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get user by ID endpoint not implemented yet")
}
