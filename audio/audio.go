package audio

import (
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

const sampleRate = 44100
const baseVolume = -1

var mtx = sync.Mutex{}
var sounds = map[string]*beep.Buffer{}
var enabled = false
var music = &beep.Ctrl{
	Streamer: beep.Loop(-1, emptySound{}),
	Paused:   false,
}

func InitAudio() {
	// TODO: Fix audio. Currently, audio is resulting in a lot of noise.
	if runtime.GOOS == "windows" {
		log.Printf("Disable audio on windows!")
		return
	}

	wg := sync.WaitGroup{}

	_ = filepath.Walk("./assets/audio", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			var streamer beep.StreamSeekCloser
			var format beep.Format

			if !info.IsDir() {
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

	if err := speaker.Init(sampleRate, 200); err != nil {
		panic(err)
	}

	wg.Wait()

	speaker.Play(music)

	enabled = true
}

func Play(key string, volumeModifier ...float64) {
	if !enabled || runtime.GOOS == "windows" {
		return
	}

	if val, ok := sounds[key]; ok {
		volume := &effects.Volume{
			Streamer: val.Streamer(0, val.Len()),
			Base:     2,
			Volume:   baseVolume,
			Silent:   false,
		}

		if len(volumeModifier) > 0 {
			volume.Volume += volumeModifier[0]
		}

		speaker.Play(volume)
	}
}

func PlayMusic(key string) {
	if !enabled || runtime.GOOS == "windows" {
		return
	}

	if val, ok := sounds[key]; ok {
		volume := &effects.Volume{
			Streamer: beep.Loop(-1, val.Streamer(0, val.Len())),
			Base:     2,
			Volume:   baseVolume - 2,
			Silent:   false,
		}

		speaker.Lock()
		music.Streamer = volume
		speaker.Unlock()
	}
}
