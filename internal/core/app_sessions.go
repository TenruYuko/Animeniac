package core

import (
	"context"
	"seanime/internal/api/anilist"
	"seanime/internal/database/models"
	"time"
)

// NewAnilistClient creates a new AniList client with the given token
func (a *App) NewAnilistClient(token string) anilist.AnilistClient {
	return anilist.NewAnilistClient(token)
}

// UpdateUserSession updates the AniList token for a specific user session
func (a *App) UpdateUserSession(sessionID, token string) {
	a.SessionManager.CreateOrUpdateSession(sessionID, token)
	
	// Always update the global client to ensure proper functionality
	// This ensures existing code works correctly with the token
	a.UpdateAnilistClientToken(token)
	
	// Also update the token in the database as a fallback
	a.saveTokenToDatabase(token)
}

// saveTokenToDatabase saves the token to the database as a fallback mechanism
func (a *App) saveTokenToDatabase(token string) {
	// Get the current account or create a new one
	account, err := a.Database.GetAccount()
	if err != nil || account == nil {
		// Create a new account with ID 1 for backward compatibility
		account = &models.Account{
			BaseModel: models.BaseModel{
				ID:        1,
				UpdatedAt: time.Now(),
			},
		}
	}
	
	// Update the token
	account.Token = token
	
	// Save to database
	_, err = a.Database.UpsertAccount(account)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to save token to database")
	}
}

// RemoveUserSession removes a user session
func (a *App) RemoveUserSession(sessionID string) {
	a.SessionManager.RemoveSession(sessionID)
}

// GetAnilistClient returns the AniList client for a session
func (a *App) GetAnilistClient(sessionID string) anilist.AnilistClient {
	// First try to get client from session
	client := a.SessionManager.GetAnilistClient(sessionID)
	if client != nil {
		return client
	}
	
	// Check if we have an account with the session ID
	account, err := a.Database.GetAccountBySessionID(sessionID)
	if err == nil && account != nil && account.Token != "" {
		// Create a new client with the account token
		client = a.NewAnilistClient(account.Token)
		// Update the session with this token
		a.SessionManager.CreateOrUpdateSession(sessionID, account.Token)
		// Also update the global client for backward compatibility
		a.UpdateAnilistClientToken(account.Token)
		return client
	}
	
	// Fall back to global client
	return a.AnilistClient
}

// InitOrRefreshAnilistDataForSession initializes or refreshes AniList data for a specific session
func (a *App) InitOrRefreshAnilistDataForSession(sessionID string) {
	if sessionID == "" {
		// For backward compatibility, use the global method
		a.InitOrRefreshAnilistData()
		return
	}
	
	session := a.SessionManager.GetSession(sessionID)
	if session == nil {
		a.Logger.Warn().Str("sessionID", sessionID).Msg("Cannot refresh AniList data: session not found")
		return
	}
	
	// Use the session's AniList client for refreshing data
	client := session.AnilistClient
	
	// Save the original client to restore it later
	originalClient := a.AnilistClient
	
	// Temporarily set the session's client as the global one
	a.AnilistClient = client
	a.AnilistPlatform.SetAnilistClient(client)
	
	// Refresh collections
	_, _ = a.RefreshAnimeCollection()
	_, _ = a.RefreshMangaCollection()
	
	// Restore the original client
	a.AnilistClient = originalClient
	a.AnilistPlatform.SetAnilistClient(originalClient)
}

// GetCurrentViewerForSession gets the current viewer data for a specific session
func (a *App) GetCurrentViewerForSession(sessionID string) (*anilist.GetViewer, error) {
	client := a.GetAnilistClient(sessionID)
	
	return client.GetViewer(context.Background())
}
