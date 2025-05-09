package handlers

import (
	"context"
	"errors"
	"seanime/internal/database/models"
	sessionmiddleware "seanime/internal/middleware"
	"seanime/internal/util"
	"time"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

// HandleLogin
//
//	@summary logs in the user by saving the JWT token in the database.
//	@desc This is called when the JWT token is obtained from AniList after logging in with redirection on the client.
//	@desc It also fetches the Viewer data from AniList and saves it in the database.
//	@desc It creates a new handlers.Status and refreshes App modules.
//	@route /api/v1/auth/login [POST]
//	@returns handlers.Status
func (h *Handler) HandleLogin(c echo.Context) error {

	type body struct {
		Token string `json:"token"`
	}

	var b body

	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}
	
	// Get the session ID from the context
	sessionID := h.getSessionID(c)

	// Create a temporary AniList client with the token to validate it
	tempClient := h.App.NewAnilistClient(b.Token)

	// Get viewer data from AniList
	getViewer, err := tempClient.GetViewer(context.Background())
	if err != nil {
		h.App.Logger.Error().Msg("Could not authenticate to AniList")
		return h.RespondWithError(c, err)
	}

	if len(getViewer.Viewer.Name) == 0 {
		return h.RespondWithError(c, errors.New("could not find user"))
	}

	// Marshal viewer data
	bytes, err := json.Marshal(getViewer.Viewer)
	if err != nil {
		h.App.Logger.Err(err).Msg("scan: could not save local files")
	}

	// Save account data in database with the session ID
	_, err = h.App.Database.UpsertAccount(&models.Account{
		BaseModel: models.BaseModel{
			UpdatedAt: time.Now(),
		},
		Username:    getViewer.Viewer.Name,
		Token:       b.Token,
		Viewer:      bytes,
		SessionID:   sessionID,
		LastLoginAt: time.Now(),
	})

	if err != nil {
		return h.RespondWithError(c, err)
	}

	// Associate this token with the current session
	h.App.UpdateUserSession(sessionID, b.Token)

	h.App.Logger.Info().Str("username", getViewer.Viewer.Name).Str("sessionID", sessionID).Msg("User authenticated to AniList")

	// Refresh the session cookie
	h.refreshSession(c)

	// Create a new status
	status := h.NewStatus(c)

	// Only refresh AniList data for this specific user session
	h.App.InitOrRefreshAnilistDataForSession(sessionID)

	// Continue to refresh modules for backward compatibility
	h.App.InitOrRefreshModules()

	go func() {
		defer util.HandlePanicThen(func() {})
		h.App.InitOrRefreshTorrentstreamSettings()
		h.App.InitOrRefreshMediastreamSettings()
		h.App.InitOrRefreshDebridSettings()
	}()

	// Return new status
	return h.RespondWithData(c, status)
}

// Helper methods for session management
func (h *Handler) getSessionID(c echo.Context) string {
	return sessionmiddleware.GetSessionID(c)
}

func (h *Handler) refreshSession(c echo.Context) {
	sessionmiddleware.RefreshSession(c)
}

func (h *Handler) clearSession(c echo.Context) {
	sessionmiddleware.ClearSession(c)
}

// HandleLogout
//
//	@summary logs out the user by removing JWT token from the database.
//	@desc It removes JWT token and Viewer data from the database.
//	@desc It creates a new handlers.Status and refreshes App modules.
//	@route /api/v1/auth/logout [POST]
//	@returns handlers.Status
func (h *Handler) HandleLogout(c echo.Context) error {
	// Get the session ID
	sessionID := h.getSessionID(c)

	// Get the current account to log info about who is logging out
	currentAccount, _ := h.App.Database.GetAccountBySessionID(sessionID)
	if currentAccount != nil {
		h.App.Logger.Info().Str("username", currentAccount.Username).Str("sessionID", sessionID).Msg("User logging out")
	}

	// Clear the user's token but keep the session ID
	_, err := h.App.Database.UpsertAccount(&models.Account{
		BaseModel: models.BaseModel{
			UpdatedAt: time.Now(),
		},
		Username:  "",
		Token:     "",
		Viewer:    nil,
		SessionID: sessionID, // Keep the same session ID
	})

	if err != nil {
		return h.RespondWithError(c, err)
	}

	// Remove this session's token from the app
	h.App.RemoveUserSession(sessionID)

	// Clear the session cookie
	h.clearSession(c)

	status := h.NewStatus(c)

	// These are kept for backward compatibility
	h.App.InitOrRefreshModules()
	h.App.InitOrRefreshAnilistData()

	return h.RespondWithData(c, status)
}
