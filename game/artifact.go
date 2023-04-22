package game

import (
	"encoding/gob"
	"github.com/BigJk/project_gonzo/luhelp"
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
