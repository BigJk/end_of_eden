//go:build js
// +build js

package settings

// Settings represents the loaded game settings.
type Settings struct {
	Volume float64  `json:"volume"`
	Mods   []string `json:"mods"`
}

// LoadedSettings represents the loaded game settings.
var LoadedSettings Settings

// LoadSettings loads the game settings from settings.json or creates a new one with the default values.
func LoadSettings() error {
	return nil
}

// SaveSettings saves the settings to settings.json.
func SaveSettings() error {
	return nil
}
