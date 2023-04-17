package game

import "github.com/BigJk/project_gonzo/luhelp"

// Card represents a playable card definition.
type Card struct {
	ID          string
	Name        string
	Description string
	Color       string
	PointCost   int
	DoesExhaust bool
	NeedTarget  bool
	Callbacks   map[string]luhelp.OwnedCallback
}

// CardInstance represents a instance of a card owned by some actor.
type CardInstance struct {
	TypeID string
	GUID   string
	Owner  string
}
