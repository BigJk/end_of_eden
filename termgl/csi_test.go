package termgl

import (
	"fmt"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCSI(t *testing.T) {
	var testString string

	testString += fmt.Sprintf(termenv.CSI+termenv.EraseDisplaySeq, 20)
	testString += "HELLO WORLD"
	testString += fmt.Sprintf(termenv.CSI+termenv.CursorPositionSeq, 1, 2)
	testString += fmt.Sprintf(termenv.CSI+termenv.CursorPositionSeq, 1, 2)
	testString += "HELLO WORLD"
	testString += fmt.Sprintf(termenv.CSI+termenv.CursorPositionSeq, 1, 2)
	testString += fmt.Sprintf(termenv.CSI+termenv.CursorBackSeq, 5)

	var sequences []any
	for i := 0; i < len(testString); i++ {
		csi, ok := extractCSI(testString[i:])
		if ok {
			i += len(csi) - 1

			if res, ok := parseCSI(csi); ok {
				sequences = append(sequences, res)
			}
		}
	}

	assert.Equal(t, []any{
		EraseDisplaySeq{Type: 20},
		CursorPositionSeq{Row: 1, Col: 2},
		CursorPositionSeq{Row: 1, Col: 2},
		CursorPositionSeq{Row: 1, Col: 2},
		CursorBackSeq{Count: 5},
	}, sequences)
}
