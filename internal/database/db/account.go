package db

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
	"seanime/internal/database/models"
	"sync"
	"time"
)

// accountCache maps session IDs to account objects
var accountCache sync.Map

// UpsertAccount creates or updates an account in the database
func (db *Database) UpsertAccount(acc *models.Account) (*models.Account, error) {
	// Ensure the account has a session ID
	if acc.SessionID == "" {
		acc.SessionID = uuid.New().String()
	}
	
	// Update the last login time
	if !acc.LastLoginAt.IsZero() {
		acc.LastLoginAt = time.Now()
	}

	err := db.gormdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_id"}},
		UpdateAll: true,
	}).Create(acc).Error

	if err != nil {
		db.Logger.Error().Err(err).Msg("Failed to save account in the database")
		return nil, err
	}

	// Update the cache
	accountCache.Store(acc.SessionID, acc)

	return acc, nil
}

// GetAccountBySessionID retrieves an account by its session ID
func (db *Database) GetAccountBySessionID(sessionID string) (*models.Account, error) {
	// Check the cache first
	if cachedAcc, ok := accountCache.Load(sessionID); ok {
		return cachedAcc.(*models.Account), nil
	}

	// If not in cache, query the database
	var acc models.Account
	err := db.gormdb.Where("session_id = ?", sessionID).First(&acc).Error
	if err != nil {
		return nil, err
	}
	
	if acc.Username == "" || acc.Token == "" || acc.Viewer == nil {
		return nil, errors.New("account does not exist or is not authenticated")
	}

	// Update the cache
	accountCache.Store(sessionID, &acc)

	return &acc, nil
}

// GetAccount returns the legacy account (for backward compatibility)
func (db *Database) GetAccount() (*models.Account, error) {
	// Try to get the first account in the database
	var acc models.Account
	err := db.gormdb.First(&acc).Error
	if err != nil {
		return nil, err
	}
	
	if acc.Username == "" || acc.Token == "" || acc.Viewer == nil {
		return nil, errors.New("account does not exist")
	}

	return &acc, nil
}

// GetAnilistToken retrieves the AniList token from the account with the given session ID or returns an empty string
func (db *Database) GetAnilistTokenBySessionID(sessionID string) string {
	acc, err := db.GetAccountBySessionID(sessionID)
	if err != nil {
		return ""
	}
	return acc.Token
}

// GetAnilistToken retrieves the AniList token from the account or returns an empty string (legacy method)
func (db *Database) GetAnilistToken() string {
	acc, err := db.GetAccount()
	if err != nil {
		return ""
	}
	return acc.Token
}
