package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// GCP Database configuration via Cloud SQL Proxy
	host := "localhost"
	port := "5434"
	user := "postgres"
	password := `u6"Ml6%XD7cSg9q]`
	dbname := "tru_s3"
	sslmode := "disable"

	// Create DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	fmt.Printf("Connecting to database: %s@%s:%s/%s\n", user, host, port, dbname)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to database!")

	// Get list of tables
	var tables []string
	result := db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`).Scan(&tables)

	if result.Error != nil {
		log.Fatalf("Failed to get table list: %v", result.Error)
	}

	fmt.Printf("\nFound %d tables in the database:\n", len(tables))
	fmt.Println("=====================================")
	for i, table := range tables {
		fmt.Printf("%d. %s\n", i+1, table)
	}

	// Get detailed table information
	fmt.Println("\nTable Details:")
	fmt.Println("=====================================")
	for _, table := range tables {
		var count int64
		db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		fmt.Printf("%-25s: %d rows\n", table, count)
	}

	// Check for specific indexes
	fmt.Println("\nIndexes:")
	fmt.Println("=====================================")
	var indexes []struct {
		TableName string `json:"table_name"`
		IndexName string `json:"index_name"`
	}
	
	db.Raw(`
		SELECT 
			t.relname AS table_name,
			i.relname AS index_name
		FROM pg_index ix
		JOIN pg_class t ON t.oid = ix.indrelid
		JOIN pg_class i ON i.oid = ix.indexrelid
		WHERE t.relkind = 'r'
		AND t.relname NOT LIKE 'pg_%'
		ORDER BY t.relname, i.relname
	`).Scan(&indexes)

	currentTable := ""
	for _, idx := range indexes {
		if idx.TableName != currentTable {
			fmt.Printf("\n%s:\n", idx.TableName)
			currentTable = idx.TableName
		}
		fmt.Printf("  - %s\n", idx.IndexName)
	}

	fmt.Println("\nDatabase inspection completed!")
}