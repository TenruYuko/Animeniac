package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	// SessionCookieName is the name of the cookie that stores the session ID
	SessionCookieName = "seanime_session"
	// SessionCookieMaxAge is the maximum age of the session cookie in seconds (30 days)
	SessionCookieMaxAge = 60 * 60 * 24 * 30
)

// SessionMiddleware adds session handling capabilities to the Echo server
func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if the user already has a session
			cookie, err := c.Cookie(SessionCookieName)
			
			// If no cookie exists or it's expired, create a new session
			if err != nil || cookie.Value == "" {
				// Generate a new session ID
				sessionID := uuid.New().String()
				
				// Create a new cookie
				newCookie := &http.Cookie{
					Name:     SessionCookieName,
					Value:    sessionID,
					Path:     "/",
					MaxAge:   SessionCookieMaxAge,
					HttpOnly: true,
					Secure:   c.Request().TLS != nil, // Secure if HTTPS
					SameSite: http.SameSiteLaxMode,
				}
				
				// Set the cookie
				c.SetCookie(newCookie)
				
				// Store session ID in context
				c.Set("sessionID", sessionID)
			} else {
				// Use existing session ID
				c.Set("sessionID", cookie.Value)
			}
			
			// Continue with the request
			return next(c)
		}
	}
}

// GetSessionID retrieves the session ID from the context
func GetSessionID(c echo.Context) string {
	sessionID, ok := c.Get("sessionID").(string)
	if !ok || sessionID == "" {
		// This shouldn't happen if the middleware is properly set up
		// but just in case, generate a new session ID
		return uuid.New().String()
	}
	return sessionID
}

// RefreshSession extends the session cookie lifetime
func RefreshSession(c echo.Context) {
	sessionID := GetSessionID(c)
	
	// Create a new cookie with the same session ID but refreshed expiry
	newCookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   SessionCookieMaxAge,
		HttpOnly: true,
		Secure:   c.Request().TLS != nil, // Secure if HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	
	// Set the refreshed cookie
	c.SetCookie(newCookie)
}

// ClearSession clears the session cookie
func ClearSession(c echo.Context) {
	// Create an expired cookie to clear the session
	expiredCookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   c.Request().TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-24 * time.Hour), // Set expiry to the past
	}
	
	// Set the expired cookie
	c.SetCookie(expiredCookie)
}
