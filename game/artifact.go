package game

import (
	"encoding/gob"
	"github.com/BigJk/end_of_eden/luhelp"
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
}

type ArtifactInstance struct {
	TypeID string
	GUID   string
	Owner  string
}
