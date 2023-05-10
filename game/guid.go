package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// NewGuid creates a new guid with the given tags. Tags are optional and only used for debugging purposes.
func NewGuid(tags ...string) string {
	return strings.Join(append(tags, fmt.Sprint(time.Now().UnixMilli()), fmt.Sprint(rand.Intn(100000))), "-")
}
