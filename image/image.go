package image

import (
	"errors"
	"github.com/BigJk/imeji"
	"github.com/BigJk/imeji/charmaps"
	"github.com/muesli/termenv"
	"os"
	"path/filepath"
)

// TODO: Better decoupling in relation to session

var searchPaths []string

func init() {
	ResetSearchPaths()
}

func AddSearchPaths(paths ...string) {
	searchPaths = append(searchPaths, paths...)
}

func ResetSearchPaths() {
	searchPaths = []string{"./assets/images/"}
}

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

	for i := range searchPaths {
		res, err := imeji.FileString(filepath.Join(searchPaths[i], name), options...)
		if err == nil {
			return res, nil
		}
	}
	return "", errors.New("could not load image")
}
