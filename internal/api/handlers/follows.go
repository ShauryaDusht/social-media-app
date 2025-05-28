package handlers

import (
	"net/http"

	"social-media-app/internal/utils"

	"github.com/gin-gonic/gin"
)

// Follows a user
func FollowUser(c *gin.Context) {
	// TODO
	utils.ErrorResponse(c, http.StatusNotImplemented, "Follow user endpoint not implemented yet")
}

// Unfollows a user
func UnfollowUser(c *gin.Context) {
	// TODO
	utils.ErrorResponse(c, http.StatusNotImplemented, "Unfollow user endpoint not implemented yet")
}

// Gets user's followers
func GetFollowers(c *gin.Context) {
	// TODO
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get followers endpoint not implemented yet")
}

// Gets whom the user follows
func GetFollowing(c *gin.Context) {
	// TODO
	utils.ErrorResponse(c, http.StatusNotImplemented, "Get following endpoint not implemented yet")
}
