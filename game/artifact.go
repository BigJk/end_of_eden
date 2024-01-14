package game

import (
	"encoding/gob"
	"github.com/BigJk/end_of_eden/internal/lua/luhelp"
	"github.com/samber/lo"
	"strings"
)

func init() {
	gob.Register(ArtifactInstance{})
}

type Artifact struct {
	ID          string
	Name        string
	Description string
	Tags        []string
	Order       int
	Price       int
	Callbacks   map[string]luhelp.OwnedCallback
	Test        luhelp.OwnedCallback
	BaseGame    bool
}

func (a Artifact) PublicTags() []string {
	return lo.Filter(a.Tags, func(s string, i int) bool {
		return !strings.HasPrefix(s, "_")
	})
}

type ArtifactInstance struct {
	TypeID string
	GUID   string
	Owner  string
}
