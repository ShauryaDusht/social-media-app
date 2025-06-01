package handlers

import (
	"net/http"
	"strconv"

	"social-media-app/internal/models"
	"social-media-app/internal/services"
	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
)

var postService *services.PostService

func InitPostHandler(service *services.PostService) {
	postService = service
}

func parseLimitOffset(c *gin.Context) (int, int) {
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
	return limit, offset
}

func GetPosts(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	limit, offset := parseLimitOffset(c)
	posts, err := postService.GetAll(limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Posts retrieved successfully", posts)
}

func CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var post models.CreatePostRequest
	if err := c.BindJSON(&post); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := postService.CreatePost(userID.(uint), post.Content, post.ImageURL)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Post created successfully", nil)
}

func GetPostByID(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := postService.GetByID(uint(postID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post retrieved successfully", post)
}

func UpdatePost(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.UpdatePostRequest
	if err := c.BindJSON(&post); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = postService.Update(&models.Post{ID: uint(postID), Content: post.Content, ImageURL: post.ImageURL})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post updated successfully", nil)
}

func DeletePost(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = postService.Delete(uint(postID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post deleted successfully", nil)
}

func GetUserPosts(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit, offset := parseLimitOffset(c)
	posts, err := postService.GetByUserID(uint(userID), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User posts retrieved successfully", posts)
}

func GetTimeline(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	limit, offset := parseLimitOffset(c)
	posts, err := postService.GetTimeline(userID.(uint), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Timeline retrieved successfully", posts)
}
