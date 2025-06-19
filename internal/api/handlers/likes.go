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

	likedBy := userID.(uint)

	// Check if already liked
	isLiked, err := like_service.IsPostLikedByUser(likeRequest.PostID, likedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error checking like status")
		return
	}

	if isLiked {
		utils.ErrorResponse(c, http.StatusBadRequest, "Post already liked by this user")
		return
	}

	if err := like_service.LikePost(likedBy, likeRequest.PostID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Get updated like count
	likeCount, err := like_service.GetLikeCount(likeRequest.PostID)
	if err != nil {
		likeCount = 0
	}

	// Get list of users who liked this post
	likedUserIDs, err := like_service.GetLikedUserIDs(likeRequest.PostID)
	if err != nil {
		likedUserIDs = []uint{}
	}

	utils.SuccessResponse(c, http.StatusOK, "Post liked successfully", gin.H{
		"post_id":    likeRequest.PostID,
		"like_count": likeCount,
		"liked_by":   likedUserIDs,
		"is_liked":   true,
	})
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

	likedBy := userID.(uint)

	// Check if actually liked
	isLiked, err := like_service.IsPostLikedByUser(uint(postID), likedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error checking like status")
		return
	}

	if !isLiked {
		utils.ErrorResponse(c, http.StatusBadRequest, "Post not liked by this user")
		return
	}

	err = like_service.UnlikePost(likedBy, uint(postID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Get updated like count
	likeCount, err := like_service.GetLikeCount(uint(postID))
	if err != nil {
		likeCount = 0
	}

	// Get list of users who liked this post
	likedUserIDs, err := like_service.GetLikedUserIDs(uint(postID))
	if err != nil {
		likedUserIDs = []uint{}
	}

	utils.SuccessResponse(c, http.StatusOK, "Post unliked successfully", gin.H{
		"post_id":    postID,
		"like_count": likeCount,
		"liked_by":   likedUserIDs,
		"is_liked":   false,
	})
}

// GetPostLikes returns all users who liked a specific post
func GetPostLikes(c *gin.Context) {
	postIDParam := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	likes, err := like_service.GetPostLikes(uint(postID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error fetching likes")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Likes retrieved successfully", gin.H{
		"likes": likes,
		"count": len(likes),
	})
}
