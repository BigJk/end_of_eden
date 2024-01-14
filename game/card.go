package game

import (
	"encoding/gob"
	"github.com/BigJk/end_of_eden/internal/lua/luhelp"
	"github.com/samber/lo"
	"strings"
)

func init() {
	gob.Register(CardInstance{})
}

// Card represents a playable card definition.
type Card struct {
	ID          string
	Name        string
	Description string
	Tags        []string
	State       luhelp.OwnedCallback
	Color       string
	PointCost   int
	MaxLevel    int
	DoesExhaust bool
	DoesConsume bool
	NeedTarget  bool
	Price       int
	Callbacks   map[string]luhelp.OwnedCallback
	Test        luhelp.OwnedCallback
	BaseGame    bool
}

func (c Card) PublicTags() []string {
	return lo.Filter(c.Tags, func(s string, i int) bool {
		return !strings.HasPrefix(s, "_")
	})
}

// CardInstance represents an instance of a card owned by some actor.
type CardInstance struct {
	TypeID string
	GUID   string
	Level  int
	Owner  string
}

func (c CardInstance) IsNone() bool {
	return len(c.GUID) == 0
}
