package localization

// Global is the default localization instance.
var Global = New()

// SetCurrent sets the current locale in the default localization instance.
func SetCurrent(locale string) {
	Global.SetCurrent(locale)
}

// GetCurrent returns the current locale in the default localization instance.
func GetCurrent() string {
	return Global.GetCurrent()
}

// G returns the translation for the given key in the default localization instance.
func G(key string, defaults ...string) string {
	return Global.G(key, defaults...)
}
