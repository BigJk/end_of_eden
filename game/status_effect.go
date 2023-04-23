package game

import (
	"encoding/gob"
	"github.com/BigJk/end_of_eden/luhelp"
)

func init() {
	gob.Register(StatusEffectInstance{})
}

type DecayBehaviour string

const (
	DecayAll  = DecayBehaviour("DecayAll")
	DecayOne  = DecayBehaviour("DecayOne")
	DecayNone = DecayBehaviour("DecayNone")
)

type StatusEffect struct {
	ID          string
	Name        string
	Description string
	State       luhelp.OwnedCallback
	Look        string
	Foreground  string
	Order       int
	CanStack    bool
	Decay       DecayBehaviour
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
