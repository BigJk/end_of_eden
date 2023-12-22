package image

import (
	"fmt"
	"github.com/BigJk/imeji"
)

type Options struct {
	maxWidth int
	width    int
	height   int
	tag      string
}

func (o Options) String() string {
	return fmt.Sprintf("mw%d-w%d-h%d-%s", o.maxWidth, o.width, o.height, o.tag)
}

type Option func(options *Options) imeji.Option

func WithMaxWidth(maxWidth int) Option {
	return func(options *Options) imeji.Option {
		options.maxWidth = maxWidth
		return imeji.WithMaxWidth(maxWidth)
	}
}

func WithResize(width int, height int) Option {
	return func(options *Options) imeji.Option {
		options.width = width
		options.height = height
		return imeji.WithResize(width, height)
	}
}
