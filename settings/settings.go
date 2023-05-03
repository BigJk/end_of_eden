package settings

import (
	"encoding/json"
	"os"
)

// Settings represents the loaded game settings.
type Settings struct {
	Volume float64  `json:"volume"`
	Mods   []string `json:"mods"`
}

// LoadedSettings represents the loaded game settings.
var LoadedSettings Settings

// LoadSettings loads the game settings from settings.json or creates a new one with the default values.
func LoadSettings() error {
	data, err := os.ReadFile("./settings.json")
	if err != nil {
		LoadedSettings.Volume = 1
		LoadedSettings.Mods = []string{}
		if err := SaveSettings(); err != nil {
			return err
		}
		return nil
	}

	return json.Unmarshal(data, &LoadedSettings)
}

// SaveSettings saves the settings to settings.json.
func SaveSettings() error {
	data, err := json.MarshalIndent(LoadedSettings, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile("./settings.json", data, 0666)
}
