package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"seanime/internal/database/models"
	"time"
)

// MigrateToMultiUserSupport migrates the database to support multiple users by adding the session_id column
// and setting an initial session ID for existing accounts
func (db *Database) MigrateToMultiUserSupport() error {
	// Check if we need to update the schema
	if err := db.gormdb.AutoMigrate(&models.Account{}); err != nil {
		db.Logger.Error().Err(err).Msg("Failed to migrate Account model")
		return err
	}

	// Find existing accounts without session IDs and update them
	var accounts []models.Account
	if err := db.gormdb.Where("session_id = ? OR session_id IS NULL", "").Find(&accounts).Error; err != nil && err != gorm.ErrRecordNotFound {
		db.Logger.Error().Err(err).Msg("Failed to find accounts for migration")
		return err
	}

	// Update each account with a session ID
	for _, acc := range accounts {
		// Generate a new session ID
		sessionID := uuid.New().String()

		// Update the account
		acc.SessionID = sessionID
		acc.LastLoginAt = time.Now()

		if err := db.gormdb.Save(&acc).Error; err != nil {
			db.Logger.Error().Err(err).Msg("Failed to update account during migration")
			return err
		}

		db.Logger.Info().Str("username", acc.Username).Str("sessionID", sessionID).Msg("Migrated account to multi-user support")
	}

	return nil
}
