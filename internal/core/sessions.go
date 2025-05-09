package core

import (
	"seanime/internal/api/anilist"
	"sync"
)

// UserSession represents a user's session with their AniList client
type UserSession struct {
	SessionID     string
	AnilistToken  string
	AnilistClient anilist.AnilistClient
}

// SessionManager manages user sessions
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*UserSession
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*UserSession),
	}
}

// GetSession retrieves a user session by session ID
func (m *SessionManager) GetSession(sessionID string) *UserSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if session, ok := m.sessions[sessionID]; ok {
		return session
	}
	return nil
}

// CreateOrUpdateSession creates or updates a user session
func (m *SessionManager) CreateOrUpdateSession(sessionID, token string) *UserSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Create a new AniList client for this session
	client := anilist.NewAnilistClient(token)
	
	session := &UserSession{
		SessionID:     sessionID,
		AnilistToken:  token,
		AnilistClient: client,
	}
	
	m.sessions[sessionID] = session
	return session
}

// RemoveSession removes a user session
func (m *SessionManager) RemoveSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.sessions, sessionID)
}

// HasSession checks if a session exists
func (m *SessionManager) HasSession(sessionID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	_, ok := m.sessions[sessionID]
	return ok
}

// GetAnilistClient returns the AniList client for a session if it exists
func (m *SessionManager) GetAnilistClient(sessionID string) anilist.AnilistClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if session, ok := m.sessions[sessionID]; ok && session.AnilistClient != nil {
		return session.AnilistClient
	}
	return nil
}

// SyncSessions ensures all session data is properly synced
func (m *SessionManager) SyncSessions() {
	// This is a placeholder for future sync functionality
	// Currently just ensures all sessions are properly initialized
	m.mu.RLock()
	defer m.mu.RUnlock()
}
