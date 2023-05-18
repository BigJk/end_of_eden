package settings

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

func GetGlobal() Settings {
	return global
}

func SetSettings(s Settings) {
	global = s
}

func LoadSettings() error {
	return global.LoadSettings()
}

func SaveSettings() error {
	return global.SaveSettings()
}

func Get(key string) any {
	return global.Get(key)
}

func GetString(key string) string {
	return global.GetString(key)
}

func GetStrings(key string) []string {
	return global.GetStrings(key)
}

func GetInt(key string) int {
	return global.GetInt(key)
}

func GetFloat(key string) float64 {
	return global.GetFloat(key)
}

func GetBool(key string) bool {
	return global.GetBool(key)
}

func Set(key string, value any) {
	global.Set(key, value)
}

func GetKeys() []string {
	return global.GetKeys()
}
