package domain

import (
	"context"
	"errors"
)

var (
	ErrFileNotFound      = errors.New("file not found")
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrInvalidFileName   = errors.New("invalid file name")
	ErrInvalidFileSize   = errors.New("invalid file size")
)

// FileRepository defines the contract for file storage operations
type FileRepository interface {
	// Create creates a new file
	Create(ctx context.Context, req *CreateFileRequest) (*File, error)

	// GetByID retrieves a file by its ID (path)
	GetByID(ctx context.Context, id string) (*File, error)

	// GetContent retrieves file content with metadata
	GetContent(ctx context.Context, id string) (*FileData, error)

	// List retrieves a list of files based on query parameters
	List(ctx context.Context, query *FileQuery) ([]*File, error)

	// Update updates an existing file
	Update(ctx context.Context, id string, req *UpdateFileRequest) (*File, error)

	// Delete deletes a file by its ID
	Delete(ctx context.Context, id string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, id string) (bool, error)
}
