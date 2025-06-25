package interfaces

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/TRu-S3/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HackathonHandler struct {
	*BaseHandler
}

func NewHackathonHandler(db *gorm.DB) *HackathonHandler {
	return &HackathonHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

type CreateHackathonRequest struct {
	Name                 string    `json:"name" binding:"required"`
	Description          string    `json:"description"`
	StartDate            string    `json:"start_date" binding:"required"`
	EndDate              string    `json:"end_date" binding:"required"`
	RegistrationStart    string    `json:"registration_start" binding:"required"`
	RegistrationDeadline string    `json:"registration_deadline" binding:"required"`
	MaxParticipants      int       `json:"max_participants" binding:"min=0"`
	Location             string    `json:"location"`
	Organizer            string    `json:"organizer" binding:"required"`
	ContactEmail         string    `json:"contact_email"`
	PrizeInfo            string    `json:"prize_info"`
	Rules                string    `json:"rules"`
	TechStack            string    `json:"tech_stack"`
	IsPublic             *bool     `json:"is_public"`
	BannerURL            string    `json:"banner_url"`
	WebsiteURL           string    `json:"website_url"`
}

type UpdateHackathonRequest struct {
	Name                 *string `json:"name,omitempty"`
	Description          *string `json:"description,omitempty"`
	StartDate            *string `json:"start_date,omitempty"`
	EndDate              *string `json:"end_date,omitempty"`
	RegistrationStart    *string `json:"registration_start,omitempty"`
	RegistrationDeadline *string `json:"registration_deadline,omitempty"`
	MaxParticipants      *int    `json:"max_participants,omitempty" binding:"omitempty,min=0"`
	Location             *string `json:"location,omitempty"`
	Organizer            *string `json:"organizer,omitempty"`
	ContactEmail         *string `json:"contact_email,omitempty"`
	PrizeInfo            *string `json:"prize_info,omitempty"`
	Rules                *string `json:"rules,omitempty"`
	TechStack            *string `json:"tech_stack,omitempty"`
	Status               *string `json:"status,omitempty"`
	IsPublic             *bool   `json:"is_public,omitempty"`
	BannerURL            *string `json:"banner_url,omitempty"`
	WebsiteURL           *string `json:"website_url,omitempty"`
}

type CreateParticipantRequest struct {
	HackathonID uint   `json:"hackathon_id" binding:"required"`
	UserID      uint   `json:"user_id" binding:"required"`
	TeamName    string `json:"team_name"`
	Role        string `json:"role"`
	Notes       string `json:"notes"`
}

// CreateHackathon handles POST /api/v1/hackathons
func (h *HackathonHandler) CreateHackathon(c *gin.Context) {
	var req CreateHackathonRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Parse and validate dates using utility functions
	startDate, ok := utils.ParseDateRFC3339(c, req.StartDate, "start_date")
	if !ok {
		return
	}

	endDate, ok := utils.ParseDateRFC3339(c, req.EndDate, "end_date")
	if !ok {
		return
	}

	registrationStart, ok := utils.ParseDateRFC3339(c, req.RegistrationStart, "registration_start")
	if !ok {
		return
	}

	registrationDeadline, ok := utils.ParseDateRFC3339(c, req.RegistrationDeadline, "registration_deadline")
	if !ok {
		return
	}

	// Validate date ranges using utility functions
	if !utils.ValidateDateRange(c, startDate, endDate, "start_date", "end_date") {
		return
	}
	if !utils.ValidateDateRange(c, registrationStart, registrationDeadline, "registration_start", "registration_deadline") {
		return
	}
	if !utils.ValidateDateRange(c, registrationDeadline, startDate, "registration_deadline", "start_date") {
		return
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	hackathon := database.Hackathon{
		Name:                 req.Name,
		Description:          req.Description,
		StartDate:            startDate,
		EndDate:              endDate,
		RegistrationStart:    registrationStart,
		RegistrationDeadline: registrationDeadline,
		MaxParticipants:      req.MaxParticipants,
		Location:             req.Location,
		Organizer:            req.Organizer,
		ContactEmail:         req.ContactEmail,
		PrizeInfo:            req.PrizeInfo,
		Rules:                req.Rules,
		TechStack:            req.TechStack,
		Status:               "upcoming",
		IsPublic:             isPublic,
		BannerURL:            req.BannerURL,
		WebsiteURL:           req.WebsiteURL,
	}

	if err := h.GetDatabase().Create(&hackathon).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to create hackathon")
		return
	}

	h.HandleCreated(c, hackathon)
}

// ListHackathons handles GET /api/v1/hackathons
func (h *HackathonHandler) ListHackathons(c *gin.Context) {
	var hackathons []database.Hackathon
	
	// Parse pagination parameters using utility function
	params := utils.ParsePagination(c)

	query := h.GetDatabase().Model(&database.Hackathon{})

	// Filter by status if specified
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by organizer if specified
	if organizer := c.Query("organizer"); organizer != "" {
		query = query.Where("organizer ILIKE ?", "%"+organizer+"%")
	}

	// Filter by public/private
	if isPublic := c.Query("is_public"); isPublic != "" {
		query = query.Where("is_public = ?", isPublic == "true")
	}

	// Filter by upcoming events
	if c.Query("upcoming") == "true" {
		query = query.Where("start_date > ?", time.Now())
	}

	// Filter by ongoing events
	if c.Query("ongoing") == "true" {
		now := time.Now()
		query = query.Where("start_date <= ? AND end_date >= ?", now, now)
	}

	// Filter by registration open
	if c.Query("registration_open") == "true" {
		now := time.Now()
		query = query.Where("registration_start <= ? AND registration_deadline >= ?", now, now)
	}

	// Order by start date (newest first)
	query = query.Order("start_date DESC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results
	if err := query.Offset(params.Offset).Limit(params.Limit).Find(&hackathons).Error; err != nil {
		utils.InternalErrorResponse(c, "Failed to retrieve hackathons")
		return
	}

	utils.StandardResponse(c, 200, gin.H{
		"hackathons": hackathons,
		"pagination": gin.H{
			"page":  params.Page,
			"limit": params.Limit,
			"total": total,
		},
	})
}

// GetHackathon handles GET /api/v1/hackathons/:id
func (h *HackathonHandler) GetHackathon(c *gin.Context) {
	id, ok := h.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var hackathon database.Hackathon
	query := h.GetDatabase().Preload("Participants")
	if err := query.First(&hackathon, id).Error; err != nil {
		h.HandleDBError(c, err, "hackathon")
		return
	}

	utils.StandardResponse(c, 200, hackathon)
}

// UpdateHackathon handles PUT /api/v1/hackathons/:id
func (h *HackathonHandler) UpdateHackathon(c *gin.Context) {
	id := c.Param("id")
	var hackathon database.Hackathon
	var req UpdateHackathonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find existing hackathon
	if err := h.db.First(&hackathon, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hackathon not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hackathon"})
		}
		return
	}

	// Update fields if provided
	if req.Name != nil {
		hackathon.Name = *req.Name
	}
	if req.Description != nil {
		hackathon.Description = *req.Description
	}
	if req.StartDate != nil {
		startDate, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use RFC3339 format"})
			return
		}
		hackathon.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use RFC3339 format"})
			return
		}
		hackathon.EndDate = endDate
	}
	if req.RegistrationStart != nil {
		registrationStart, err := time.Parse(time.RFC3339, *req.RegistrationStart)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration_start format. Use RFC3339 format"})
			return
		}
		hackathon.RegistrationStart = registrationStart
	}
	if req.RegistrationDeadline != nil {
		registrationDeadline, err := time.Parse(time.RFC3339, *req.RegistrationDeadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration_deadline format. Use RFC3339 format"})
			return
		}
		hackathon.RegistrationDeadline = registrationDeadline
	}
	if req.MaxParticipants != nil {
		hackathon.MaxParticipants = *req.MaxParticipants
	}
	if req.Location != nil {
		hackathon.Location = *req.Location
	}
	if req.Organizer != nil {
		hackathon.Organizer = *req.Organizer
	}
	if req.ContactEmail != nil {
		hackathon.ContactEmail = *req.ContactEmail
	}
	if req.PrizeInfo != nil {
		hackathon.PrizeInfo = *req.PrizeInfo
	}
	if req.Rules != nil {
		hackathon.Rules = *req.Rules
	}
	if req.TechStack != nil {
		hackathon.TechStack = *req.TechStack
	}
	if req.Status != nil {
		// Validate status
		validStatuses := []string{"upcoming", "ongoing", "completed", "cancelled"}
		valid := false
		for _, status := range validStatuses {
			if *req.Status == status {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: upcoming, ongoing, completed, cancelled"})
			return
		}
		hackathon.Status = *req.Status
	}
	if req.IsPublic != nil {
		hackathon.IsPublic = *req.IsPublic
	}
	if req.BannerURL != nil {
		hackathon.BannerURL = *req.BannerURL
	}
	if req.WebsiteURL != nil {
		hackathon.WebsiteURL = *req.WebsiteURL
	}

	// Validate date logic after updates
	if hackathon.EndDate.Before(hackathon.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}
	if hackathon.RegistrationDeadline.After(hackathon.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration deadline must be before start date"})
		return
	}
	if hackathon.RegistrationStart.After(hackathon.RegistrationDeadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration start must be before registration deadline"})
		return
	}

	if err := h.db.Save(&hackathon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hackathon"})
		return
	}

	c.JSON(http.StatusOK, hackathon)
}

