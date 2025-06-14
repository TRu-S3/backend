package infrastructure

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/TRu-S3/backend/internal/domain"
	"google.golang.org/api/iterator"
)

const (
	DefaultBucketName = "202506-zenn-ai-agent-hackathon"
	DefaultFolder     = "test"
)

// GCSFileRepository implements domain.FileRepository using Google Cloud Storage
type GCSFileRepository struct {
	client     *storage.Client
	bucketName string
	folder     string
}

// NewGCSFileRepository creates a new GCSFileRepository
func NewGCSFileRepository(client *storage.Client, bucketName, folder string) *GCSFileRepository {
	if bucketName == "" {
		bucketName = DefaultBucketName
	}
	if folder == "" {
		folder = DefaultFolder
	}

	return &GCSFileRepository{
		client:     client,
		bucketName: bucketName,
		folder:     folder,
	}
}

// Create creates a new file in GCS
func (r *GCSFileRepository) Create(ctx context.Context, req *domain.CreateFileRequest) (*domain.File, error) {
	objectName := r.getObjectName(req.Name)
	obj := r.client.Bucket(r.bucketName).Object(objectName)

	// Create the object writer
	w := obj.NewWriter(ctx)
	w.ContentType = req.ContentType
	w.Metadata = map[string]string{
		"original_name": req.Name,
		"created_at":    time.Now().Format(time.RFC3339),
	}

	// Write the content
	if _, err := w.Write(req.Content); err != nil {
		w.Close()
		return nil, fmt.Errorf("failed to write content: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Get the created object attributes
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	return r.attributesToFile(attrs), nil
}

// GetByID retrieves a file by its ID (object name)
func (r *GCSFileRepository) GetByID(ctx context.Context, id string) (*domain.File, error) {
	objectName := r.getObjectName(id)
	obj := r.client.Bucket(r.bucketName).Object(objectName)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, domain.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	return r.attributesToFile(attrs), nil
}

// GetContent retrieves file content with metadata
func (r *GCSFileRepository) GetContent(ctx context.Context, id string) (*domain.FileData, error) {
	objectName := r.getObjectName(id)
	obj := r.client.Bucket(r.bucketName).Object(objectName)

	// Get object attributes
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, domain.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	// Read object content
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	file := r.attributesToFile(attrs)
	return &domain.FileData{
		File:    file,
		Content: content,
	}, nil
}

// List retrieves a list of files based on query parameters
func (r *GCSFileRepository) List(ctx context.Context, query *domain.FileQuery) ([]*domain.File, error) {
	prefix := r.folder + "/"
	if query.Prefix != "" {
		prefix = r.getObjectName(query.Prefix)
	}

	it := r.client.Bucket(r.bucketName).Objects(ctx, &storage.Query{
		Prefix: prefix,
	})

	var files []*domain.File
	count := 0
	skipped := 0

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate objects: %w", err)
		}

		// Skip directories (objects ending with /)
		if attrs.Name[len(attrs.Name)-1] == '/' {
			continue
		}

		// Apply offset
		if skipped < query.Offset {
			skipped++
			continue
		}

		// Apply limit
		if query.Limit > 0 && count >= query.Limit {
			break
		}

		files = append(files, r.attributesToFile(attrs))
		count++
	}

	return files, nil
}

// Update updates an existing file
func (r *GCSFileRepository) Update(ctx context.Context, id string, req *domain.UpdateFileRequest) (*domain.File, error) {
	objectName := r.getObjectName(id)
	obj := r.client.Bucket(r.bucketName).Object(objectName)

	// Check if object exists
	_, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, domain.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	// Handle rename if new name is provided
	var newObj *storage.ObjectHandle
	if req.Name != "" && req.Name != id {
		newObjectName := r.getObjectName(req.Name)
		newObj = r.client.Bucket(r.bucketName).Object(newObjectName)
	} else {
		newObj = obj
	}

	// If content is provided, update the content
	if req.Content != nil {
		w := newObj.NewWriter(ctx)
		if req.ContentType != "" {
			w.ContentType = req.ContentType
		}
		w.Metadata = map[string]string{
			"updated_at": time.Now().Format(time.RFC3339),
		}
		if req.Name != "" {
			w.Metadata["original_name"] = req.Name
		}

		if _, err := w.Write(req.Content); err != nil {
			w.Close()
			return nil, fmt.Errorf("failed to write content: %w", err)
		}

		if err := w.Close(); err != nil {
			return nil, fmt.Errorf("failed to close writer: %w", err)
		}

		// Delete old object if renamed
		if newObj != obj {
			if err := obj.Delete(ctx); err != nil {
				return nil, fmt.Errorf("failed to delete old object: %w", err)
			}
		}
	} else if newObj != obj {
		// Only rename without content update
		copier := newObj.CopierFrom(obj)
		if _, err := copier.Run(ctx); err != nil {
			return nil, fmt.Errorf("failed to copy object: %w", err)
		}

		if err := obj.Delete(ctx); err != nil {
			return nil, fmt.Errorf("failed to delete old object: %w", err)
		}
	}

	// Get updated object attributes
	attrs, err := newObj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated object attributes: %w", err)
	}

	return r.attributesToFile(attrs), nil
}

// Delete deletes a file by its ID
func (r *GCSFileRepository) Delete(ctx context.Context, id string) error {
	objectName := r.getObjectName(id)
	obj := r.client.Bucket(r.bucketName).Object(objectName)

	err := obj.Delete(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return domain.ErrFileNotFound
		}
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// Exists checks if a file exists
func (r *GCSFileRepository) Exists(ctx context.Context, id string) (bool, error) {
	objectName := r.getObjectName(id)
	obj := r.client.Bucket(r.bucketName).Object(objectName)

	_, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		// Check for 404 error in the error message as well
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "notFound") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

// getObjectName constructs the full object name with folder prefix
func (r *GCSFileRepository) getObjectName(filename string) string {
	if r.folder == "" {
		return filename
	}
	return fmt.Sprintf("%s/%s", r.folder, filename)
}

// getFileNameFromObjectName extracts the filename from the full object name
func (r *GCSFileRepository) getFileNameFromObjectName(objectName string) string {
	if r.folder == "" {
		return objectName
	}
	prefix := r.folder + "/"
	if len(objectName) > len(prefix) && objectName[:len(prefix)] == prefix {
		return objectName[len(prefix):]
	}
	return objectName
}

// attributesToFile converts GCS object attributes to domain File
func (r *GCSFileRepository) attributesToFile(attrs *storage.ObjectAttrs) *domain.File {
	originalName := r.getFileNameFromObjectName(attrs.Name)
	if name, exists := attrs.Metadata["original_name"]; exists {
		originalName = name
	}

	return &domain.File{
		ID:          r.getFileNameFromObjectName(attrs.Name),
		Name:        originalName,
		Path:        attrs.Name,
		Size:        attrs.Size,
		ContentType: attrs.ContentType,
		CreatedAt:   attrs.Created,
		UpdatedAt:   attrs.Updated,
	}
}
