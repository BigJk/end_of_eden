package main

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/samber/lo"
	"math/rand"
)

func Shuffle[T any](rnd *rand.Rand, collection []T) []T {
	rnd.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})

	return collection
}

var Operations = map[string]func(rnd *rand.Rand, s *game.Session) string{
	"FinishPlayerTurn": func(rnd *rand.Rand, s *game.Session) string {
		s.FinishPlayerTurn()
		return "Finish player turn"
	},
	"FinishFight": func(rnd *rand.Rand, s *game.Session) string {
		s.FinishFight()
		return "Finish fight"
	},
	"CastCard": func(rnd *rand.Rand, s *game.Session) string {
		guid := Shuffle(rnd, lo.Flatten([][]string{{""}, s.GetInstances(), s.GetActors()}))[0]
		target := Shuffle(rnd, lo.Flatten([][]string{{""}, s.GetInstances(), s.GetActors()}))[0]
		s.CastCard(guid, target)
		return fmt.Sprintf("Cast card with guid '%s' on '%s'", guid, target)
	},
	"AddActorFromEnemy": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		enemyId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.Enemies)}))[0]
		s.AddActorFromEnemy(enemyId)
		return fmt.Sprintf("Added enemy '%s'", enemyId)
	},
	"SetEvent": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		eventId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.Events)}))[0]
		s.SetEvent(eventId)
		return fmt.Sprintf("Set event '%s'", eventId)
	},
	"SetGameState": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		eventId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.Events)}))[0]
		s.SetGameState(Shuffle(rnd, []game.GameState{game.GameStateGameOver, game.GameStateMerchant, game.GameStateRandom, game.GameStateEvent, game.GameStateFight, game.GameState("")})[0])
		return fmt.Sprintf("Set event '%s'", eventId)
	},
	"FinishEvent": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		eventId := Shuffle(rnd, lo.Flatten([][]string{lo.Keys(res.Events)}))[0]
		event := res.Events[eventId]
		choice := rnd.Intn(len(event.Choices) + 1)
		s.FinishEvent(choice)
		return fmt.Sprintf("Finish event '%s' with choice %d", eventId, choice)
	},
	"CleanUpFight": func(rnd *rand.Rand, s *game.Session) string {
		s.CleanUpFight()
		return "Clean up fight"
	},
	"SetupFight": func(rnd *rand.Rand, s *game.Session) string {
		s.SetupFight()
		return "Setup fight"
	},
	"SetupMerchant": func(rnd *rand.Rand, s *game.Session) string {
		s.SetupMerchant()
		return "Setup merchant"
	},
	"LeaveMerchant": func(rnd *rand.Rand, s *game.Session) string {
		s.LeaveMerchant()
		return "Leave merchant"
	},
	"GivePlayerGold": func(rnd *rand.Rand, s *game.Session) string {
		gold := rnd.Intn(100)
		s.GivePlayerGold(gold)
		return fmt.Sprintf("Give %d gold to player", gold)
	},
	"PlayerGiveActionPoints": func(rnd *rand.Rand, s *game.Session) string {
		actionPoints := rnd.Intn(5)
		s.PlayerGiveActionPoints(actionPoints)
		return fmt.Sprintf("Give %d action points to player", actionPoints)
	},
	"AddCard": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		cardId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.Cards)}))[0]
		s.GiveCard(cardId, game.PlayerActorID)
		return fmt.Sprintf("Give '%s' card to player", cardId)
	},
	"AddArtifact": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		artifactId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.Artifacts)}))[0]
		s.GiveArtifact(artifactId, game.PlayerActorID)
		return fmt.Sprintf("Give '%s' artifact to player", artifactId)
	},
	"PlayerBuyCard": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		cardId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.Cards)}))[0]
		s.PlayerBuyCard(cardId)
		return fmt.Sprintf("Buy '%s' card as player", cardId)
	},
	"PlayerBuyArtifact": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		artifactId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.Artifacts)}))[0]
		s.PlayerBuyArtifact(artifactId)
		return fmt.Sprintf("Buy '%s' artifact as player", artifactId)
	},
	"AddStatusEffect": func(rnd *rand.Rand, s *game.Session) string {
		res := s.GetResources()
		effectId := Shuffle(rnd, lo.Flatten([][]string{{""}, lo.Keys(res.StatusEffects)}))[0]
		stacks := rnd.Intn(10)
		s.GiveStatusEffect(effectId, game.PlayerActorID, stacks)
		return fmt.Sprintf("Give '%s' status effect with %d stacks to player", effectId, stacks)
	},
	"BuyUpgradeCard": func(rnd *rand.Rand, s *game.Session) string {
		cardId := Shuffle(rnd, lo.Flatten([][]string{{""}, s.GetInstances()}))[0]
		s.BuyUpgradeCard(cardId)
		return fmt.Sprintf("Buy upgrading card '%s'", cardId)
	},
	"BuyRemoveCard": func(rnd *rand.Rand, s *game.Session) string {
		cardId := Shuffle(rnd, lo.Flatten([][]string{{""}, s.GetInstances()}))[0]
		s.BuyRemoveCard(cardId)
		return fmt.Sprintf("Buy removing card '%s'", cardId)
	},
}
