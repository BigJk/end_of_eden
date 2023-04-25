package audio

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

const sampleRate = 44100
const baseVolume = -1

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

	_ = filepath.Walk("./assets/audio", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		var streamer beep.StreamSeekCloser
		var format beep.Format

		if !info.IsDir() {
			if strings.HasSuffix(path, ".mp3") {
				f, err := os.Open(path)
				if err != nil {
					return err
				}

				streamer, format, err = mp3.Decode(f)
				if err != nil {
					return err
				}
			} else if strings.HasSuffix(path, ".wav") {
				f, err := os.Open(path)
				if err != nil {
					return err
				}

				streamer, format, err = wav.Decode(f)
				if err != nil {
					return err
				}
			}
		}

		if streamer != nil {
			buf := beep.NewBuffer(beep.Format{
				SampleRate:  sampleRate,
				NumChannels: 2,
				Precision:   2,
			})
			buf.Append(beep.Resample(6, format.SampleRate, sampleRate, streamer))
			sounds[strings.Split(filepath.Base(path), ".")[0]] = buf
		}

		return nil
	})

	if err := speaker.Init(sampleRate, 200); err != nil {
		panic(err)
	}

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
