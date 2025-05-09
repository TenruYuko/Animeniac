package core

import (
	"context"
	"seanime/internal/api/anilist"
)

// NewAnilistClient creates a new AniList client with the given token
func (a *App) NewAnilistClient(token string) anilist.AnilistClient {
	return anilist.NewAnilistClient(token)
}

// UpdateUserSession updates the AniList token for a specific user session
func (a *App) UpdateUserSession(sessionID, token string) {
	a.SessionManager.CreateOrUpdateSession(sessionID, token)
	
	// For backward compatibility, also update the global client
	// This ensures existing code continues to work while we transition to session-based clients
	a.UpdateAnilistClientToken(token)
}

// RemoveUserSession removes a user session
func (a *App) RemoveUserSession(sessionID string) {
	a.SessionManager.RemoveSession(sessionID)
}

// GetAnilistClientForSession returns the AniList client for a specific session
func (a *App) GetAnilistClientForSession(sessionID string) anilist.AnilistClient {
	session := a.SessionManager.GetSession(sessionID)
	if session != nil {
		return session.AnilistClient
	}
	
	// If no session is found, return the default client
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
	client := a.GetAnilistClientForSession(sessionID)
	
	return client.GetViewer(context.Background())
}
