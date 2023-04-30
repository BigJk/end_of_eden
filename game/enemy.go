package game

import (
	"github.com/BigJk/end_of_eden/lua/luhelp"
)

// Enemy represents a definition of a enemy that can be linked from a Actor.
type Enemy struct {
	ID          string
	Name        string
	Description string
	InitialHP   int
	MaxHP       int
	Look        string
	Color       string
	Callbacks   map[string]luhelp.OwnedCallback
}
