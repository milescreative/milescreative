package migrations

import (
	"fmt"

	"mc-mono/go-server/database/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Create auth schema
	db.Exec("CREATE SCHEMA IF NOT EXISTS auth")

	// Enable UUID extension if not exists
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Set search path
	db.Exec("SET search_path TO auth")

	// Run migrations
	if err := db.AutoMigrate(&models.User{}, &models.Session{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Add indexes
	if err := addIndexes(db); err != nil {
		return fmt.Errorf("failed to add indexes: %w", err)
	}

	return nil
}

func addIndexes(db *gorm.DB) error {
	// Add any custom indexes
	return db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users (email);
		CREATE INDEX IF NOT EXISTS idx_users_provider ON auth.users (provider);
		CREATE INDEX IF NOT EXISTS idx_users_uid ON auth.users (uid);

		CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON auth.sessions (user_uid);
		CREATE INDEX IF NOT EXISTS idx_sessions_token ON auth.sessions (token_hash);
		CREATE INDEX IF NOT EXISTS idx_sessions_uid ON auth.sessions (uid);
	`).Error
}
