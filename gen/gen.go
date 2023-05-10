package gen

import (
	"log"
	"math/rand"
	"os"
	"strings"
)

var data = map[string][]string{}

// InitGen loads all data from the assets/gen folder. This function should be called on startup.
// The data is stored in a map with the type as key and a slice of strings, which are the lines
// of the file, as value.
func InitGen() {
	files, err := os.ReadDir("./assets/gen")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			bytes, err := os.ReadFile("./assets/gen/" + file.Name())
			if err != nil {
				log.Println("Error reading file:", err.Error())
			}
			data[strings.Split(file.Name(), ".")[0]] = strings.Split(string(bytes), "\n")
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
