package utils

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ValidateRequired checks if required fields are not empty
func ValidateRequired(c *gin.Context, fields map[string]string) bool {
	var missingFields []string
	
	for fieldName, value := range fields {
		if strings.TrimSpace(value) == "" {
			missingFields = append(missingFields, fieldName)
		}
	}
	
	if len(missingFields) > 0 {
		ErrorResponse(c, 400, "Missing required fields: "+strings.Join(missingFields, ", "))
		return false
	}
	
	return true
}

// ValidatePositiveInt checks if an integer is positive
func ValidatePositiveInt(c *gin.Context, value int, fieldName string) bool {
	if value <= 0 {
		ErrorResponse(c, 400, fieldName+" must be greater than 0")
		return false
	}
	return true
}

// ValidateEmail checks if email format is valid (basic validation)
func ValidateEmail(c *gin.Context, email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		ErrorResponse(c, 400, "Invalid email format")
		return false
	}
	return true
}

// ValidateDateRange checks if start date is before end date
func ValidateDateRange(c *gin.Context, startDate, endDate time.Time, startFieldName, endFieldName string) bool {
	if startDate.After(endDate) {
		ErrorResponse(c, 400, startFieldName+" must be before "+endFieldName)
		return false
	}
	return true
}

// ValidateFutureDate checks if date is in the future
func ValidateFutureDate(c *gin.Context, date time.Time, fieldName string) bool {
	if date.Before(time.Now()) {
		ErrorResponse(c, 400, fieldName+" must be in the future")
		return false
	}
	return true
}

// ValidateMaxLength checks if string doesn't exceed maximum length
func ValidateMaxLength(c *gin.Context, value string, maxLength int, fieldName string) bool {
	if len(value) > maxLength {
		ErrorResponse(c, 400, fieldName+" must not exceed "+string(rune(maxLength))+" characters")
		return false
	}
	return true
}