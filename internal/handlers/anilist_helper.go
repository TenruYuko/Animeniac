package handlers

import (
	"seanime/internal/database/models"

	"github.com/labstack/echo/v4"
)

// GetAnilistTokenForSession returns the proper AniList token for the current session
// This ensures that all operations use the correct authenticated token
func (h *Handler) GetAnilistTokenForSession(c echo.Context) string {
	// Get the session ID from the context
	sessionID := h.getSessionID(c)
	
	// Try to get the token from the database by session ID first
	account, err := h.App.Database.GetAccountBySessionID(sessionID)
	if err == nil && account != nil && account.Token != "" {
		return account.Token
	}
	
	// If we don't have a token in the session, fall back to the global account
	account, _ = h.App.Database.GetAccount()
	if account != nil && account.Token != "" {
		return account.Token
	}
	
	// Return empty string if no token is found
	return ""
}

// RefreshAnilistToken ensures the token for the current session is up-to-date
// Call this after successful authentication
func (h *Handler) RefreshAnilistToken(c echo.Context, account *models.Account) {
	if account == nil || account.Token == "" {
		return
	}
	
	// Get the session ID
	sessionID := h.getSessionID(c)
	
	// Update the session with the token
	h.App.UpdateUserSession(sessionID, account.Token)
	
	// For better compatibility, also update or create the account in the database
	account.SessionID = sessionID
	_, _ = h.App.Database.UpsertAccount(account)
}
