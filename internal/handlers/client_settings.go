package handlers

import (
	"errors"
	"seanime/internal/database/models"

	"github.com/labstack/echo/v4"
)

// HandleGetClientSettings
//
//	@summary gets browser-specific client settings
//	@desc Retrieves client settings specific to the current browser session
//	@route /api/v1/client-settings [GET]
//	@returns models.ClientSettings
func (h *Handler) HandleGetClientSettings(c echo.Context) error {
	// Get browser ID from request context
	browserId := h.GetBrowserIdFromContext(c)
	if browserId == "" {
		return h.RespondWithError(c, errors.New("browser ID not found"))
	}

	// Get settings from database
	settings, err := h.App.Database.GetClientSettings(browserId)
	if err != nil {
		h.App.Logger.Error().Err(err).Str("browserId", browserId).Msg("Failed to get client settings")
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, settings)
}

// HandleUpdateClientSettings
//
//	@summary updates browser-specific client settings
//	@desc Updates client settings specific to the current browser session
//	@route /api/v1/client-settings [PUT]
//	@returns models.ClientSettings
func (h *Handler) HandleUpdateClientSettings(c echo.Context) error {
	// Get browser ID from request context
	browserId := h.GetBrowserIdFromContext(c)
	if browserId == "" {
		return h.RespondWithError(c, errors.New("browser ID not found"))
	}

	// Bind request body
	var settings models.ClientSettings
	if err := c.Bind(&settings); err != nil {
		return h.RespondWithError(c, err)
	}

	// Ensure correct browser ID
	settings.BrowserId = browserId

	// Save to database
	if err := h.App.Database.UpsertClientSettings(&settings); err != nil {
		h.App.Logger.Error().Err(err).Str("browserId", browserId).Msg("Failed to update client settings")
		return h.RespondWithError(c, err)
	}

	// Get updated settings
	updatedSettings, err := h.App.Database.GetClientSettings(browserId)
	if err != nil {
		h.App.Logger.Error().Err(err).Str("browserId", browserId).Msg("Failed to get updated client settings")
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, updatedSettings)
}
