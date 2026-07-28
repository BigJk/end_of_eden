package gen

import (
	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/end_of_eden/system/localization"
	"log"
	"math/rand"
	"strings"
)

var data = map[string][]string{}

// InitGen loads all data from the assets/gen folder. This function should be called on startup.
// The data is stored in a map with the type as key and a slice of strings, which are the lines
// of the file, as value.
//
// Per-locale overrides: if a sub-folder named after the current locale (e.g. assets/gen/ko/)
// exists, its .txt files take precedence over the default English files with the same name.
// This lets contributors ship locale-specific loading hints and merchant dialogue without
// touching call sites.
func InitGen() {
	files, err := fs.ReadDir("./assets/gen")
	if err != nil {
		panic(err)
	}

	// First, load default (English) entries.
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			bytes, err := fs.ReadFile("./assets/gen/" + file.Name())
			if err != nil {
				log.Println("Error reading file:", err.Error())
			}
			data[strings.Split(file.Name(), ".")[0]] = strings.Split(string(bytes), "\n")
		}
	}

	// Then, overlay locale-specific entries if a matching sub-folder exists.
	locale := localization.GetCurrent()
	localeDir := "./assets/gen/" + locale
	if localeFiles, localeErr := fs.ReadDir(localeDir); localeErr == nil {
		// Locale folder exists; overlay its .txt files on top of the defaults.
		for _, file := range localeFiles {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
				bytes, err := fs.ReadFile(localeDir + "/" + file.Name())
				if err != nil {
					log.Println("Error reading file:", err.Error())
					continue
				}
				data[strings.Split(file.Name(), ".")[0]] = strings.Split(string(bytes), "\n")
			}
		}
	}
}

// Get returns all data for the given type.
func Get(t string) []string {
	return data[t]
}

// GetRandom returns a random entry for the given type.
func GetRandom(t string) string {
	selected := data[t]
	if len(selected) == 0 {
		return ""
	}
	return selected[rand.Intn(len(selected))]
}