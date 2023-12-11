package game

import (
	"github.com/BigJk/end_of_eden/lua/luhelp"
)

// EventChoice represents a possible choice in the Event.
type EventChoice struct {
	Description   string
	DescriptionFn luhelp.OwnedCallback `json:"-"`
	Callback      luhelp.OwnedCallback `json:"-"`
}

// Event represents a encounter-able event.
type Event struct {
	ID          string
	Name        string
	Description string
	Choices     []EventChoice
	OnEnter     luhelp.OwnedCallback `json:"-"`
	OnEnd       luhelp.OwnedCallback `json:"-"`
}

func (e Event) IsNone() bool {
	return len(e.ID) == 0
}
