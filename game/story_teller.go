package game

import "github.com/BigJk/end_of_eden/luhelp"

type StoryTeller struct {
	ID     string
	Active luhelp.OwnedCallback
	Decide luhelp.OwnedCallback
}
