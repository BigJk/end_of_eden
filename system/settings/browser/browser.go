//go:build js
// +build js

package browser

import (
	"syscall/js"
)

type Browser struct{}

func (b Browser) LoadSettings() error {
	return nil
}

func (b Browser) SaveSettings() error {
	return nil
}

func (b Browser) Get(key string) any {
	return js.Global().Get("settings").Call("get", key)
}

func (b Browser) GetString(key string) string {
	return js.Global().Get("settings").Call("getString", key).String()
}

func (b Browser) GetStrings(key string) []string {
	val := js.Global().Get("settings").Call("getStrings", key)
	if val.Type() == js.TypeObject {
		var result []string
		for i := 0; i < val.Length(); i++ {
			result = append(result, val.Index(i).String())
		}
		return result
	}
	return nil
}

func (b Browser) GetInt(key string) int {
	return js.Global().Get("settings").Call("getInt", key).Int()
}

func (b Browser) GetFloat(key string) float64 {
	return js.Global().Get("settings").Call("getFloat", key).Float()
}

func (b Browser) GetBool(key string) bool {
	val := js.Global().Get("settings").Call("getBool", key)
	if val.Type() == js.TypeBoolean {
		return val.Bool()
	}
	return false
}

func (b Browser) Set(key string, value any) {
	js.Global().Get("settings").Call("set", key, value)
}

func (b Browser) GetKeys() []string {
	val := js.Global().Get("settings").Call("getKeys")
	if val.Type() == js.TypeObject {
		var result []string
		for i := 0; i < val.Length(); i++ {
			result = append(result, val.Index(i).String())
		}
		return result
	}
	return nil
}

func (b Browser) SetDefault(key string, value any) {
	js.Global().Get("settings").Call("setDefault", key, value)
}
