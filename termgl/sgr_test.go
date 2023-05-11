package termgl

import (
	"bytes"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSGR(t *testing.T) {
	buf := &bytes.Buffer{}
	lip := lipgloss.NewRenderer(buf, termenv.WithProfile(termenv.TrueColor))
	testString := lip.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff00ff")).Render("Hello World") + "asdasdasdasdasd" + lip.NewStyle().Italic(true).Background(lipgloss.Color("#ff00ff")).Render("Hello World")

	var sequences []any
	for i := 0; i < len(testString); i++ {
		sgr, ok := extractSGR(testString[i:])
		if ok {
			i += len(sgr) - 1

			if res, ok := parseSGR(sgr); ok {
				sequences = append(sequences, res...)
			}
		}
	}

	assert.Equal(t, []any{
		SGRBold{},
		SGRFgTrueColor{R: 255, G: 0, B: 255},
		SGRReset{},
		SGRItalic{},
		SGRBgTrueColor{R: 255, G: 0, B: 255},
		SGRReset{},
	}, sequences)
}
