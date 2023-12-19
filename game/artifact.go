package game

import (
	"encoding/gob"
	"github.com/BigJk/end_of_eden/lua/luhelp"
)

func init() {
	gob.Register(ArtifactInstance{})
}

type Artifact struct {
	ID          string
	Name        string
	Description string
	Order       int
	Price       int
	Callbacks   map[string]luhelp.OwnedCallback
	Test        luhelp.OwnedCallback
	BaseGame    bool
}

type ArtifactInstance struct {
	TypeID string
	GUID   string
	Owner  string
}
