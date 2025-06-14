package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	// Server Configuration
	Port    string
	GinMode string

	// GCP Configuration
	GCSBucketName                string
	GCSFolder                    string
	GoogleCloudProject           string
	GoogleApplicationCredentials string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// Server Configuration
		Port:    getEnvWithDefault("PORT", "8080"),
		GinMode: getEnvWithDefault("GIN_MODE", "debug"),

		// GCP Configuration
		GCSBucketName:                getEnvWithDefault("GCS_BUCKET_NAME", "202506-zenn-ai-agent-hackathon"),
		GCSFolder:                    getEnvWithDefault("GCS_FOLDER", "test"),
		GoogleCloudProject:           os.Getenv("GOOGLE_CLOUD_PROJECT"),
		GoogleApplicationCredentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.GinMode == "debug"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.GinMode == "release"
}

// GetPortInt returns port as integer
func (c *Config) GetPortInt() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
	}
	return port
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// getEnvWithDefault gets environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
