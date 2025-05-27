package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "Validation failed",
		Error:   "Validation errors: " + joinErrors(errors),
	})
}

// joinErrors joins multiple error messages
func joinErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}

	result := errors[0]
	for i := 1; i < len(errors); i++ {
		result += ", " + errors[i]
	}
	return result
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *gin.Context, data interface{}, page, limit int, total int64) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Data retrieved successfully",
		Data: PaginationResponse{
			Data:       data,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
