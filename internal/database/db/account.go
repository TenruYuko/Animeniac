package db

import (
	"errors"
	"gorm.io/gorm/clause"
	"seanime/internal/database/models"
)

// Map of browser IDs to account cache
var accountCacheMap = make(map[string]*models.Account)

// UpsertAccount creates or updates an account with the given browser ID
func (db *Database) UpsertAccount(acc *models.Account) (*models.Account, error) {
	// Check if we need to find an existing account by browser ID
	if acc.BrowserId != "" {
		// Look for existing account with this browser ID
		var existingAcc models.Account
		result := db.gormdb.Where("browser_id = ?", acc.BrowserId).First(&existingAcc)
		if result.Error == nil {
			// If exists, update ID to match existing record
			acc.ID = existingAcc.ID
		}
	}

	// Create or update the account
	err := db.gormdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(acc).Error

	if err != nil {
		db.Logger.Error().Err(err).Msg("Failed to save account in the database")
		return nil, err
	}

	// Update cache if browser ID is provided
	if acc.BrowserId != "" {
		accountCacheMap[acc.BrowserId] = acc
	}

	return acc, nil
}

// GetAccount retrieves an account by browser ID or the last account if no browser ID is provided
func (db *Database) GetAccount(browserId string) (*models.Account, error) {
	// If browser ID is provided and exists in cache, return it
	if browserId != "" {
		if acc, ok := accountCacheMap[browserId]; ok {
			return acc, nil
		}

		// Try to find account with this browser ID in database
		var acc models.Account
		result := db.gormdb.Where("browser_id = ?", browserId).First(&acc)
		if result.Error == nil {
			// Valid account found with this browser ID
			if acc.Username != "" && acc.Token != "" && acc.Viewer != nil {
				accountCacheMap[browserId] = &acc
				return &acc, nil
			}
		}
	}

	// If no browser ID provided or not found, return last account (for backward compatibility)
	var acc models.Account
	err := db.gormdb.Last(&acc).Error
	if err != nil {
		return nil, err
	}
	if acc.Username == "" || acc.Token == "" || acc.Viewer == nil {
		return nil, errors.New("account does not exist")
	}

	// For backward compatibility, if this account has no browser ID but is valid,
	// and we were given a browser ID, assign that browser ID to this account
	if browserId != "" && acc.BrowserId == "" {
		acc.BrowserId = browserId
		db.UpsertAccount(&acc)
	}

	return &acc, err
}

// GetAnilistToken retrieves the AniList token from the account or returns an empty string
func (db *Database) GetAnilistToken(browserId string) string {
	acc, err := db.GetAccount(browserId)
	if err != nil {
		return ""
	}
	return acc.Token
}
