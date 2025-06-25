package interfaces

import (
	"net/http"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/TRu-S3/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	*BaseHandler
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Gmail string `json:"gmail" binding:"required,email"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Gmail *string `json:"gmail,omitempty"`
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Validate email format
	if !utils.ValidateEmail(c, req.Gmail) {
		return
	}

	user := database.User{
		Name:  req.Name,
		Gmail: req.Gmail,
	}

	if err := h.GetDatabase().Create(&user).Error; err != nil {
		if isUniqueConstraintError(err, "gmail") {
			utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists")
			return
		}
		utils.InternalErrorResponse(c, "Failed to create user")
		return
	}

	h.HandleCreated(c, user)
}

// ListUsers handles GET /api/v1/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []database.User

	// Parse pagination parameters
	params := utils.ParsePagination(c)

	query := h.GetDatabase().Model(&database.User{})

	// Filter by name if specified
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Filter by email if specified
	if email := c.Query("gmail"); email != "" {
		query = query.Where("gmail ILIKE ?", "%"+email+"%")
	}

	// Order by creation date (newest first)
	query = query.Order("created_at DESC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results with profile
	if err := query.Offset(params.Offset).Limit(params.Limit).Preload("Profile").Find(&users).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to retrieve users")
		return
	}

	utils.StandardResponse(c, 200, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":  params.Page,
			"limit": params.Limit,
			"total": total,
		},
	})
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var user database.User
	query := h.GetDatabase().Preload("Profile").Preload("Profile.Tag").Preload("Bookmarks")
	if err := query.First(&user, id).Error; err != nil {
		h.HandleDBError(c, err, "user")
		return
	}

	utils.StandardResponse(c, 200, user)
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var user database.User
	var req UpdateUserRequest

	if !h.BindJSON(c, &req) {
		return
	}

	// Find existing user
	if err := h.GetDatabase().First(&user, id).Error; err != nil {
		h.HandleDBError(c, err, "user")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Gmail != nil {
		if !utils.ValidateEmail(c, *req.Gmail) {
			return
		}
		user.Gmail = *req.Gmail
	}

	if err := h.GetDatabase().Save(&user).Error; err != nil {
		if isUniqueConstraintError(err, "gmail") {
			utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists")
			return
		}
		utils.InternalErrorResponse(c, "Failed to update user")
		return
	}

	utils.StandardResponse(c, 200, user)
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var user database.User

	if err := h.GetDatabase().First(&user, id).Error; err != nil {
		h.HandleDBError(c, err, "user")
		return
	}

	if err := h.GetDatabase().Delete(&user).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to delete user")
		return
	}

	utils.SuccessResponse(c, "User deleted successfully")
}

// Helper function to check for unique constraint errors
func isUniqueConstraintError(err error, field string) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return (field == "gmail" && (contains(errStr, "duplicate key") || contains(errStr, "unique constraint"))) &&
		contains(errStr, field)
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
