//go:build js
// +build js

// Package audio handles all audio playback. It uses the beep library to play audio files.
package audio

import (
	"path/filepath"
	"strings"
	"syscall/js"

	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/end_of_eden/system/settings"
)

// InitAudio initializes the audio system. Loads all audio files from the assets/audio folder.
func InitAudio() {}

// Play plays a sound effect. If the sound effect is not loaded, nothing will happen.
func Play(key string, volumeModifier ...float64) {
	fs.Walk("./assets/audio", func(path string, isDir bool) error {
		if !isDir && strings.HasPrefix(filepath.Base(path), key) {
			js.Global().Call("playSound", path)
		}
		return nil
	})
}

// PlayMusic plays a music track. If the music track is not loaded, nothing will happen.
func PlayMusic(key string) {
	if settings.GetFloat("volume") == 0 || !settings.GetBool("audio") {
		return
	}

	fs.Walk("./assets/audio", func(path string, isDir bool) error {
		if !isDir && strings.HasPrefix(filepath.Base(path), key) {
			js.Global().Call("loopMusic", path)
		}
		return nil
	})
}
