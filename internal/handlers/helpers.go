package handlers

import (
	"seanime/internal/api/anilist"
	"github.com/labstack/echo/v4"
)

// GetAnilistClientFromContext retrieves the appropriate AniList client for this request
// based on browser ID from context or cookie
func (h *Handler) GetAnilistClientFromContext(c echo.Context) anilist.AnilistClient {
	// Get browser ID from context
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

	// Get browser-specific AniList client
	if browserId != "" {
		return h.App.GetAnilistClientForBrowser(browserId)
	}

	// Fallback to default AniList client
	return h.App.AnilistClient
}
