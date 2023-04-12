package audio

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var sounds = map[string]*beep.Buffer{}

func InitAudio() {
	_ = filepath.Walk("./assets/audio", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".mp3") {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			streamer, format, err := mp3.Decode(f)
			if err != nil {
				return err
			}

			buf := beep.NewBuffer(beep.Format{
				SampleRate:  24000,
				NumChannels: 2,
				Precision:   2,
			})
			buf.Append(beep.Resample(6, format.SampleRate, 24000, streamer))
			sounds[filepath.Base(path)] = buf
		}

		return nil
	})

	if err := speaker.Init(24000, 24000/2); err != nil {
		panic(err)
	}
}

func Play(key string) {
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
