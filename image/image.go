package image

import (
	"github.com/BigJk/imeji"
	"github.com/BigJk/imeji/charmaps"
	"github.com/muesli/termenv"
	"os"
)

// Fetch fetches an image from ./assets/images and converts it to an ansi string.
//
// env EOE_IMG_SIMPLE = 1: Forces the usage of a simpler character set. Can be used for stock cmd on windows.
func Fetch(name string, options ...imeji.Option) (string, error) {
	// TODO: Find a way to handle ssh nicely

	if os.Getenv("EOE_IMG_SIMPLE") == "1" {
		options = append(options, imeji.WithPattern(charmaps.BlocksBasic))
	}

	switch termenv.DefaultOutput().Profile {
	case termenv.TrueColor:
		options = append(options, imeji.WithTrueColor())
	case termenv.ANSI:
		options = append(options, imeji.WithANSI())
	case termenv.ANSI256:
		options = append(options, imeji.WithANSI256())
	default:
		// TODO: should this be the default fallback?
		options = append(options, imeji.WithTrueColor())
	}

	res, err := imeji.FileString("./assets/images/"+name, options...)
	if err != nil {
		return "", err
	}
	return res, nil
}
