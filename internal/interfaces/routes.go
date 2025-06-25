package interfaces

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(r *gin.Engine, fileHandler *FileHandler, contestHandler *ContestHandler, bookmarkHandler *BookmarkHandler, hackathonHandler *HackathonHandler) {
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

		// Contest routes
		contests := v1.Group("/contests")
		{
			contests.POST("", contestHandler.CreateContest)   // Create contest
			contests.GET("", contestHandler.ListContests)     // List contests
			contests.GET("/:id", contestHandler.GetContest)   // Get contest by ID
			contests.PUT("/:id", contestHandler.UpdateContest) // Update contest
			contests.DELETE("/:id", contestHandler.DeleteContest) // Delete contest
		}

		// Bookmark routes
		bookmarks := v1.Group("/bookmarks")
		{
			bookmarks.POST("", bookmarkHandler.CreateBookmark)   // Create bookmark
			bookmarks.GET("", bookmarkHandler.ListBookmarks)     // List bookmarks
			bookmarks.DELETE("/:id", bookmarkHandler.DeleteBookmark) // Delete bookmark
		}

		// Hackathon routes
		hackathons := v1.Group("/hackathons")
		{
			hackathons.POST("", hackathonHandler.CreateHackathon)   // Create hackathon
			hackathons.GET("", hackathonHandler.ListHackathons)     // List hackathons
			hackathons.GET("/:id", hackathonHandler.GetHackathon)   // Get hackathon by ID
			hackathons.PUT("/:id", hackathonHandler.UpdateHackathon) // Update hackathon
			hackathons.DELETE("/:id", hackathonHandler.DeleteHackathon) // Delete hackathon
			
			// Participant routes
			hackathons.POST("/:id/participants", hackathonHandler.CreateParticipant)   // Register for hackathon
			hackathons.GET("/:id/participants", hackathonHandler.ListParticipants)     // List participants
			hackathons.DELETE("/:id/participants/:participant_id", hackathonHandler.DeleteParticipant) // Remove participant
		}
	}
}
