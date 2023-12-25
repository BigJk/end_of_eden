package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/log"
	"github.com/samber/lo"
	"os"
	"strings"
)

func setupClean(session *game.Session) {
	session.UpdatePlayer(func(actor *game.Actor) bool {
		actor.HP = 100000
		actor.MaxHP = 100000
		actor.Gold = 0
		return true
	})

	player := session.GetPlayer()
	for _, artifact := range player.Artifacts.ToSlice() {
		session.RemoveArtifact(artifact)
	}
	for _, cards := range player.Cards.ToSlice() {
		session.RemoveCard(cards)
	}
}

func setupFight(session *game.Session) {
	session.SetGameState(game.GameStateFight)
	session.SetupFight()
}

func main() {
	modsFlag := flag.String("mods", "", "mods to load (e.g. 'my-mod,test-mod,another-mod')")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		fmt.Println("End Of Eden :: Tester")
		fmt.Println("The tester tests all artifacts, cards and status effects based on their test function.")
		fmt.Println()
		flag.PrintDefaults()
		return
	}

	mods := lo.Map(strings.Split(*modsFlag, ","), func(item string, index int) string {
		return strings.TrimSpace(item)
	})
	mods = lo.Filter(mods, func(item string, index int) bool {
		return len(item) > 0
	})

	allPassed := true
	session := game.NewSession(game.WithMods(mods))
	resources := session.GetResources()

	fmt.Println("--- Testing artifacts...")

	for _, artifact := range resources.Artifacts {
		if artifact.Test != nil {
			setupClean(session)
			session.GiveArtifact(artifact.ID, game.PlayerActorID)
			setupFight(session)

			res, err := artifact.Test.Call()
			if err != nil {
				log.Error("Error while testing artifact", "id", artifact.ID, "err", err)
			} else {
				switch res := res.(type) {
				case string:
					log.Error("Error while testing artifact", "id", artifact.ID, "res", res)
					allPassed = false
				default:
					log.Info("Tested artifact successfully", "id", artifact.ID)
				}
			}
		} else {
			log.Warn("Artifact has no test function", "id", artifact.ID)
		}
	}

	fmt.Println("\n--- Testing cards...")

	for _, card := range resources.Cards {
		if card.Test != nil {
			setupClean(session)
			session.GiveCard(card.ID, game.PlayerActorID)
			setupFight(session)

			res, err := card.Test.Call()
			if err != nil {
				log.Error("Error while testing card", "id", card.ID, "err", err)
			} else {
				switch res := res.(type) {
				case string:
					log.Error("Error while testing card", "id", card.ID, "res", res)
					allPassed = false
				default:
					log.Info("Tested card successfully", "id", card.ID)
				}
			}
		} else {
			log.Warn("Card has no test function", "id", card.ID)
		}
	}

	fmt.Println("\n--- Testing status effects...")
	for _, statusEffect := range resources.StatusEffects {
		if statusEffect.Test != nil {
			setupClean(session)
			session.GiveStatusEffect(statusEffect.ID, game.PlayerActorID, 1)
			setupFight(session)

			res, err := statusEffect.Test.Call()
			if err != nil {
				log.Error("Error while testing status effect", "id", statusEffect.ID, "err", err)
			} else {
				switch res := res.(type) {
				case string:
					log.Error("Error while testing status effect", "id", statusEffect.ID, "res", res)
					allPassed = false
				default:
					log.Info("Tested status effect successfully", "id", statusEffect.ID)
				}
			}
		} else {
			log.Warn("Status effect has no test function", "id", statusEffect.ID)
		}
	}

	if allPassed {
		fmt.Println("\n--- " + style.GreenText.Render("All tests passed!"))
	} else {
		fmt.Println("\n--- " + style.RedText.Render("Some tests failed!"))
		os.Exit(1)
	}
}
