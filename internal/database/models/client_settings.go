package models

import "time"

// ClientSettings stores browser-specific client settings
type ClientSettings struct {
	BaseModel

	// BrowserId uniquely identifies the browser session
	BrowserId string `json:"browserId" gorm:"index"`

	// UI theme preferences
	DarkMode           *bool   `json:"darkMode"`
	AccentColor        *string `json:"accentColor"`
	CustomAccentColor  *string `json:"customAccentColor"`
	UIAnimations       *bool   `json:"uiAnimations"`
	BlurEffects        *bool   `json:"blurEffects"`
	CompactView        *bool   `json:"compactView"`
	ShowAdult          *bool   `json:"showAdult"`
	
	// Anime list view preferences
	DefaultListTab     *string `json:"defaultListTab"`
	DefaultSortOrder   *string `json:"defaultSortOrder"`
	ListViewMode       *string `json:"listViewMode"`
	
	// Media player preferences
	DefaultAudioTrack  *string `json:"defaultAudioTrack"`
	DefaultSubTrack    *string `json:"defaultSubTrack"`
	AutoplayNext       *bool   `json:"autoplayNext"`
	PreferredResolution *string `json:"preferredResolution"`
	
	// Notification preferences
	NotifyUpdates      *bool   `json:"notifyUpdates"`
	NotifyNewEpisodes  *bool   `json:"notifyNewEpisodes"`
	
	// Extra data as JSON string for future extensibility
	ExtraData          string  `json:"extraData"`
}

// NewDefaultClientSettings creates client settings with default values
func NewDefaultClientSettings(browserId string) *ClientSettings {
	darkMode := true
	accentColor := "blue"
	uiAnimations := true
	blurEffects := true
	compactView := false
	showAdult := false
	defaultListTab := "current"
	defaultSortOrder := "progress"
	listViewMode := "grid"
	autoplayNext := true
	notifyUpdates := true
	notifyNewEpisodes := true
	
	return &ClientSettings{
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		BrowserId:          browserId,
		DarkMode:           &darkMode,
		AccentColor:        &accentColor,
		UIAnimations:       &uiAnimations,
		BlurEffects:        &blurEffects,
		CompactView:        &compactView,
		ShowAdult:          &showAdult,
		DefaultListTab:     &defaultListTab,
		DefaultSortOrder:   &defaultSortOrder,
		ListViewMode:       &listViewMode,
		AutoplayNext:       &autoplayNext,
		NotifyUpdates:      &notifyUpdates,
		NotifyNewEpisodes:  &notifyNewEpisodes,
	}
}
