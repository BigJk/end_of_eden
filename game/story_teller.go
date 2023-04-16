package game

import "github.com/BigJk/project_gonzo/luhelp"

type StoryTeller struct {
	ID     string
	Active luhelp.OwnedCallback
	Decide luhelp.OwnedCallback
}
