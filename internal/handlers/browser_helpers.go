package handlers

import (
	"github.com/labstack/echo/v4"
)

// GetBrowserIdFromContext extracts the browser ID from the context
// Returns empty string if no browser ID is found
func (h *Handler) GetBrowserIdFromContext(c echo.Context) string {
	browserId := ""
	if ctxBrowserId := c.Get("BrowserId"); ctxBrowserId != nil {
		if id, ok := ctxBrowserId.(string); ok && id != "" {
			browserId = id
		}
	}

	// If no browser ID in context, try cookie
	if browserId == "" {
		cookie, err := c.Cookie("Seanime-Browser-Id")
		if err == nil && cookie.Value != "" {
			browserId = cookie.Value
		}
	}

	return browserId
}
