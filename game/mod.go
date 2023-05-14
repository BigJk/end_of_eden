package game

import (
	"encoding/json"
	"github.com/BigJk/end_of_eden/fs"
	"path/filepath"
)

type Mod struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Version     string `json:"version"`
	URL         string `json:"url"`
}

func ModDescription(folder string) (Mod, error) {
	data, err := fs.ReadFile(filepath.Join(folder, "/meta.json"))
	if err != nil {
		return Mod{}, err
	}

	var mod Mod
	if err := json.Unmarshal(data, &mod); err != nil {
		return Mod{}, err
	}

	return mod, nil
}
