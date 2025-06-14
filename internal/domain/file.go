package domain

import (
	"time"
)

// File represents a file entity in the domain
type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FileData represents file content with metadata
type FileData struct {
	File    *File  `json:"file"`
	Content []byte `json:"content"`
}

// CreateFileRequest represents a request to create a file
type CreateFileRequest struct {
	Name        string `json:"name" binding:"required"`
	Content     []byte `json:"content" binding:"required"`
	ContentType string `json:"content_type"`
}

// UpdateFileRequest represents a request to update a file
type UpdateFileRequest struct {
	Name        string `json:"name"`
	Content     []byte `json:"content"`
	ContentType string `json:"content_type"`
}

// FileQuery represents query parameters for listing files
type FileQuery struct {
	Prefix string `json:"prefix"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
