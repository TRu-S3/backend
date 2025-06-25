package interfaces

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(r *gin.Engine, fileHandler *FileHandler, contestHandler *ContestHandler, bookmarkHandler *BookmarkHandler, hackathonHandler *HackathonHandler, userHandler *UserHandler, tagHandler *TagHandler, profileHandler *ProfileHandler, matchingHandler *MatchingHandler) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)                    // Create user
			users.GET("", userHandler.ListUsers)                      // List users
			users.GET("/:id", userHandler.GetUser)                    // Get user by ID
			users.PUT("/:id", userHandler.UpdateUser)                 // Update user
			users.DELETE("/:id", userHandler.DeleteUser)              // Delete user
			users.GET("/:user_id/matches", matchingHandler.GetUserMatches) // Get user matches
		}

		// Tag routes
		tags := v1.Group("/tags")
		{
			tags.POST("", tagHandler.CreateTag)        // Create tag
			tags.GET("", tagHandler.ListTags)          // List tags
			tags.GET("/:id", tagHandler.GetTag)        // Get tag by ID
			tags.PUT("/:id", tagHandler.UpdateTag)     // Update tag
			tags.DELETE("/:id", tagHandler.DeleteTag)  // Delete tag
		}

		// Profile routes
		profiles := v1.Group("/profiles")
		{
			profiles.POST("", profileHandler.CreateProfile)                    // Create profile
			profiles.GET("", profileHandler.ListProfiles)                      // List profiles
			profiles.GET("/:id", profileHandler.GetProfile)                    // Get profile by ID
			profiles.PUT("/:id", profileHandler.UpdateProfile)                 // Update profile
			profiles.DELETE("/:id", profileHandler.DeleteProfile)              // Delete profile
			profiles.GET("/user/:user_id", profileHandler.GetProfileByUserID)  // Get profile by user ID
		}

		// Matching routes
		matchings := v1.Group("/matchings")
		{
			matchings.POST("", matchingHandler.CreateMatching)      // Create matching
			matchings.GET("", matchingHandler.ListMatchings)        // List matchings
			matchings.GET("/:id", matchingHandler.GetMatching)      // Get matching by ID
			matchings.PUT("/:id", matchingHandler.UpdateMatching)   // Update matching
			matchings.DELETE("/:id", matchingHandler.DeleteMatching) // Delete matching
		}

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
			bookmarks.POST("", bookmarkHandler.CreateBookmark)     // Create bookmark
			bookmarks.GET("", bookmarkHandler.ListBookmarks)       // List bookmarks
			bookmarks.PUT("/:id", bookmarkHandler.UpdateBookmark)  // Update bookmark
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
			hackathons.POST("/:id/participants", hackathonHandler.CreateParticipant)      // Register for hackathon
			hackathons.GET("/:id/participants", hackathonHandler.ListParticipants)        // List participants
			hackathons.PUT("/:id/participants/:participant_id", hackathonHandler.UpdateParticipant) // Update participant
			hackathons.DELETE("/:id/participants/:participant_id", hackathonHandler.DeleteParticipant) // Remove participant
		}
	}
}
