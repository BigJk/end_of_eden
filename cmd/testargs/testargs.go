package testargs

import (
	"flag"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui/menus/gameview"
	"github.com/BigJk/end_of_eden/ui/root"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
)

type TestArgs struct {
	Cards     *string
	Enemies   *string
	Artifacts *string
	GameState *string
	Event     *string
}

// ApplyArgs applies the test setup to the game based on the given cli arguments.
func (ta *TestArgs) ApplyArgs(baseModel tea.Model, zones *zone.Manager) tea.Model {
	if len(*ta.Cards) > 0 || len(*ta.Enemies) > 0 || len(*ta.Artifacts) > 0 || len(*ta.GameState) > 0 || len(*ta.Event) > 0 {
		session := game.NewSession(game.WithLogging(log.Default()), game.WithMods(settings.GetStrings("mods")), lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil))
		session.SetGameState(game.GameStateFight)
		session.GetPlayer().Cards.Clear()

		if len(*ta.Enemies) == 0 {
			*ta.Enemies = "DUMMY,DUMMY,DUMMY"
		}

		lo.ForEach(strings.Split(*ta.Cards, ","), func(item string, index int) {
			if len(item) == 0 {
				return
			}
			session.GiveCard(item, game.PlayerActorID)
		})

		lo.ForEach(strings.Split(*ta.Enemies, ","), func(item string, index int) {
			if len(item) == 0 {
				return
			}
			session.AddActorFromEnemy(item)
		})

		lo.ForEach(strings.Split(*ta.Artifacts, ","), func(item string, index int) {
			if len(item) == 0 {
				return
			}
			session.GiveArtifact(item, game.PlayerActorID)
		})

		session.SetupFight()

		if len(*ta.Event) > 0 {
			session.SetGameState(game.GameStateEvent)
			session.SetEvent(*ta.Event)
		}

		if len(*ta.GameState) > 0 {
			session.SetGameState(game.GameState(*ta.GameState))
		}

		return baseModel.(root.Model).PushModel(gameview.New(baseModel, zones, session))
	}
	return baseModel
}

func New() TestArgs {
	testCards := flag.String("cards", "", "test cards")
	testEnemies := flag.String("enemies", "", "test enemies")
	testArtifacts := flag.String("artifacts", "", "test artifacts")
	testGameState := flag.String("game_state", "", "test game state")
	testEvent := flag.String("event", "", "test event")

	return TestArgs{
		Cards:     testCards,
		Enemies:   testEnemies,
		Artifacts: testArtifacts,
		GameState: testGameState,
		Event:     testEvent,
	}
}
