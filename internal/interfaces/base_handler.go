package interfaces

import (
	"strconv"

	"github.com/TRu-S3/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BaseHandler provides common handler functionality
type BaseHandler struct {
	db *gorm.DB
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(db *gorm.DB) *BaseHandler {
	return &BaseHandler{db: db}
}

// ParseIDParam extracts and validates ID parameter from URL
func (h *BaseHandler) ParseIDParam(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	if idStr == "" {
		utils.BadRequestResponse(c, paramName+" parameter is required")
		return 0, false
	}
	
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid "+paramName+" format")
		return 0, false
	}
	
	return uint(id), true
}

// BindJSON binds JSON request body and handles errors
func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		utils.BadRequestResponse(c, "Invalid JSON format: "+err.Error())
		return false
	}
	return true
}

// GetDatabase returns the database instance
func (h *BaseHandler) GetDatabase() *gorm.DB {
	return h.db
}

// GetWithPagination performs paginated database query
func (h *BaseHandler) GetWithPagination(c *gin.Context, model interface{}, result interface{}) error {
	params := utils.ParsePagination(c)
	return h.db.Offset(params.Offset).Limit(params.Limit).Find(result).Error
}

// HandleNotFound handles record not found cases consistently
func (h *BaseHandler) HandleNotFound(c *gin.Context, entityName string) {
	utils.NotFoundResponse(c, entityName)
}

// HandleDBError handles database errors consistently
func (h *BaseHandler) HandleDBError(c *gin.Context, err error, entityName string) {
	utils.HandleDBError(c, err, entityName)
}

// HandleSuccess handles successful operations consistently
func (h *BaseHandler) HandleSuccess(c *gin.Context, message string) {
	utils.SuccessResponse(c, message)
}

// HandleCreated handles successful creation consistently
func (h *BaseHandler) HandleCreated(c *gin.Context, data interface{}) {
	utils.CreatedResponse(c, data)
}