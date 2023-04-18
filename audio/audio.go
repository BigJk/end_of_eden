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

var sounds = map[string]*beep.Buffer{}
var enabled = false

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
				SampleRate:  24000,
				NumChannels: 2,
				Precision:   2,
			})
			buf.Append(beep.Resample(6, format.SampleRate, 24000, streamer))
			sounds[strings.Split(filepath.Base(path), ".")[0]] = buf
		}

		return nil
	})

	if err := speaker.Init(24000, 100); err != nil {
		panic(err)
	}

	enabled = true
}

func Play(key string) {
	if !enabled || runtime.GOOS == "windows" {
		return
	}

	if val, ok := sounds[key]; ok {
		volume := &effects.Volume{
			Streamer: val.Streamer(0, val.Len()),
			Base:     2,
			Volume:   -1,
			Silent:   false,
		}

		speaker.Play(volume)
	}
}
