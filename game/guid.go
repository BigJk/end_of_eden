package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func NewGuid(tags ...string) string {
	return strings.Join(append(tags, fmt.Sprint(time.Now().UnixMilli()), fmt.Sprint(rand.Intn(100000))), "-")
}
