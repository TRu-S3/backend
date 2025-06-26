package interfaces

import (
	"net/http"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/TRu-S3/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TagHandler struct {
	*BaseHandler
}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagRequest struct {
	Name *string `json:"name,omitempty"`
}

// CreateTag handles POST /api/v1/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Validate tag name length
	if !utils.ValidateMaxLength(c, req.Name, 50, "name") {
		return
	}

	tag := database.Tag{
		Name: req.Name,
	}

	if err := h.GetDatabase().Create(&tag).Error; err != nil {
		if isUniqueConstraintError(err, "name") {
			utils.ErrorResponse(c, http.StatusConflict, "Tag with this name already exists")
			return
		}
		utils.InternalErrorResponse(c, "Failed to create tag")
		return
	}

	h.HandleCreated(c, tag)
}

// ListTags handles GET /api/v1/tags
func (h *TagHandler) ListTags(c *gin.Context) {
	var tags []database.Tag

	// Parse pagination parameters
	params := utils.ParsePagination(c)

	query := h.GetDatabase().Model(&database.Tag{})

	// Filter by name if specified
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Order by name
	query = query.Order("name ASC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results
	if err := query.Offset(params.Offset).Limit(params.Limit).Find(&tags).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to retrieve tags")
		return
	}

	utils.StandardResponse(c, 200, gin.H{
		"tags": tags,
		"pagination": gin.H{
			"page":  params.Page,
			"limit": params.Limit,
			"total": total,
		},
	})
}

// GetTag handles GET /api/v1/tags/:id
func (h *TagHandler) GetTag(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var tag database.Tag
	if err := h.GetDatabase().First(&tag, id).Error; err != nil {
		h.HandleDBError(c, err, "tag")
		return
	}

	utils.StandardResponse(c, 200, tag)
}

// UpdateTag handles PUT /api/v1/tags/:id
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var tag database.Tag
	var req UpdateTagRequest

	if !h.BindJSON(c, &req) {
		return
	}

	// Find existing tag
	if err := h.GetDatabase().First(&tag, id).Error; err != nil {
		h.HandleDBError(c, err, "tag")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		if !utils.ValidateMaxLength(c, *req.Name, 50, "name") {
			return
		}
		tag.Name = *req.Name
	}

	if err := h.GetDatabase().Save(&tag).Error; err != nil {
		if isUniqueConstraintError(err, "name") {
			utils.ErrorResponse(c, http.StatusConflict, "Tag with this name already exists")
			return
		}
		utils.InternalErrorResponse(c, "Failed to update tag")
		return
	}

	utils.StandardResponse(c, 200, tag)
}

// DeleteTag handles DELETE /api/v1/tags/:id
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var tag database.Tag

	if err := h.GetDatabase().First(&tag, id).Error; err != nil {
		h.HandleDBError(c, err, "tag")
		return
	}

	// Check if tag is being used by any profiles
	var profileCount int64
	h.GetDatabase().Model(&database.Profile{}).Where("tag_id = ?", id).Count(&profileCount)
	
	if profileCount > 0 {
		utils.ErrorResponse(c, http.StatusConflict, "Cannot delete tag: it is being used by profiles")
		return
	}

	if err := h.GetDatabase().Delete(&tag).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to delete tag")
		return
	}

	utils.SuccessResponse(c, "Tag deleted successfully")
}