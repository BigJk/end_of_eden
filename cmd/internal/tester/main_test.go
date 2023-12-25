package main

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/samber/lo"
	"os"
	"strings"
	"testing"
)

// TestGame tests all artifacts, cards and status effects based on their test function.
// This is similar to the CLI tester, but this uses the native go testing framework and
// is therefore easier to integrate into CI pipelines.
func TestGame(t *testing.T) {
	mods := lo.Map(strings.Split(os.Getenv("EOE_TESTER_MODS"), ","), func(item string, index int) string {
		return strings.TrimSpace(item)
	})
	mods = lo.Filter(mods, func(item string, index int) bool {
		return len(item) > 0
	})

	session := game.NewSession(game.WithMods(mods))
	resources := session.GetResources()

	for _, artifact := range resources.Artifacts {
		if artifact.Test != nil {
			setupClean(session)
			session.GiveArtifact(artifact.ID, game.PlayerActorID)
			setupFight(session)

			t.Run(fmt.Sprintf("Artifact:%s", artifact.ID), func(t *testing.T) {
				res, err := artifact.Test.Call()
				if err != nil {
					t.Errorf("Error while testing artifact: %s", err.Error())
				} else {
					switch res := res.(type) {
					case string:
						t.Errorf("Error while testing artifact: %s", res)
					}
				}
			})
		}
	}

	for _, card := range resources.Cards {
		if card.Test != nil {
			setupClean(session)
			session.GiveCard(card.ID, game.PlayerActorID)
			setupFight(session)

			t.Run(fmt.Sprintf("Card:%s", card.ID), func(t *testing.T) {
				res, err := card.Test.Call()
				if err != nil {
					t.Errorf("Error while testing card: %s", err.Error())
				} else {
					switch res := res.(type) {
					case string:
						t.Errorf("Error while testing card: %s", res)
					}
				}
			})
		}
	}

	for _, statusEffect := range resources.StatusEffects {
		if statusEffect.Test != nil {
			setupClean(session)
			setupFight(session)
			session.GiveStatusEffect(statusEffect.ID, game.PlayerActorID, 1)

			t.Run(fmt.Sprintf("StatusEffect:%s", statusEffect.ID), func(t *testing.T) {
				res, err := statusEffect.Test.Call()
				if err != nil {
					t.Errorf("Error while testing status effect: %s", err.Error())
				} else {
					switch res := res.(type) {
					case string:
						t.Errorf("Error while testing status effect: %s", res)
					}
				}
			})
		}
	}
}
