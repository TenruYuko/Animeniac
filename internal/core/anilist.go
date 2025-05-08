package core

import (
	"errors"
	"seanime/internal/api/anilist"
	"seanime/internal/database/models"
)

// GetAccount retrieves the account for the specified browser ID or the default account
func (a *App) GetAccount(browserId string) (*models.Account, error) {
	// If a browser ID is provided, get the specific account
	if browserId != "" {
		account, err := a.Database.GetAccount(browserId)
		if err == nil && account != nil {
			if account.Username == "" {
				return nil, errors.New("no username was found for this browser")
			}
			if account.Token == "" {
				return nil, errors.New("no token was found for this browser")
			}
			return account, nil
		}
	}

	// Fallback to the cached account (if any)
	if a.account == nil {
		return nil, nil
	}

	if a.account.Username == "" {
		return nil, errors.New("no username was found")
	}

	if a.account.Token == "" {
		return nil, errors.New("no token was found")
	}

	return a.account, nil
}

// GetAccountToken retrieves the AniList token for the specified browser ID or the default account
func (a *App) GetAccountToken(browserId string) string {
	if browserId != "" {
		return a.Database.GetAnilistToken(browserId)
	}

	if a.account == nil {
		return ""
	}

	return a.account.Token
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// UpdateAnilistClientToken will update the Anilist Client Wrapper token.
// This function should be called when a user logs in
func (a *App) UpdateAnilistClientToken(token string) {
	a.AnilistClient = anilist.NewAnilistClient(token)
	a.AnilistPlatform.SetAnilistClient(a.AnilistClient) // Update Anilist Client Wrapper in Platform
}

// GetAnilistClientForBrowser retrieves or creates an AniList client for the specified browser ID
// This allows different browsers to have their own AniList sessions
func (a *App) GetAnilistClientForBrowser(browserId string) anilist.AnilistClient {
	// If no browser ID is provided, return the default client
	if browserId == "" {
		return a.AnilistClient
	}

	// Get the token for this browser
	token := a.GetAccountToken(browserId)
	
	// If no token is found, return the default client
	if token == "" {
		return a.AnilistClient
	}
	
	// Create a new client with the browser-specific token
	return anilist.NewAnilistClient(token)
}

// GetAnimeCollection returns the user's Anilist collection if it in the cache, otherwise it queries Anilist for the user's collection.
// When bypassCache is true, it will always query Anilist for the user's collection
func (a *App) GetAnimeCollection(bypassCache bool) (*anilist.AnimeCollection, error) {
	return a.AnilistPlatform.GetAnimeCollection(bypassCache)
}

// GetRawAnimeCollection is the same as GetAnimeCollection but returns the raw collection that includes custom lists
func (a *App) GetRawAnimeCollection(bypassCache bool) (*anilist.AnimeCollection, error) {
	return a.AnilistPlatform.GetRawAnimeCollection(bypassCache)
}

// RefreshAnimeCollection queries Anilist for the user's collection
func (a *App) RefreshAnimeCollection() (*anilist.AnimeCollection, error) {
	ret, err := a.AnilistPlatform.RefreshAnimeCollection()

	if err != nil {
		return nil, err
	}

	// Save the collection to PlaybackManager
	a.PlaybackManager.SetAnimeCollection(ret)

	// Save the collection to AutoDownloader
	a.AutoDownloader.SetAnimeCollection(ret)

	a.SyncManager.SetAnimeCollection(ret)

	return ret, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetMangaCollection is the same as GetAnimeCollection but for manga
func (a *App) GetMangaCollection(bypassCache bool) (*anilist.MangaCollection, error) {
	return a.AnilistPlatform.GetMangaCollection(bypassCache)
}

// GetRawMangaCollection does not exclude custom lists
func (a *App) GetRawMangaCollection(bypassCache bool) (*anilist.MangaCollection, error) {
	return a.AnilistPlatform.GetRawMangaCollection(bypassCache)
}

// RefreshMangaCollection queries Anilist for the user's manga collection
func (a *App) RefreshMangaCollection() (*anilist.MangaCollection, error) {
	mc, err := a.AnilistPlatform.RefreshMangaCollection()

	if err != nil {
		return nil, err
	}

	a.SyncManager.SetMangaCollection(mc)

	return mc, nil
}
