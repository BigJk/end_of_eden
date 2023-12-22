// Package settings provides a simple interface for storing and retrieving game settings.
package settings

// Settings is an interface for a game settings store.
type Settings interface {
	LoadSettings() error
	SaveSettings() error
	Get(key string) any
	GetString(key string) string
	GetStrings(key string) []string
	GetInt(key string) int
	GetFloat(key string) float64
	GetBool(key string) bool
	Set(key string, value any)
	GetKeys() []string
}

func init() {
	global = empty{}
}

var global Settings

// GetGlobal returns the global settings store.
func GetGlobal() Settings {
	return global
}

// SetSettings sets the global settings store.
func SetSettings(s Settings) {
	global = s
}

// LoadSettings loads settings from the global settings store.
func LoadSettings() error {
	return global.LoadSettings()
}

// SaveSettings saves settings to the global settings store.
func SaveSettings() error {
	return global.SaveSettings()
}

// Get returns a value from the global settings store.
func Get(key string) any {
	return global.Get(key)
}

// GetString returns a string from the global settings store.
func GetString(key string) string {
	return global.GetString(key)
}

// GetStrings returns a slice of strings from the global settings store.
func GetStrings(key string) []string {
	return global.GetStrings(key)
}

// GetInt returns an int from the global settings store.
func GetInt(key string) int {
	return global.GetInt(key)
}

// GetFloat returns a float64 from the global settings store.
func GetFloat(key string) float64 {
	return global.GetFloat(key)
}

// GetBool returns a bool from the global settings store.
func GetBool(key string) bool {
	return global.GetBool(key)
}

// Set sets a value in the global settings store.
func Set(key string, value any) {
	global.Set(key, value)
}

// GetKeys returns the keys in the global settings store.
func GetKeys() []string {
	return global.GetKeys()
}
