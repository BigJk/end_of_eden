package settings

import (
	"encoding/json"
	"os"
)

type Settings struct {
	Volume float64  `json:"volume"`
	Mods   []string `json:"mods"`
}

var LoadedSettings Settings

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

func SaveSettings() error {
	data, err := json.MarshalIndent(LoadedSettings, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile("./settings.json", data, 0666)
}
