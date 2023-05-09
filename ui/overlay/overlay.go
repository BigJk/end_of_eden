package overlay

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
)

// Code borrowed and cut down from @mrusme and https://github.com/charmbracelet/lipgloss/pull/102

// Split a string into lines, additionally returning the size of the widest line.
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := ansi.PrintableRuneWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}

// PlaceOverlay places overlay on top of background.
func PlaceOverlay(x, y int, overlay, background string) string {
	overlayLines, overlayWidth := getLines(overlay)
	backgroundLines, backgroundWidth := getLines(background)
	backgroundHeight := len(backgroundLines)
	overlayHeight := len(overlayLines)

	if overlayWidth >= backgroundWidth && overlayHeight >= backgroundHeight {
		return overlay
	}

	x = clamp(x, 0, backgroundWidth-overlayWidth)
	y = clamp(y, 0, backgroundHeight-overlayHeight)

	var b strings.Builder
	for i, backgroundLine := range backgroundLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < y || i >= y+overlayHeight {
			b.WriteString(backgroundLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := truncate.String(backgroundLine, uint(x))
			pos = ansi.PrintableRuneWidth(left)
			b.WriteString(left)
			if pos < x {
				pos = x
			}
		}

		overlayLine := overlayLines[i-y]
		b.WriteString(overlayLine)
		pos += ansi.PrintableRuneWidth(overlayLine)

		right := cutLeft(backgroundLine, pos)
		b.WriteString(right)
	}

	return b.String()
}

// cutLeft cuts printable characters from the left.
// This function is heavily based on muesli's ansi and truncate packages.
func cutLeft(s string, cutWidth int) string {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)

	for _, c := range s {
		var w int
		if c == ansi.Marker || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if ansi.IsTerminator(c) {
				isAnsi = false
				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
					ab.Reset()
				}
			}
		} else {
			w = runewidth.RuneWidth(c)
		}

		if pos >= cutWidth {
			if b.Len() == 0 {
				if ab.Len() > 0 {
					b.Write(ab.Bytes())
				}
				if pos-cutWidth > 1 {
					b.WriteByte(' ')
					continue
				}
			}
			b.WriteRune(c)
		}
		pos += w
	}

	return b.String()
}

func clamp(v, lower, upper int) int {
	return min(max(v, lower), upper)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
