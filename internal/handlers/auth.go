package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"seanime/internal/database/models"
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
		Token     string `json:"token"`
		BrowserId string `json:"browserId"` // Added to identify browser-specific sessions
	}

	var b body

	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}

	// Get browserId from request body or create a cookie if not present
	browserId := b.BrowserId
	if browserId == "" {
		// Check if we have a browser cookie
		cookie, err := c.Cookie("Seanime-Browser-Id")
		if err == nil && cookie.Value != "" {
			browserId = cookie.Value
		} else {
			// Generate a new browser ID using crypto/rand
			u, err := randomHex(16)
			if err == nil {
				browserId = u
			} else {
				// Fallback to timestamp-based ID
				browserId = "browser-" + time.Now().Format("20060102150405")
			}
			
			// Set browser ID cookie - expires in 1 year
			cookie := new(http.Cookie)
			cookie.Name = "Seanime-Browser-Id"
			cookie.Value = browserId
			cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
			cookie.Path = "/"
			cookie.HttpOnly = true
			c.SetCookie(cookie)
		}
	}

	// Set a new AniList client by passing to JWT token
	h.App.UpdateAnilistClientToken(b.Token)

	// Get viewer data from AniList
	getViewer, err := h.App.AnilistClient.GetViewer(context.Background())
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

	// Find any existing account with this browser ID
	var accountId uint = 0
	existingAccount, err := h.App.Database.GetAccount(browserId)
	if err == nil && existingAccount != nil {
		accountId = existingAccount.ID
	} else {
		// Create a new ID for this account
		// Query for the max ID in the database
		h.App.Logger.Debug().Msg("Creating new account ID")
		accountId = 1 // Start with ID 1 if no accounts exist
	}

	// Save account data in database with browser ID
	_, err = h.App.Database.UpsertAccount(&models.Account{
		BaseModel: models.BaseModel{
			ID:        accountId,
			UpdatedAt: time.Now(),
		},
		Username:  getViewer.Viewer.Name,
		Token:     b.Token,
		Viewer:    bytes,
		BrowserId: browserId,
	})

	if err != nil {
		return h.RespondWithError(c, err)
	}

	h.App.Logger.Info().Str("browserId", browserId).Str("username", getViewer.Viewer.Name).Msg("app: Authenticated to AniList")

	// Create a new status with the browser ID
	status := h.NewStatus(c)

	// Set the browser ID as a context value for this request
	c.Set("BrowserId", browserId)

	h.App.InitOrRefreshAnilistData()

	h.App.InitOrRefreshModules()

	go func() {
		var err error
		if err = retry(3, time.Second*3, func() error { return nil }); err != nil {
			h.App.Logger.Err(err).Msg("Failed to retry initialization")
		}
		h.App.InitOrRefreshTorrentstreamSettings()
		h.App.InitOrRefreshMediastreamSettings()
		h.App.InitOrRefreshDebridSettings()
	}()

	// Return new status
	return h.RespondWithData(c, status)

}

// HandleLogout
//
//	@summary logs out the user by removing JWT token from the database.
//	@desc It removes JWT token and Viewer data from the database.
//	@desc It creates a new handlers.Status and refreshes App modules.
//	@route /api/v1/auth/logout [POST]
//	@returns handlers.Status
func (h *Handler) HandleLogout(c echo.Context) error {
	// Get browser ID from cookie
	browserId := ""
	cookie, err := c.Cookie("Seanime-Browser-Id")
	if err == nil && cookie.Value != "" {
		browserId = cookie.Value
	}

	// If browser ID is provided in the context (set by middleware)
	if ctxBrowserId := c.Get("BrowserId"); ctxBrowserId != nil {
		if id, ok := ctxBrowserId.(string); ok && id != "" {
			browserId = id
		}
	}

	// If we have a browser ID, only log out that specific browser
	if browserId != "" {
		// Get the account associated with this browser ID
		account, err := h.App.Database.GetAccount(browserId)
		if err == nil && account != nil {
			// Update only this browser's session
			_, err = h.App.Database.UpsertAccount(&models.Account{
				BaseModel: models.BaseModel{
					ID:        account.ID,
					UpdatedAt: time.Now(),
				},
				Username:  "",
				Token:     "",
				Viewer:    nil,
				BrowserId: browserId,
			})

			if err != nil {
				return h.RespondWithError(c, err)
			}

			h.App.Logger.Info().Str("browserId", browserId).Msg("Logged out specific browser from AniList")
		} else {
			// Browser ID provided but no account found, log the issue
			h.App.Logger.Warn().Str("browserId", browserId).Msg("No account found for browser ID during logout")
		}
	} else {
		// Legacy fallback: update the default account (ID=1)
		_, err := h.App.Database.UpsertAccount(&models.Account{
			BaseModel: models.BaseModel{
				ID:        1,
				UpdatedAt: time.Now(),
			},
			Username: "",
			Token:    "",
			Viewer:   nil,
		})

		if err != nil {
			return h.RespondWithError(c, err)
		}

		h.App.Logger.Info().Msg("Logged out default account from AniList")
	}

	status := h.NewStatus(c)

	h.App.InitOrRefreshModules()

	h.App.InitOrRefreshAnilistData()

	return h.RespondWithData(c, status)
}
// randomHex generates a random hex string with the specified number of bytes
func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
// retry executes the provided function up to maxRetries times with the specified delay between attempts
func retry(maxRetries int, delay time.Duration, fn func() error) (err error) {
	for i := 0; i < maxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return err
}
