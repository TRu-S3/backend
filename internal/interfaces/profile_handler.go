package interfaces

import (
	"net/http"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/TRu-S3/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	*BaseHandler
}

func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

type CreateProfileRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	TagID    *uint  `json:"tag_id,omitempty"`
	Bio      string `json:"bio"`
	Age      *int   `json:"age,omitempty"`
	Location string `json:"location"`
}

type UpdateProfileRequest struct {
	TagID    *uint   `json:"tag_id,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Age      *int    `json:"age,omitempty"`
	Location *string `json:"location,omitempty"`
}

// CreateProfile handles POST /api/v1/profiles
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var req CreateProfileRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Validate age if provided
	if req.Age != nil && !utils.ValidatePositiveInt(c, *req.Age, "age") {
		return
	}

	// Validate bio length
	if !utils.ValidateMaxLength(c, req.Bio, 1000, "bio") {
		return
	}

	// Validate location length
	if !utils.ValidateMaxLength(c, req.Location, 100, "location") {
		return
	}

	// Check if user exists
	var user database.User
	if err := h.GetDatabase().First(&user, req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to validate user")
		return
	}

	// Check if profile already exists for this user
	var existingProfile database.Profile
	if err := h.GetDatabase().Where("user_id = ?", req.UserID).First(&existingProfile).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Profile already exists for this user")
		return
	}

	// Check if tag exists if provided
	if req.TagID != nil {
		var tag database.Tag
		if err := h.GetDatabase().First(&tag, *req.TagID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.ErrorResponse(c, http.StatusNotFound, "Tag not found")
				return
			}
			utils.InternalErrorResponse(c, "Failed to validate tag")
			return
		}
	}

	profile := database.Profile{
		UserID:   req.UserID,
		Bio:      req.Bio,
		Location: req.Location,
	}

	if req.TagID != nil {
		profile.TagID = *req.TagID
	}

	if req.Age != nil {
		profile.Age = *req.Age
	}

	if err := h.GetDatabase().Create(&profile).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to create profile")
		return
	}

	// Load relationships for response
	h.GetDatabase().Preload("User").Preload("Tag").First(&profile, profile.ID)

	h.HandleCreated(c, profile)
}

// ListProfiles handles GET /api/v1/profiles
func (h *ProfileHandler) ListProfiles(c *gin.Context) {
	var profiles []database.Profile

	// Parse pagination parameters
	params := utils.ParsePagination(c)

	query := h.GetDatabase().Model(&database.Profile{})

	// Filter by user_id if specified
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Filter by tag_id if specified
	if tagID := c.Query("tag_id"); tagID != "" {
		query = query.Where("tag_id = ?", tagID)
	}

	// Filter by location if specified
	if location := c.Query("location"); location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	// Filter by age range if specified
	if minAge := c.Query("min_age"); minAge != "" {
		query = query.Where("age >= ?", minAge)
	}
	if maxAge := c.Query("max_age"); maxAge != "" {
		query = query.Where("age <= ?", maxAge)
	}

	// Order by creation date (newest first)
	query = query.Order("created_at DESC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results with relationships
	if err := query.Offset(params.Offset).Limit(params.Limit).
		Preload("User").Preload("Tag").Find(&profiles).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to retrieve profiles")
		return
	}

	utils.StandardResponse(c, 200, gin.H{
		"profiles": profiles,
		"pagination": gin.H{
			"page":  params.Page,
			"limit": params.Limit,
			"total": total,
		},
	})
}

// GetProfile handles GET /api/v1/profiles/:id
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var profile database.Profile
	query := h.GetDatabase().Preload("User").Preload("Tag")
	if err := query.First(&profile, id).Error; err != nil {
		h.HandleDBError(c, err, "profile")
		return
	}

	utils.StandardResponse(c, 200, profile)
}

// GetProfileByUserID handles GET /api/v1/profiles/user/:user_id
func (h *ProfileHandler) GetProfileByUserID(c *gin.Context) {
	userID, ok := h.ParseIDParam(c, "user_id")
	if !ok {
		return
	}

	var profile database.Profile
	query := h.GetDatabase().Preload("User").Preload("Tag")
	if err := query.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		h.HandleDBError(c, err, "profile")
		return
	}

	utils.StandardResponse(c, 200, profile)
}

// UpdateProfile handles PUT /api/v1/profiles/:id
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var profile database.Profile
	var req UpdateProfileRequest

	if !h.BindJSON(c, &req) {
		return
	}

	// Find existing profile
	if err := h.GetDatabase().First(&profile, id).Error; err != nil {
		h.HandleDBError(c, err, "profile")
		return
	}

	// Update fields if provided
	if req.TagID != nil {
		// Check if tag exists
		var tag database.Tag
		if err := h.GetDatabase().First(&tag, *req.TagID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.ErrorResponse(c, http.StatusNotFound, "Tag not found")
				return
			}
			utils.InternalErrorResponse(c, "Failed to validate tag")
			return
		}
		profile.TagID = *req.TagID
	}

	if req.Bio != nil {
		if !utils.ValidateMaxLength(c, *req.Bio, 1000, "bio") {
			return
		}
		profile.Bio = *req.Bio
	}

	if req.Age != nil {
		if !utils.ValidatePositiveInt(c, *req.Age, "age") {
			return
		}
		profile.Age = *req.Age
	}

	if req.Location != nil {
		if !utils.ValidateMaxLength(c, *req.Location, 100, "location") {
			return
		}
		profile.Location = *req.Location
	}

	if err := h.GetDatabase().Save(&profile).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to update profile")
		return
	}

	// Load relationships for response
	h.GetDatabase().Preload("User").Preload("Tag").First(&profile, profile.ID)

	utils.StandardResponse(c, 200, profile)
}

// DeleteProfile handles DELETE /api/v1/profiles/:id
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var profile database.Profile

	if err := h.GetDatabase().First(&profile, id).Error; err != nil {
		h.HandleDBError(c, err, "profile")
		return
	}

	if err := h.GetDatabase().Delete(&profile).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to delete profile")
		return
	}

	utils.SuccessResponse(c, "Profile deleted successfully")
}