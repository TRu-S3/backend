package interfaces

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/TRu-S3/backend/internal/application"
	"github.com/TRu-S3/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// FileHandler handles HTTP requests for file operations
type FileHandler struct {
	fileService *application.FileService
}

// NewFileHandler creates a new FileHandler
func NewFileHandler(fileService *application.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// CreateFile handles file creation requests
// POST /files
func (h *FileHandler) CreateFile(c *gin.Context) {
	// Get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Get files from form
	files := form.File["file"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	file := files[0]

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Read file content
	content, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// Create file request
	req := &domain.CreateFileRequest{
		Name:        file.Filename,
		Content:     content,
		ContentType: file.Header.Get("Content-Type"),
	}

	// Create file through service
	createdFile, err := h.fileService.CreateFile(c.Request.Context(), req)
	if err != nil {
		if err == domain.ErrFileAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "File already exists"})
			return
		}
		if err == domain.ErrInvalidFileName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
			return
		}
		log.Printf("Failed to create file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}

	c.JSON(http.StatusCreated, createdFile)
}

// GetFile handles file retrieval requests
// GET /files/:id
func (h *FileHandler) GetFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrFileNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
		return
	}

	c.JSON(http.StatusOK, file)
}

// DownloadFile handles file download requests
// GET /files/:id/download
func (h *FileHandler) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	fileData, err := h.fileService.GetFileContent(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrFileNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file content"})
		return
	}

	// Set response headers
	c.Header("Content-Disposition", "attachment; filename="+fileData.File.Name)
	c.Header("Content-Type", fileData.File.ContentType)
	c.Header("Content-Length", strconv.FormatInt(fileData.File.Size, 10))

	// Write content
	c.Data(http.StatusOK, fileData.File.ContentType, fileData.Content)
}

// ListFiles handles file listing requests
// GET /files
func (h *FileHandler) ListFiles(c *gin.Context) {
	// Parse query parameters
	prefix := c.Query("prefix")
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	query := &domain.FileQuery{
		Prefix: prefix,
		Limit:  limit,
		Offset: offset,
	}

	files, err := h.fileService.ListFiles(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
		"count": len(files),
	})
}

// UpdateFile handles file update requests
// PUT /files/:id
func (h *FileHandler) UpdateFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	// Get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	req := &domain.UpdateFileRequest{}

	// Get new name if provided
	if names, exists := form.Value["name"]; exists && len(names) > 0 {
		req.Name = names[0]
	}

	// Get new file if provided
	if files, exists := form.File["file"]; exists && len(files) > 0 {
		file := files[0]

		// Open the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		// Read file content
		content, err := io.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
			return
		}

		req.Content = content
		req.ContentType = file.Header.Get("Content-Type")
	}

	// Update file through service
	updatedFile, err := h.fileService.UpdateFile(c.Request.Context(), id, req)
	if err != nil {
		if err == domain.ErrFileNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		if err == domain.ErrInvalidFileName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
		return
	}

	c.JSON(http.StatusOK, updatedFile)
}

// DeleteFile handles file deletion requests
// DELETE /files/:id
func (h *FileHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	err := h.fileService.DeleteFile(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrFileNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
