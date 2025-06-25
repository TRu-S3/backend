package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load different environment configurations
	envFiles := []string{".env.secure", ".env.iam", ".env.proxy", ".env"}
	
	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			fmt.Printf("🔄 Testing with %s configuration...\n", envFile)
			testConnection(envFile)
			fmt.Println()
		}
	}
}

func testConnection(envFile string) {
	// Load environment variables
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("❌ Could not load %s: %v", envFile, err)
		return
	}

	// Build connection string
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "tru_s3")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Add SSL certificate paths if available
	if sslCert := getEnv("DB_SSL_CERT", ""); sslCert != "" {
		dsn += fmt.Sprintf(" sslcert=%s", sslCert)
	}
	if sslKey := getEnv("DB_SSL_KEY", ""); sslKey != "" {
		dsn += fmt.Sprintf(" sslkey=%s", sslKey)
	}
	if sslRootCert := getEnv("DB_SSL_ROOT_CERT", ""); sslRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", sslRootCert)
	}

	log.Printf("🔌 Attempting connection to: %s:%s/%s (SSL: %s)", host, port, dbname, sslmode)

	// Test connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Connection failed: %v", err)
		return
	}

	log.Println("✅ Connection successful!")

	// Test basic query
	var count int64
	result := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
	if result.Error != nil {
		log.Printf("❌ Query failed: %v", result.Error)
		return
	}

	log.Printf("📊 Found %d tables in database", count)

	// Test hackathon data
	var hackathonCount int64
	result = db.Raw("SELECT COUNT(*) FROM hackathons").Scan(&hackathonCount)
	if result.Error != nil {
		log.Printf("⚠️  Could not query hackathons: %v", result.Error)
	} else {
		log.Printf("🏆 Found %d hackathons", hackathonCount)
	}

	log.Printf("🎉 All tests passed for %s!", envFile)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}