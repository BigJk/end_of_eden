package game

import "github.com/BigJk/project_gonzo/luhelp"

// Card represents a playable card definition.
type Card struct {
	ID          string
	Name        string
	Description string
	State       luhelp.OwnedCallback
	Color       string
	PointCost   int
	MaxLevel    int
	DoesExhaust bool
	NeedTarget  bool
	Price       int
	Callbacks   map[string]luhelp.OwnedCallback
}

// CardInstance represents a instance of a card owned by some actor.
type CardInstance struct {
	TypeID string
	GUID   string
	Level  int
	Owner  string
}
