package termgl

import (
	"github.com/muesli/termenv"
	"strconv"
	"strings"
)

type CursorUpSeq struct {
	Count int
}

type CursorDownSeq struct {
	Count int
}

type CursorForwardSeq struct {
	Count int
}

type CursorBackSeq struct {
	Count int
}

type CursorNextLineSeq struct {
	Count int
}

type CursorPreviousLineSeq struct {
	Count int
}

type CursorHorizontalSeq struct {
	Count int
}

type CursorPositionSeq struct {
	Row int
	Col int
}

type EraseDisplaySeq struct {
	Type int
}

type EraseLineSeq struct {
	Type int
}

type ScrollUpSeq struct {
	Count int
}

type ScrollDownSeq struct {
	Count int
}

type SaveCursorPositionSeq struct{}

type RestoreCursorPositionSeq struct{}

type ChangeScrollingRegionSeq struct {
	Top    int
	Bottom int
}

type InsertLineSeq struct {
	Count int
}

type DeleteLineSeq struct {
	Count int
}

// extractCSI extracts a CSI sequence from the beginning of a string.
// It returns the sequence without any suffix, and a boolean indicating
// whether a sequence was found.
func extractCSI(s string) (string, bool) {
	if !strings.HasPrefix(s, termenv.CSI) {
		return "", false
	}

	s = s[len(termenv.CSI):]
	if len(s) == 0 {
		return "", false
	}

	for i, c := range s {
		if c >= '@' && c <= '~' {
			return termenv.CSI + s[:i+1], true
		}
	}

	return "", false
}

// parseCSI parses a CSI sequence and returns a struct representing the sequence.
func parseCSI(s string) (any, bool) {
	if !strings.HasPrefix(s, termenv.CSI) {
		return nil, false
	}

	s = s[len(termenv.CSI):]
	if len(s) == 0 {
		return nil, false
	}

	switch s[len(s)-1] {
	case 'A':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return CursorUpSeq{Count: count}, true
		}
	case 'B':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return CursorDownSeq{Count: count}, true
		}
	case 'C':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return CursorForwardSeq{Count: count}, true
		}
	case 'D':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return CursorBackSeq{Count: count}, true
		}
	case 'E':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return CursorNextLineSeq{Count: count}, true
		}
	case 'F':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return CursorPreviousLineSeq{Count: count}, true
		}
	case 'G':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return CursorHorizontalSeq{Count: count}, true
		}
	case 'H':
		if strings.Contains(s, ";") {
			parts := strings.Split(s[:len(s)-1], ";")
			if len(parts) != 2 {
				return nil, false
			}
			row, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, false
			}
			col, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, false
			}
			return CursorPositionSeq{Row: row, Col: col}, true
		}
		return nil, false
	case 'J':
		if t, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return EraseDisplaySeq{Type: t}, true
		}
	case 'K':
		if t, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return EraseLineSeq{Type: t}, true
		}
	case 'S':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return ScrollUpSeq{Count: count}, true
		}
	case 'T':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return ScrollDownSeq{Count: count}, true
		}
	case 's':
		if len(s) == 1 {
			return SaveCursorPositionSeq{}, true
		}
	case 'u':
		if len(s) == 1 {
			return RestoreCursorPositionSeq{}, true
		}
	case 'r':
		// TODO: implement
	case 'L':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return InsertLineSeq{Count: count}, true
		}
	case 'M':
		if count, err := strconv.Atoi(s[:len(s)-1]); err == nil {
			return DeleteLineSeq{Count: count}, true
		}
	}

	return nil, false
}
