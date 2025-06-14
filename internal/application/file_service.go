package application

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/TRu-S3/backend/internal/domain"
)

// FileService represents the application service for file operations
type FileService struct {
	fileRepo domain.FileRepository
}

// NewFileService creates a new FileService
func NewFileService(fileRepo domain.FileRepository) *FileService {
	return &FileService{
		fileRepo: fileRepo,
	}
}

// CreateFile creates a new file
func (s *FileService) CreateFile(ctx context.Context, req *domain.CreateFileRequest) (*domain.File, error) {
	// Validate file name
	if err := s.validateFileName(req.Name); err != nil {
		return nil, err
	}

	// Set default content type if not provided
	if req.ContentType == "" {
		req.ContentType = s.detectContentType(req.Name)
	}

	// Check if file already exists
	exists, err := s.fileRepo.Exists(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check file existence: %w", err)
	}
	if exists {
		return nil, domain.ErrFileAlreadyExists
	}

	// Create the file
	file, err := s.fileRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return file, nil
}

// GetFile retrieves a file by its ID
func (s *FileService) GetFile(ctx context.Context, id string) (*domain.File, error) {
	if id == "" {
		return nil, domain.ErrInvalidFileName
	}

	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return file, nil
}

// GetFileContent retrieves file content with metadata
func (s *FileService) GetFileContent(ctx context.Context, id string) (*domain.FileData, error) {
	if id == "" {
		return nil, domain.ErrInvalidFileName
	}

	fileData, err := s.fileRepo.GetContent(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	return fileData, nil
}

// ListFiles retrieves a list of files based on query parameters
func (s *FileService) ListFiles(ctx context.Context, query *domain.FileQuery) ([]*domain.File, error) {
	// Set default limits
	if query.Limit <= 0 {
		query.Limit = 100
	}
	if query.Limit > 1000 {
		query.Limit = 1000
	}

	files, err := s.fileRepo.List(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// UpdateFile updates an existing file
func (s *FileService) UpdateFile(ctx context.Context, id string, req *domain.UpdateFileRequest) (*domain.File, error) {
	if id == "" {
		return nil, domain.ErrInvalidFileName
	}

	// Validate new file name if provided
	if req.Name != "" {
		if err := s.validateFileName(req.Name); err != nil {
			return nil, err
		}
	}

	// Set default content type if not provided but content is provided
	if req.Content != nil && req.ContentType == "" {
		if req.Name != "" {
			req.ContentType = s.detectContentType(req.Name)
		} else {
			// Get existing file to determine content type
			existingFile, err := s.fileRepo.GetByID(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("failed to get existing file: %w", err)
			}
			req.ContentType = s.detectContentType(existingFile.Name)
		}
	}

	// Update the file
	file, err := s.fileRepo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	return file, nil
}

// DeleteFile deletes a file by its ID
func (s *FileService) DeleteFile(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidFileName
	}

	err := s.fileRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// validateFileName validates the file name
func (s *FileService) validateFileName(name string) error {
	if name == "" {
		return domain.ErrInvalidFileName
	}

	// Check for invalid characters
	invalidChars := []string{"..", "//", "\\", "<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return domain.ErrInvalidFileName
		}
	}

	// Check if it's a valid path
	if !filepath.IsAbs(name) && strings.HasPrefix(name, "/") {
		return domain.ErrInvalidFileName
	}

	return nil
}

// detectContentType detects content type based on file extension
func (s *FileService) detectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	contentTypes := map[string]string{
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".mp4":  "video/mp4",
		".zip":  "application/zip",
		".gz":   "application/gzip",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}
