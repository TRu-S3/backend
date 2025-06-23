package interfaces

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(r *gin.Engine, fileHandler *FileHandler) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// File routes
		files := v1.Group("/files")
		{
			files.POST("", fileHandler.CreateFile)               // Create file
			files.GET("", fileHandler.ListFiles)                 // List files
			files.GET("/:id", fileHandler.GetFile)               // Get file metadata
			files.GET("/:id/download", fileHandler.DownloadFile) // Download file content
			files.PUT("/:id", fileHandler.UpdateFile)            // Update file
			files.DELETE("/:id", fileHandler.DeleteFile)         // Delete file
		}
	}
}
