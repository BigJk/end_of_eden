package game

import (
	"github.com/BigJk/end_of_eden/internal/lua/luhelp"
)

// Enemy represents a definition of a enemy that can be linked from a Actor.
type Enemy struct {
	ID          string
	Name        string
	Description string
	InitialHP   int
	MaxHP       int
	Gold        int
	Look        string
	Color       string
	Intend      luhelp.OwnedCallback
	Callbacks   map[string]luhelp.OwnedCallback
	Test        luhelp.OwnedCallback
	BaseGame    bool
}
