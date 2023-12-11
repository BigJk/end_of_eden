package game

import (
	"github.com/BigJk/end_of_eden/lua/luhelp"
)

type StoryTeller struct {
	ID     string
	Active luhelp.OwnedCallback `json:"-"`
	Decide luhelp.OwnedCallback `json:"-"`
}
