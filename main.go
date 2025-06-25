package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TRu-S3/backend/internal/application"
	"github.com/TRu-S3/backend/internal/config"
	"github.com/TRu-S3/backend/internal/database"
	"github.com/TRu-S3/backend/internal/infrastructure"
	"github.com/TRu-S3/backend/internal/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}
	// Display configuration info
	log.Printf("Configuration loaded:")
	log.Printf("  Port: %s", cfg.Port)
	log.Printf("  Gin Mode: %s", cfg.GinMode)
	log.Printf("  GCS Bucket: %s", cfg.GCSBucketName)
	log.Printf("  GCS Folder: %s", cfg.GCSFolder)
	log.Printf("  Database: %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	if cfg.GoogleCloudProject != "" {
		log.Printf("  GCP Project: %s", cfg.GoogleCloudProject)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := database.Migrate(database.GetDB()); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get configuration from config instead of environment variables directly
	bucketName := cfg.GCSBucketName
	folder := cfg.GCSFolder

	// Create GCS client
	gcsClient, err := infrastructure.NewGCSClient(ctx)
	if err != nil {
		log.Fatal("Failed to create GCS client:", err)
	}

	// Create repository
	fileRepo := infrastructure.NewGCSFileRepository(gcsClient, bucketName, folder)

	// Create service
	fileService := application.NewFileService(fileRepo)

	// Create handler
	fileHandler := interfaces.NewFileHandler(fileService)

	// Create contest handler
	contestHandler := interfaces.NewContestHandler(database.GetDB())

	// Create bookmark handler
	bookmarkHandler := interfaces.NewBookmarkHandler(database.GetDB())

	// Create hackathon handler
	hackathonHandler := interfaces.NewHackathonHandler(database.GetDB())

	// Set Gin mode from configuration
	gin.SetMode(cfg.GinMode)

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add a simple health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "TRu-S3 Backend is running!",
			"status":  "healthy",
		})
	})

	// Add a health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Setup API routes
	interfaces.SetupRoutes(r, fileHandler, contestHandler, bookmarkHandler, hackathonHandler)

	// Create HTTP server with port from configuration
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on :%s (mode: %s)", cfg.Port, cfg.GinMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give a timeout of 30 seconds to shutdown gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
