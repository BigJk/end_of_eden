package game

import "github.com/BigJk/project_gonzo/luhelp"

type StatusEffect struct {
	ID          string
	Name        string
	Description string
	State       luhelp.OwnedCallback
	Look        string
	Foreground  string
	Background  string
	Order       int
	CanStack    bool
	Rounds      int
	Callbacks   map[string]luhelp.OwnedCallback
}

type StatusEffectInstance struct {
	GUID       string
	TypeID     string
	Owner      string
	Stacks     int
	RoundsLeft int
}
