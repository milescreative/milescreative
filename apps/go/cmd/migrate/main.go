package main

import (
	"fmt"
	"log"
	"os"

	"mc-mono/go-server/database/migrations"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable connect_timeout=5",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Add debug logging
	fmt.Printf("Attempting to connect to: %s:%s\n", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Add this for more detailed logging
	})
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Migrations completed successfully")
}
