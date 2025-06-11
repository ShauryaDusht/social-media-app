package handlers

import (
	"net/http"
	"strconv"

	"social-media-app/internal/models"
	"social-media-app/internal/services"
	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
)

var like_service *services.LikeService

func InitLikeHandler(s *services.LikeService) {
	like_service = s
}
func LikePost(c *gin.Context) {
	var likeRequest models.LikeRequest
	if err := c.ShouldBindJSON(&likeRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := like_service.LikePost(userID.(uint), likeRequest.PostID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post liked successfully", nil)
}

func UnlikePost(c *gin.Context) {
	postIDParam := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err = like_service.UnlikePost(userID.(uint), uint(postID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post unliked successfully", nil)
}
