package interfaces

import (
	"net/http"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/TRu-S3/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatchingHandler struct {
	*BaseHandler
}

func NewMatchingHandler(db *gorm.DB) *MatchingHandler {
	return &MatchingHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

type CreateMatchingRequest struct {
	User1ID uint   `json:"user1_id" binding:"required"`
	User2ID uint   `json:"user2_id" binding:"required"`
	Status  string `json:"status"`
}

type UpdateMatchingRequest struct {
	Status *string `json:"status,omitempty"`
}

// CreateMatching handles POST /api/v1/matchings
func (h *MatchingHandler) CreateMatching(c *gin.Context) {
	var req CreateMatchingRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Validate that user1_id and user2_id are different
	if req.User1ID == req.User2ID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot create matching with the same user")
		return
	}

	// Check if users exist
	var user1, user2 database.User
	if err := h.GetDatabase().First(&user1, req.User1ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User1 not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to validate user1")
		return
	}

	if err := h.GetDatabase().First(&user2, req.User2ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User2 not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to validate user2")
		return
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "pending"
	}

	// Validate status
	validStatuses := []string{"pending", "accepted", "rejected", "blocked"}
	if !isValidStatus(req.Status, validStatuses) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status. Must be one of: pending, accepted, rejected, blocked")
		return
	}

	// Check if matching already exists (in either direction)
	var existingMatching database.Matching
	if err := h.GetDatabase().Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		req.User1ID, req.User2ID, req.User2ID, req.User1ID,
	).First(&existingMatching).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Matching already exists between these users")
		return
	}

	matching := database.Matching{
		User1ID: req.User1ID,
		User2ID: req.User2ID,
		Status:  req.Status,
	}

	if err := h.GetDatabase().Create(&matching).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to create matching")
		return
	}

	// Load relationships for response
	h.GetDatabase().Preload("User1").Preload("User2").First(&matching, matching.ID)

	h.HandleCreated(c, matching)
}

// ListMatchings handles GET /api/v1/matchings
func (h *MatchingHandler) ListMatchings(c *gin.Context) {
	var matchings []database.Matching

	// Parse pagination parameters
	params := utils.ParsePagination(c)

	query := h.GetDatabase().Model(&database.Matching{})

	// Filter by user_id (either user1_id or user2_id)
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user1_id = ? OR user2_id = ?", userID, userID)
	}

	// Filter by user1_id if specified
	if user1ID := c.Query("user1_id"); user1ID != "" {
		query = query.Where("user1_id = ?", user1ID)
	}

	// Filter by user2_id if specified
	if user2ID := c.Query("user2_id"); user2ID != "" {
		query = query.Where("user2_id = ?", user2ID)
	}

	// Filter by status if specified
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Order by creation date (newest first)
	query = query.Order("created_at DESC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results with relationships
	if err := query.Offset(params.Offset).Limit(params.Limit).
		Preload("User1").Preload("User2").Find(&matchings).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to retrieve matchings")
		return
	}

	utils.StandardResponse(c, 200, gin.H{
		"matchings": matchings,
		"pagination": gin.H{
			"page":  params.Page,
			"limit": params.Limit,
			"total": total,
		},
	})
}

// GetMatching handles GET /api/v1/matchings/:id
func (h *MatchingHandler) GetMatching(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var matching database.Matching
	query := h.GetDatabase().Preload("User1").Preload("User2")
	if err := query.First(&matching, id).Error; err != nil {
		h.HandleDBError(c, err, "matching")
		return
	}

	utils.StandardResponse(c, 200, matching)
}

// UpdateMatching handles PUT /api/v1/matchings/:id
func (h *MatchingHandler) UpdateMatching(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var matching database.Matching
	var req UpdateMatchingRequest

	if !h.BindJSON(c, &req) {
		return
	}

	// Find existing matching
	if err := h.GetDatabase().First(&matching, id).Error; err != nil {
		h.HandleDBError(c, err, "matching")
		return
	}

	// Update fields if provided
	if req.Status != nil {
		// Validate status
		validStatuses := []string{"pending", "accepted", "rejected", "blocked"}
		if !isValidStatus(*req.Status, validStatuses) {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status. Must be one of: pending, accepted, rejected, blocked")
			return
		}
		matching.Status = *req.Status
	}

	if err := h.GetDatabase().Save(&matching).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to update matching")
		return
	}

	// Load relationships for response
	h.GetDatabase().Preload("User1").Preload("User2").First(&matching, matching.ID)

	utils.StandardResponse(c, 200, matching)
}

// DeleteMatching handles DELETE /api/v1/matchings/:id
func (h *MatchingHandler) DeleteMatching(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var matching database.Matching

	if err := h.GetDatabase().First(&matching, id).Error; err != nil {
		h.HandleDBError(c, err, "matching")
		return
	}

	if err := h.GetDatabase().Delete(&matching).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to delete matching")
		return
	}

	utils.SuccessResponse(c, "Matching deleted successfully")
}

// GetUserMatches handles GET /api/v1/users/:user_id/matches
func (h *MatchingHandler) GetUserMatches(c *gin.Context) {
	userID, ok := h.ParseIDParam(c, "user_id")
	if !ok {
		return
	}

	// Parse pagination parameters
	params := utils.ParsePagination(c)

	// Filter by status if specified
	status := c.DefaultQuery("status", "accepted")

	var matchings []database.Matching
	query := h.GetDatabase().Model(&database.Matching{}).
		Where("(user1_id = ? OR user2_id = ?) AND status = ?", userID, userID, status).
		Order("created_at DESC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results with relationships
	if err := query.Offset(params.Offset).Limit(params.Limit).
		Preload("User1").Preload("User2").Find(&matchings).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to retrieve user matches")
		return
	}

	utils.StandardResponse(c, 200, gin.H{
		"matches": matchings,
		"pagination": gin.H{
			"page":  params.Page,
			"limit": params.Limit,
			"total": total,
		},
	})
}

// Helper function to validate status
func isValidStatus(status string, validStatuses []string) bool {
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}