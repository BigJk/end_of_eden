package game

import (
	"github.com/BigJk/end_of_eden/lua/luhelp"
)

// EventChoice represents a possible choice in the Event.
type EventChoice struct {
	Description   string
	DescriptionFn luhelp.OwnedCallback
	Callback      luhelp.OwnedCallback
}

// Event represents a encounter-able event.
type Event struct {
	ID          string
	Name        string
	Description string
	Choices     []EventChoice
	OnEnter     luhelp.OwnedCallback
	OnEnd       luhelp.OwnedCallback
	Test        luhelp.OwnedCallback
	BaseGame    bool
}

func (e Event) IsNone() bool {
	return len(e.ID) == 0
}
