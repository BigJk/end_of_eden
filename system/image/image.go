// Package image provides a simple interface for loading images from ./assets/images
// and mods and converting them to ansi strings.
package image

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/imeji"
	"github.com/BigJk/imeji/charmaps"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

// TODO: Better decoupling in relation to session

func buildOption(options ...Option) (Options, []imeji.Option) {
	// TODO: Find a way to handle ssh nicely

	var data Options
	var imejiOptions []imeji.Option

	// Build imeji options
	if os.Getenv("EOE_IMG_SIMPLE") == "1" {
		imejiOptions = append(imejiOptions, imeji.WithPattern(charmaps.BlocksBasic))
		data.tag += "simple"
	}

	if len(os.Getenv("EOE_IMG_PATTERN")) > 0 {
		patternStr := strings.Split(os.Getenv("EOE_IMG_PATTERN"), ",")

		var pattern [][]charmaps.Pattern
		for i := range patternStr {
			if val, ok := charmaps.CharMaps[strings.TrimSpace(strings.ToLower(patternStr[i]))]; ok {
				pattern = append(pattern, val)
			}
		}

		imejiOptions = append(imejiOptions, imeji.WithPattern(pattern...))
		data.tag += os.Getenv("EOE_IMG_PATTERN")
	}

	if runtime.GOOS == "js" {
		imejiOptions = append(imejiOptions, imeji.WithTrueColor())
		data.tag += "truecolor"
	} else {
		switch termenv.DefaultOutput().Profile {
		case termenv.TrueColor:
			imejiOptions = append(imejiOptions, imeji.WithTrueColor())
			data.tag += "truecolor"
		case termenv.ANSI:
			imejiOptions = append(imejiOptions, imeji.WithANSI())
			data.tag += "ansi"
		case termenv.ANSI256:
			imejiOptions = append(imejiOptions, imeji.WithANSI256())
			data.tag += "ansi256"
		default:
			// TODO: should this be the default fallback?
			imejiOptions = append(imejiOptions, imeji.WithTrueColor())
		}
	}

	// Build image options
	for i := range options {
		imejiOptions = append(imejiOptions, options[i](&data))
	}

	return data, imejiOptions
}

// Fetch fetches an image from ./assets/images and converts it to an ansi string.
//
// env EOE_IMG_SIMPLE = 1: Forces the usage of a simpler character set. Can be used for stock cmd on windows.
func Fetch(name string, options ...Option) (string, error) {
	data, imejiOptions := buildOption(options...)
	hash, err := hash(name, data)
	if err != nil {
		return "", err
	}

	// If we are running in the browser, we reduce the color pair count to 2
	// as the algorithm is very slow running in WASM.
	if runtime.GOOS == "js" {
		imejiOptions = append(imejiOptions, imeji.WithColorPairMax(2))
	}

	if res, err := getCache(hash); err == nil {
		return string(res.([]byte)), nil
	}

	for i := range searchPaths {
		file, err := fs.ReadFile(filepath.Join(searchPaths[i], name))
		if err != nil {
			continue
		}

		image, _, err := image.Decode(bytes.NewBuffer(file))
		if err != nil {
			continue
		}

		res, err := imeji.ImageString(image, imejiOptions...)
		if err == nil {
			if err := setCache(hash, res); err != nil {
				log.Warn("could not cache image: %s", err)
			}
			return res, nil
		}
	}
	return "", errors.New("could not load image")
}

// FetchAnimation fetches an animated gif from ./assets/images and converts it to ansi strings for each frame.
//
// env EOE_IMG_SIMPLE = 1: Forces the usage of a simpler character set. Can be used for stock cmd on windows.
func FetchAnimation(name string, options ...Option) ([]string, error) {
	if !strings.HasSuffix(name, ".gif") {
		return nil, errors.New("could not load image")
	}

	data, imejiOptions := buildOption(options...)
	hash, err := hash(name, data)
	if err != nil {
		return nil, err
	}

	// If we are running in the browser, we reduce the color pair count to 2
	// as the algorithm is very slow running in WASM.
	if runtime.GOOS == "js" {
		imejiOptions = append(imejiOptions, imeji.WithColorPairMax(2))
	}

	if res, err := getCache(hash); err == nil {
		return res.([]string), nil
	}

	var frames []string
	for i := range searchPaths {
		fileBytes, err := fs.ReadFile(filepath.Join(searchPaths[i], name))
		if err != nil {
			continue
		}

		g, err := gif.DecodeAll(bytes.NewBuffer(fileBytes))
		if err != nil {
			continue
		}

		imgWidth, imgHeight := getGifDimensions(g)

		overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		draw.Draw(overpaintImage, overpaintImage.Bounds(), g.Image[0], image.ZP, draw.Src)

		for _, srcImg := range g.Image {
			draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.ZP, draw.Over)
			img, err := imeji.ImageString(overpaintImage, imejiOptions...)
			if err != nil {
				return nil, err
			}

			frames = append(frames, img)
		}

		if err := setCache(hash, frames); err != nil {
			log.Warn("could not cache image: %s", err)
		}

		return frames, nil
	}
	return nil, errors.New("could not load image")
}

func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}
