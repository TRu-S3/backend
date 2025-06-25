package interfaces

import (
	"net/http"
	"strconv"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookmarkHandler struct {
	db *gorm.DB
}

func NewBookmarkHandler(db *gorm.DB) *BookmarkHandler {
	return &BookmarkHandler{db: db}
}

type CreateBookmarkRequest struct {
	UserID           uint `json:"user_id" binding:"required"`
	BookmarkedUserID uint `json:"bookmarked_user_id" binding:"required"`
}

func (h *BookmarkHandler) CreateBookmark(c *gin.Context) {
	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID == req.BookmarkedUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot bookmark yourself"})
		return
	}

	bookmark := database.Bookmark{
		UserID:           req.UserID,
		BookmarkedUserID: req.BookmarkedUserID,
	}

	if err := h.db.Create(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
		return
	}

	c.JSON(http.StatusCreated, bookmark)
}

func (h *BookmarkHandler) ListBookmarks(c *gin.Context) {
	var bookmarks []database.Bookmark
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	query := h.db.Model(&database.Bookmark{})

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query = query.Order("created_at DESC")

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Preload("User").Preload("BookmarkedUser").Find(&bookmarks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookmarks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bookmarks": bookmarks,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *BookmarkHandler) DeleteBookmark(c *gin.Context) {
	id := c.Param("id")
	var bookmark database.Bookmark

	if err := h.db.First(&bookmark, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookmark"})
		}
		return
	}

	if err := h.db.Delete(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookmark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
}