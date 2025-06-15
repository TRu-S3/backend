package config

import (
	"fmt"
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

	// Database Configuration
	DBHost                 string
	DBPort                 string
	DBName                 string
	DBUser                 string
	DBPassword             string
	DBSSLMode              string
	DBMaxOpenConns         int
	DBMaxIdleConns         int
	CloudSQLConnectionName string
	UseCloudSQLProxy       bool
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

		// Database Configuration
		DBHost:                 getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:                 getEnvWithDefault("DB_PORT", "5432"),
		DBName:                 getEnvWithDefault("DB_NAME", "tru_s3"),
		DBUser:                 getEnvWithDefault("DB_USER", "postgres"),
		DBPassword:             getEnvWithDefault("DB_PASSWORD", "postgres123"),
		DBSSLMode:              getEnvWithDefault("DB_SSL_MODE", "disable"),
		DBMaxOpenConns:         getEnvIntWithDefault("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:         getEnvIntWithDefault("DB_MAX_IDLE_CONNS", 5),
		CloudSQLConnectionName: os.Getenv("CLOUD_SQL_CONNECTION_NAME"),
		UseCloudSQLProxy:       getEnvBoolWithDefault("USE_CLOUD_SQL_PROXY", false),
	}
}

// GetDatabaseDSN returns the database DSN for connection
func (c *Config) GetDatabaseDSN() string {
	if c.UseCloudSQLProxy && c.CloudSQLConnectionName != "" {
		// Cloud SQL Proxy connection
		return fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=%s",
			c.CloudSQLConnectionName, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
	}
	// Regular PostgreSQL connection
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
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

// getEnvIntWithDefault gets environment variable as integer with default value
func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBoolWithDefault gets environment variable as boolean with default value
func getEnvBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
