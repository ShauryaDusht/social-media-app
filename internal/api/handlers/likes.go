package handlers

import (
	"net/http"

	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
)

// Likes a post
func LikePost(c *gin.Context) {
	// TODO
	utils.ErrorResponse(c, http.StatusNotImplemented, "Like post endpoint not implemented yet")
}

// Unlikes a post
func UnlikePost(c *gin.Context) {
	// TODO
	utils.ErrorResponse(c, http.StatusNotImplemented, "Unlike post endpoint not implemented yet")
}
