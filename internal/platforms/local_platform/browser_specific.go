package local_platform

import (
	"seanime/internal/api/anilist"
)

// GetAnimeCollectionWithClient retrieves anime collection using a specific AniList client
// For the local platform, this is the same as GetAnimeCollection
func (lp *LocalPlatform) GetAnimeCollectionWithClient(client anilist.AnilistClient, bypassCache bool) (*anilist.AnimeCollection, error) {
	// Simply return the local collection, ignoring client parameter since we're in offline mode
	return lp.GetAnimeCollection(bypassCache)
}

// GetMangaCollectionWithClient retrieves manga collection using a specific AniList client
// For the local platform, this is the same as GetMangaCollection
func (lp *LocalPlatform) GetMangaCollectionWithClient(client anilist.AnilistClient, bypassCache bool) (*anilist.MangaCollection, error) {
	// Simply return the local collection, ignoring client parameter since we're in offline mode
	return lp.GetMangaCollection(bypassCache)
}

// GetRawAnimeCollectionWithClient retrieves raw anime collection using a specific AniList client
// For the local platform, this is the same as GetRawAnimeCollection
func (lp *LocalPlatform) GetRawAnimeCollectionWithClient(client anilist.AnilistClient, bypassCache bool) (*anilist.AnimeCollection, error) {
	// Simply return the local collection, ignoring client parameter since we're in offline mode
	return lp.GetRawAnimeCollection(bypassCache)
}
