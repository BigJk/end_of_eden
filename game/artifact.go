package game

import "github.com/BigJk/project_gonzo/luhelp"

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
