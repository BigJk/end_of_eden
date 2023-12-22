package localization

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Localization struct {
	current string
	locals  map[string]map[string]string
}

func New() *Localization {
	loc := &Localization{current: "en", locals: make(map[string]map[string]string)}
	loc.locals["en"] = make(map[string]string)
	return loc
}

// SetCurrent sets the current locale.
func (l *Localization) SetCurrent(locale string) {
	l.current = locale
}

// GetCurrent returns the current locale.
func (l *Localization) GetCurrent() string {
	return l.current
}

// Add adds a new locale with the given translations. If the locale
// already exists, the translations will be added to the existing ones.
func (l *Localization) Add(locale string, translations map[string]string) {
	if _, ok := l.locals[locale]; !ok {
		l.locals[locale] = make(map[string]string)
	}

	for k, v := range translations {
		l.locals[locale][k] = v
	}
}

// AddFolder adds all locales from the given folder. Will walk through all
// sub-folders and add all .yaml and .yml files.
func (l *Localization) AddFolder(folder string) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		return l.AddFile(path)
	})
}

// AddFile adds a new locale with the given translations from the given file.
func (l *Localization) AddFile(file string) error {
	var parsed map[string]map[string]any

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}

	var walk func(found map[string]string, m map[string]any, prefix string)
	walk = func(found map[string]string, m map[string]any, prefix string) {
		for k, v := range m {
			if m, ok := v.(map[string]any); ok {
				if prefix != "" {
					walk(found, m, prefix+"."+k)
				} else {
					walk(found, m, k)
				}
			} else {
				if prefix != "" {
					found[prefix+"."+k] = fmt.Sprint(v)
				} else {
					found[k] = fmt.Sprint(v)
				}
			}
		}
	}

	for local, v := range parsed {
		if len(local) != 2 {
			return fmt.Errorf("invalid locale: %s", local)
		}
		found := make(map[string]string)
		walk(found, v, "")
		l.Add(local, found)
	}

	return nil
}

// GetLocales returns a list of all available locales.
func (l *Localization) GetLocales() []string {
	var locales []string
	for k := range l.locals {
		locales = append(locales, k)
	}
	return locales
}

// Get returns the translation for the given key in the given locale. If the key is not found,
// it will fall back to the "en" locale and if the key is not found there, it will return the key
// or the given defaults.
func (l *Localization) Get(locale, key string, defaults ...string) string {
	if _, ok := l.locals[locale]; !ok {
		if locale != "en" {
			return l.Get("en", key, defaults...)
		}
		if len(defaults) > 0 {
			return strings.Join(defaults, " ")
		}
		return key
	}

	if _, ok := l.locals[locale][key]; !ok {
		if locale != "en" {
			return l.Get("en", key, defaults...)
		}
		if len(defaults) > 0 {
			return strings.Join(defaults, " ")
		}
		return key
	}

	return l.locals[locale][key]
}

// G is a shortcut for Get with the current locale.
func (l *Localization) G(key string, defaults ...string) string {
	return l.Get(l.current, key, defaults...)
}
