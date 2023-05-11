package termgl

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"os"
)

func LoadFace(file string, dpi float64, size float64) font.Face {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}

	return face
}
