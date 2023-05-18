package settings

type empty struct{}

func (e empty) LoadSettings() error {
	return nil
}

func (e empty) SaveSettings() error {
	return nil
}

func (e empty) Get(key string) any {
	return nil
}

func (e empty) GetString(key string) string {
	return ""
}

func (e empty) GetStrings(key string) []string {
	return []string{}
}

func (e empty) GetInt(key string) int {
	return 0
}

func (e empty) GetFloat(key string) float64 {
	return 0
}

func (e empty) GetBool(key string) bool {
	return false
}

func (e empty) Set(key string, value any) {

}

func (e empty) GetKeys() []string {
	return []string{}
}
