package interfaces

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContestHandler struct {
	db *gorm.DB
}

func NewContestHandler(db *gorm.DB) *ContestHandler {
	return &ContestHandler{db: db}
}

type CreateContestRequest struct {
	BackendQuota        int    `json:"backend_quota" binding:"min=0"`
	FrontendQuota       int    `json:"frontend_quota" binding:"min=0"`
	AIQuota             int    `json:"ai_quota" binding:"min=0"`
	ApplicationDeadline string `json:"application_deadline" binding:"required"`
	Purpose             string `json:"purpose" binding:"required"`
	Message             string `json:"message" binding:"required"`
	AuthorID            uint   `json:"author_id" binding:"required"`
}

type UpdateContestRequest struct {
	BackendQuota        *int    `json:"backend_quota,omitempty" binding:"omitempty,min=0"`
	FrontendQuota       *int    `json:"frontend_quota,omitempty" binding:"omitempty,min=0"`
	AIQuota             *int    `json:"ai_quota,omitempty" binding:"omitempty,min=0"`
	ApplicationDeadline *string `json:"application_deadline,omitempty"`
	Purpose             *string `json:"purpose,omitempty"`
	Message             *string `json:"message,omitempty"`
	Title               *string `json:"title,omitempty"`
	Description         *string `json:"description,omitempty"`
}

// CreateContest handles POST /api/v1/contests
func (h *ContestHandler) CreateContest(c *gin.Context) {
	var req CreateContestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse deadline
	deadline, err := time.Parse(time.RFC3339, req.ApplicationDeadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339 format (e.g., 2023-12-31T23:59:59Z)"})
		return
	}

	contest := database.Contest{
		BackendQuota:        req.BackendQuota,
		FrontendQuota:       req.FrontendQuota,
		AIQuota:             req.AIQuota,
		ApplicationDeadline: deadline,
		Purpose:             req.Purpose,
		Message:             req.Message,
		AuthorID:            req.AuthorID,
	}

	if err := h.db.Create(&contest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contest"})
		return
	}

	c.JSON(http.StatusCreated, contest)
}

// ListContests handles GET /api/v1/contests
func (h *ContestHandler) ListContests(c *gin.Context) {
	var contests []database.Contest
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	offset := (page - 1) * limit

	query := h.db.Model(&database.Contest{})

	// Filter by author if specified
	if authorID := c.Query("author_id"); authorID != "" {
		query = query.Where("author_id = ?", authorID)
	}

	// Filter by active contests (deadline not passed)
	if c.Query("active") == "true" {
		query = query.Where("application_deadline > ?", time.Now())
	}

	// Order by creation date (newest first)
	query = query.Order("created_at DESC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Find(&contests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contests": contests,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetContest handles GET /api/v1/contests/:id
func (h *ContestHandler) GetContest(c *gin.Context) {
	id := c.Param("id")
	var contest database.Contest

	if err := h.db.First(&contest, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contest"})
		}
		return
	}

	c.JSON(http.StatusOK, contest)
}

// UpdateContest handles PUT /api/v1/contests/:id
func (h *ContestHandler) UpdateContest(c *gin.Context) {
	id := c.Param("id")
	var contest database.Contest
	var req UpdateContestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find existing contest
	if err := h.db.First(&contest, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contest"})
		}
		return
	}

	// Update fields if provided
	if req.BackendQuota != nil {
		contest.BackendQuota = *req.BackendQuota
	}
	if req.FrontendQuota != nil {
		contest.FrontendQuota = *req.FrontendQuota
	}
	if req.AIQuota != nil {
		contest.AIQuota = *req.AIQuota
	}
	if req.ApplicationDeadline != nil {
		deadline, err := time.Parse(time.RFC3339, *req.ApplicationDeadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339 format"})
			return
		}
		contest.ApplicationDeadline = deadline
	}
	if req.Purpose != nil {
		contest.Purpose = *req.Purpose
	}
	if req.Message != nil {
		contest.Message = *req.Message
	}
	if req.Title != nil {
		contest.Title = *req.Title
	}
	if req.Description != nil {
		contest.Description = *req.Description
	}

	if err := h.db.Save(&contest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contest"})
		return
	}

	c.JSON(http.StatusOK, contest)
}

// DeleteContest handles DELETE /api/v1/contests/:id
func (h *ContestHandler) DeleteContest(c *gin.Context) {
	id := c.Param("id")
	var contest database.Contest

	if err := h.db.First(&contest, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contest"})
		}
		return
	}

	if err := h.db.Delete(&contest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contest deleted successfully"})
}