package db

import (
	"errors"
	"seanime/internal/database/models"
	"time"

	"gorm.io/gorm"
)

// GetClientSettings retrieves client settings for a specific browser
// If browserId is empty, returns nil
func (d *Database) GetClientSettings(browserId string) (*models.ClientSettings, error) {
	if browserId == "" {
		return nil, errors.New("browser ID cannot be empty")
	}

	var settings models.ClientSettings
	result := d.Gorm().Where("browser_id = ?", browserId).First(&settings)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create default settings if not found
			defaultSettings := models.NewDefaultClientSettings(browserId)
			if err := d.UpsertClientSettings(defaultSettings); err != nil {
				return nil, err
			}
			return defaultSettings, nil
		}
		return nil, result.Error
	}

	return &settings, nil
}

// UpsertClientSettings creates or updates client settings for a specific browser
func (d *Database) UpsertClientSettings(settings *models.ClientSettings) error {
	if settings.BrowserId == "" {
		return errors.New("browser ID cannot be empty")
	}

	// Check if settings exist for this browser
	var existingSettings models.ClientSettings
	result := d.Gorm().Where("browser_id = ?", settings.BrowserId).First(&existingSettings)
	
	settings.UpdatedAt = time.Now()
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new settings
		settings.CreatedAt = time.Now()
		if result := d.Gorm().Create(settings); result.Error != nil {
			return result.Error
		}
	} else if result.Error != nil {
		return result.Error
	} else {
		// Update existing settings
		settings.ID = existingSettings.ID
		settings.CreatedAt = existingSettings.CreatedAt
		if result := d.Gorm().Save(settings); result.Error != nil {
			return result.Error
		}
	}
	
	return nil
}

// DeleteClientSettings deletes client settings for a specific browser
func (d *Database) DeleteClientSettings(browserId string) error {
	if browserId == "" {
		return errors.New("browser ID cannot be empty")
	}
	
	result := d.Gorm().Where("browser_id = ?", browserId).Delete(&models.ClientSettings{})
	return result.Error
}
