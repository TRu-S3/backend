package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

// ParsePagination extracts pagination parameters from query string
func ParsePagination(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	// Ensure minimum values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	// Limit maximum page size
	if limit > 100 {
		limit = 100
	}
	
	offset := (page - 1) * limit
	
	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// HandleDBError provides consistent database error handling
func HandleDBError(c *gin.Context, err error, entityName string) {
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": entityName + " not found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve " + entityName})
	}
}

// ParseDateRFC3339 parses date string in RFC3339 format with error handling
func ParseDateRFC3339(c *gin.Context, dateStr, fieldName string) (time.Time, bool) {
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fieldName + " is required"})
		return time.Time{}, false
	}
	
	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid " + fieldName + " format. Use RFC3339 format (e.g., 2025-01-01T10:00:00Z)",
		})
		return time.Time{}, false
	}
	
	return parsedTime, true
}

// StandardResponse provides consistent JSON response formatting
func StandardResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// SuccessResponse provides consistent success response formatting
func SuccessResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// ErrorResponse provides consistent error response formatting
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// CreatedResponse provides consistent creation response formatting
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// BadRequestResponse provides consistent bad request response formatting
func BadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

// InternalErrorResponse provides consistent internal error response formatting
func InternalErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

// NotFoundResponse provides consistent not found response formatting
func NotFoundResponse(c *gin.Context, entityName string) {
	c.JSON(http.StatusNotFound, gin.H{"error": entityName + " not found"})
}