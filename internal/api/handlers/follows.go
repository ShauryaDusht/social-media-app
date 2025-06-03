package handlers

import (
	"net/http"
	"strconv"

	"social-media-app/internal/services"
	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
)

var followService *services.FollowService

func InitFollowHandler(s *services.FollowService) {
	followService = s
}
func FollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	targetIDParam := c.Param("id")
	targetID, err := strconv.ParseUint(targetIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	err = followService.FollowUser(userID.(uint), uint(targetID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User followed successfully", nil)
}

func UnfollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	targetIDParam := c.Param("id")
	targetID, err := strconv.ParseUint(targetIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	err = followService.UnfollowUser(userID.(uint), uint(targetID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User unfollowed successfully", nil)
}

func GetFollowers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	followers, err := followService.GetFollowers(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Followers retrieved successfully", followers)
}

func GetFollowing(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	following, err := followService.GetFollowing(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Following retrieved successfully", following)
}
