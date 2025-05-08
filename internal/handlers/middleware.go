package handlers

import (
	"github.com/labstack/echo/v4"
)

// BrowserIDMiddleware extracts the browser ID from the cookie or headers and adds it to the context
func BrowserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check for browser ID in the following order:
		// 1. X-Browser-ID header (set by client API)
		// 2. Seanime-Browser-Id cookie
		// 3. Seanime-Client-Id cookie (used by websocket)
		browserId := ""
		
		// Check if header exists (set by our client-side API)
		headerBrowserId := c.Request().Header.Get("X-Browser-ID")
		if headerBrowserId != "" {
			browserId = headerBrowserId
		} else {
			// Check if Seanime-Browser-Id cookie exists
			cookie, err := c.Cookie("Seanime-Browser-Id")
			if err == nil && cookie.Value != "" {
				browserId = cookie.Value
			} else {
				// As a fallback, check for Seanime-Client-Id cookie which is set by websocket provider
				clientIdCookie, err := c.Cookie("Seanime-Client-Id")
				if err == nil && clientIdCookie.Value != "" {
					browserId = clientIdCookie.Value
				}
			}
		}

		// If we found a valid browser ID, add it to the context
		if browserId != "" {
			c.Set("BrowserId", browserId)
		}

		return next(c)
	}
}
