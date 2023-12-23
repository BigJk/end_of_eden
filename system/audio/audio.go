//go:build !no_audio && !js
// +build !no_audio,!js

// Package audio handles all audio playback. It uses the beep library to play audio files.
package audio

import (
	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/end_of_eden/system/settings"

	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

const sampleRate = 48000
const baseVolume = -1

var mtx = sync.Mutex{}
var sounds = map[string]*beep.Buffer{}
var enabled = false
var allLoaded = false
var queuedSong = ""
var music = &beep.Ctrl{
	Streamer: beep.Loop(-1, emptySound{}),
	Paused:   false,
}

// InitAudio initializes the audio system. Loads all audio files from the assets/audio folder.
func InitAudio() {
	go func() {
		wg := &sync.WaitGroup{}

		_ = fs.Walk("./assets/audio", func(path string, isDir bool) error {
			wg.Add(1)
			go func() {
				defer wg.Done()

				var streamer beep.StreamSeekCloser
				var format beep.Format

				if !isDir {
					if strings.HasSuffix(path, ".mp3") {
						f, err := os.Open(path)
						if err != nil {
							log.Println("Audio error:", err)
							return
						}

						streamer, format, err = mp3.Decode(f)
						if err != nil {
							log.Println("Audio error:", err)
							return
						}
					} else if strings.HasSuffix(path, ".wav") {
						f, err := os.Open(path)
						if err != nil {
							log.Println("Audio error:", err)
							return
						}

						streamer, format, err = wav.Decode(f)
						if err != nil {
							log.Println("Audio error:", err)
							return
						}
					}
				}

				if streamer != nil {
					buf := beep.NewBuffer(beep.Format{
						SampleRate:  sampleRate,
						NumChannels: 2,
						Precision:   2,
					})

					if format.SampleRate == sampleRate {
						buf.Append(streamer)
					} else {
						buf.Append(beep.Resample(3, format.SampleRate, sampleRate, streamer))
					}

					mtx.Lock()
					sounds[strings.Split(filepath.Base(path), ".")[0]] = buf
					mtx.Unlock()
				}
			}()

			return nil
		})

		wg.Wait()

		mtx.Lock()
		allLoaded = true
		mtx.Unlock()
	}()

	bufferSize := 200
	if runtime.GOOS == "windows" {
		// TODO: investigate why windows is misbehaving with audio
		bufferSize = sampleRate / 20
	}

	if err := speaker.Init(sampleRate, bufferSize); err != nil {
		panic(err)
	}

	speaker.Play(music)

	enabled = true
}

// Play plays a sound effect. If the sound effect is not loaded, nothing will happen.
func Play(key string, volumeModifier ...float64) {
	if !enabled {
		return
	}

	if settings.GetFloat("volume") == 0 {
		return
	}

	if val, ok := sounds[key]; ok {
		volume := &effects.Volume{
			Streamer: val.Streamer(0, val.Len()),
			Base:     2,
			Volume:   baseVolume - (1-settings.GetFloat("volume"))*5,
			Silent:   false,
		}

		if len(volumeModifier) > 0 {
			volume.Volume += volumeModifier[0]
		}

		speaker.Play(volume)
	}
}

// PlayMusic plays a music track. If the music track is not loaded, nothing will happen.
func PlayMusic(key string) {
	if !enabled {
		return
	}

	if settings.GetFloat("volume") == 0 {
		return
	}

	// If not all audio files are loaded, yet we will remember which song was requested
	// and play them when the loading is done.
	if !allLoaded {
		if len(queuedSong) == 0 {
			go func() {
				for !allLoaded {
					time.Sleep(time.Millisecond * 100)
				}

				mtx.Lock()
				PlayMusic(queuedSong)
				mtx.Unlock()
			}()
		}

		mtx.Lock()
		queuedSong = key
		mtx.Unlock()
	}

	if val, ok := sounds[key]; ok {
		volume := &effects.Volume{
			Streamer: beep.Loop(-1, val.Streamer(0, val.Len())),
			Base:     2,
			Volume:   baseVolume - 2 - (1-settings.GetFloat("volume"))*5,
			Silent:   false,
		}

		speaker.Lock()
		music.Streamer = volume
		speaker.Unlock()
	}
}
