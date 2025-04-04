package database

import (
	"fmt"
	"log"
	"os"

	"mc-mono/go-server/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Create the auth schema if it doesn't exist
	db.Exec("CREATE SCHEMA IF NOT EXISTS auth")

	// Set the search path to auth schema
	db.Exec("SET search_path TO auth")

	// Auto migrate models with schema prefix
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	DB = db
}
