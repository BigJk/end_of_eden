package settings

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

// TODO: Make this more generic so that we don't need extra settings for the _gl version.

// Settings represents the loaded game settings.
type Settings struct {
	Volume float64  `toml:"volume"`
	Mods   []string `toml:"mods"`
}

// LoadedSettings represents the loaded game settings.
var LoadedSettings Settings

// LoadSettings loads the game settings from settings.json or creates a new one with the default values.
func LoadSettings() error {
	data, err := os.ReadFile("./settings.toml")
	if err != nil {
		LoadedSettings.Volume = 1
		LoadedSettings.Mods = []string{}
		if err := SaveSettings(); err != nil {
			return err
		}
		return nil
	}

	return toml.Unmarshal(data, &LoadedSettings)
}

// SaveSettings saves the settings to settings.json.
func SaveSettings() error {
	data, err := toml.Marshal(LoadedSettings)
	if err != nil {
		return err
	}

	return os.WriteFile("./settings.toml", data, 0666)
}