// DeleteHackathon handles DELETE /api/v1/hackathons/:id
func (h *HackathonHandler) DeleteHackathon(c *gin.Context) {
	id := c.Param("id")
	var hackathon database.Hackathon

	if err := h.db.First(&hackathon, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hackathon not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hackathon"})
		}
		return
	}

	if err := h.db.Delete(&hackathon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hackathon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hackathon deleted successfully"})
}

// CreateParticipant handles POST /api/v1/hackathons/:id/participants
func (h *HackathonHandler) CreateParticipant(c *gin.Context) {
	hackathonID := c.Param("id")
	var req CreateParticipantRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate hackathon exists
	var hackathon database.Hackathon
	if err := h.db.First(&hackathon, hackathonID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hackathon not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hackathon"})
		}
		return
	}

	// Check if registration is open
	now := time.Now()
	if now.Before(hackathon.RegistrationStart) || now.After(hackathon.RegistrationDeadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration is not open for this hackathon"})
		return
	}

	// Check if hackathon is full
	if hackathon.MaxParticipants > 0 {
		var participantCount int64
		h.db.Model(&database.HackathonParticipant{}).Where("hackathon_id = ? AND status != ?", hackathonID, "cancelled").Count(&participantCount)
		if int(participantCount) >= hackathon.MaxParticipants {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Hackathon is full"})
			return
		}
	}

	// Override hackathon_id from URL parameter
	hackathonIDParsed, _ := strconv.ParseUint(hackathonID, 10, 32)
	req.HackathonID = uint(hackathonIDParsed)

	participant := database.HackathonParticipant{
		HackathonID: uint(req.HackathonID),
		UserID:      req.UserID,
		TeamName:    req.TeamName,
		Role:        req.Role,
		Status:      "registered",
		Notes:       req.Notes,
	}

	if err := h.db.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create participant"})
		return
	}

	// Load relationships for response
	h.db.Preload("User").Preload("Hackathon").First(&participant, participant.ID)

	c.JSON(http.StatusCreated, participant)
}

// ListParticipants handles GET /api/v1/hackathons/:id/participants
func (h *HackathonHandler) ListParticipants(c *gin.Context) {
	hackathonID := c.Param("id")
	var participants []database.HackathonParticipant
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	query := h.db.Model(&database.HackathonParticipant{}).Where("hackathon_id = ?", hackathonID)

	// Filter by status if specified
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by team name if specified
	if teamName := c.Query("team_name"); teamName != "" {
		query = query.Where("team_name ILIKE ?", "%"+teamName+"%")
	}

	query = query.Order("registration_date DESC")

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Preload("User").Find(&participants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve participants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"participants": participants,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// DeleteParticipant handles DELETE /api/v1/hackathons/:id/participants/:participant_id
func (h *HackathonHandler) DeleteParticipant(c *gin.Context) {
	hackathonID := c.Param("id")
	participantID := c.Param("participant_id")
	
	var participant database.HackathonParticipant

	if err := h.db.Where("hackathon_id = ? AND id = ?", hackathonID, participantID).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Participant not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve participant"})
		}
		return
	}

	if err := h.db.Delete(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant deleted successfully"})
}