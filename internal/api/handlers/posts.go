package handlers

import (
	"net/http"

	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
)

// TODO

// Gets all posts/timeline
func GetPosts(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get posts endpoint not implemented yet")
}

// Creates a new post
func CreatePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Create post endpoint not implemented yet")
}

// Gets a specific post
func GetPostByID(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get post by ID endpoint not implemented yet")
}

// Updates a post
func UpdatePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Update post endpoint not implemented yet")
}

// Deletes a post
func DeletePost(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Delete post endpoint not implemented yet")
}

// Gets posts by a specific user
func GetUserPosts(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get user posts endpoint not implemented yet")
}

// Gets personalized timeline
func GetTimeline(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get timeline endpoint not implemented yet")
}
