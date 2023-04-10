package game

import (
	"fmt"
	"github.com/BigJk/project_gonzo/game/interaction"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"sort"
	"time"
)

type GameState string

const (
	GameStateFight    = GameState("FIGHT")
	GameStateMerchant = GameState("MERCHANT")
	GameStateEvent    = GameState("EVENT")
	GameStateRandom   = GameState("RANDOM")
)

const (
	PointsPerRound = 3
	DrawSize       = 3
)

// FightState represents the current state of the fight in regard to the
// deck of the player.
type FightState struct {
	Description   string
	CurrentPoints int
	Deck          []string
	Hand          []string
	Used          []string
	Exhausted     []string
}

// Session represents the state inside a game session.
type Session struct {
	luaState     *lua.LState
	resources    *ResourcesManager
	interactions interaction.Provider

	state         GameState
	actors        map[string]*Actor
	instances     map[string]any
	stagesCleared int
	currentEvent  *Event
	currentFight  FightState

	Logs []LogEntry
}

func NewSession(options ...func(s *Session)) *Session {
	session := &Session{
		interactions: &interaction.EmptyProvider{},
		state:        GameStateEvent,
		actors: map[string]*Actor{
			PlayerActorID: NewActor(PlayerActorID),
		},
		instances:     map[string]any{},
		stagesCleared: 0,
	}

	session.luaState = SessionAdapter(session)
	session.resources = NewResourcesManager(session.luaState)
	session.SetEvent("START")

	for i := range options {
		options[i](session)
	}

	return session
}

func (s *Session) WithInteractions(provider interaction.Provider) {
	s.interactions = provider
}

func (s *Session) WithAlternativeStartEvent(id string) {
	s.SetEvent(id)
}

//
// Game State Functions
//

func (s *Session) GetGameState() GameState {
	return s.state
}

func (s *Session) SetGameState(state GameState) {
	s.state = state

	switch s.state {
	case GameStateFight:
		s.SetupFight()
	}
}

func (s *Session) SetEvent(id string) {
	s.currentEvent = s.resources.Events[id]

	if s.currentEvent != nil && s.currentEvent.OnEnter != nil {
		_, _ = s.resources.Events[id].OnEnter()
	}
}

func (s *Session) GetEvent() *Event {
	if s.currentEvent == nil {
		return nil
	}
	return &*s.currentEvent
}

func (s *Session) SetupFight() {
	s.currentFight.CurrentPoints = PointsPerRound
	s.currentFight.Deck = lo.Shuffle(s.GetPlayer().Cards.ToSlice())
	s.currentFight.Hand = []string{}
	s.currentFight.Exhausted = []string{}

	s.PlayerDrawCard(DrawSize)
}
func (s *Session) GetFight() FightState {
	return s.currentFight
}

func (s *Session) FinishPlayerTurn() {
	s.currentFight.CurrentPoints = PointsPerRound

}

func (s *Session) FinishFight() bool {
	if s.GetOpponentCount(PlayerActorID) == 0 {
		s.stagesCleared += 1

		// If a event is already set we switch to it
		if s.currentEvent != nil {
			s.SetGameState(GameStateEvent)
		} else if s.stagesCleared%10 == 0 {
			s.SetEvent("CHOICE")
		} else {
			s.SetGameState(GameStateRandom)
		}
	}
	return false
}

func (s *Session) FinishEvent(choice int) {
	if s.currentEvent == nil || s.state != GameStateEvent {
		return
	}

	s.RemoveNonPlayer()

	event := s.currentEvent
	s.currentEvent = nil

	// If choice was selected and valid we try to use the next game state from the choice.
	if choice >= 0 && choice < len(event.Choices) {
		nextState, _ := event.Choices[choice].Callback()

		// If the choice dictates a new state we take that
		if nextState != nil {
			if len(nextState.(string)) > 0 {
				s.SetGameState(GameState(nextState.(string)))
			}
			_, _ = event.OnEnd(choice + 1)
			return
		}

		// Otherwise we allow OnEnd to dictate the new state
		nextState, _ = event.OnEnd(choice + 1)
		if nextState != nil && len(nextState.(string)) > 0 {
			s.SetGameState(GameState(nextState.(string)))
		}
		return
	}

	nextState, _ := event.OnEnd(nil)
	if nextState != nil && len(nextState.(string)) > 0 {
		s.SetGameState(GameState(nextState.(string)))
	}
}

func (s *Session) SetFightDescription(description string) {
	s.currentFight.Description = description
}

//
// Artifact Functions
//

func (s *Session) GetArtifactOrder(guid string) int {
	artInstance, ok := s.instances[guid]
	if !ok {
		return 0
	}
	switch artInstance := artInstance.(type) {
	case ArtifactInstance:
		if art, ok := s.resources.Artifacts[artInstance.TypeID]; ok {
			return art.Order
		}
	}
	return 0
}

func (s *Session) GetArtifact(guid string) *Artifact {
	artInstance, ok := s.instances[guid]
	if !ok {
		return nil
	}
	switch artInstance := artInstance.(type) {
	case ArtifactInstance:
		if art, ok := s.resources.Artifacts[artInstance.TypeID]; ok {
			return art
		}
	}
	return nil
}

func (s *Session) GiveArtifact(typeId string, owner string) string {
	instance := ArtifactInstance{
		TypeID: typeId,
		GUID:   NewGuid("ARTIFACT"),
		Owner:  owner,
	}
	s.instances[instance.GUID] = instance
	s.actors[owner].Artifacts.Add(instance.GUID)

	// Call OnPickUp callback for the new instance
	if onPickUp, ok := s.resources.Artifacts[typeId].Callbacks[CallbackOnPickUp]; ok {
		_, _ = onPickUp(owner, instance.GUID)
	}

	return instance.GUID
}

func (s *Session) RemoveArtifact(guid string) {
	instance := s.instances[guid].(ArtifactInstance)
	s.actors[instance.Owner].Artifacts.Remove(instance.GUID)
	delete(s.instances, guid)
}

//
// Card Functions
//

func (s *Session) GetCard(guid string) (*Card, *CardInstance) {
	cardInstance, ok := s.instances[guid]
	if !ok {
		return nil, nil
	}
	switch cardInstance := cardInstance.(type) {
	case CardInstance:
		if card, ok := s.resources.Cards[cardInstance.TypeID]; ok {
			return card, &cardInstance
		}
	}
	return nil, nil
}

func (s *Session) GiveCard(typeId string, owner string) string {
	instance := CardInstance{
		TypeID: typeId,
		GUID:   NewGuid("CARD"),
		Owner:  owner,
	}
	s.instances[instance.GUID] = instance
	s.actors[owner].Cards.Add(instance.GUID)
	return instance.GUID
}

func (s *Session) RemoveCard(guid string) {
	instance := s.instances[guid].(CardInstance)
	s.actors[instance.Owner].Cards.Remove(instance.GUID)
	delete(s.instances, guid)
}

func (s *Session) CastCard(guid string, target string) {
	if card, instance := s.GetCard(guid); card != nil {
		if onCast, ok := card.Callbacks[CallbackOnCast]; ok {
			// TODO: handle error
			_, _ = onCast(card.ID, guid, instance.Owner, target)
		}
	}
}

func (s *Session) GetCards(owner string) []string {
	return s.actors[owner].Cards.ToSlice()
}

func (s *Session) PlayerCastHand(i int, target string) {
	cardId := s.currentFight.Hand[i]

	// Only cast a card if points are available
	if card, _ := s.GetCard(cardId); card != nil {
		if s.currentFight.CurrentPoints < card.PointCost {
			return
		}
	}

	s.currentFight.Hand = lo.Reject(s.currentFight.Hand, func(item string, index int) bool {
		return index == i
	})
	s.currentFight.Used = append(s.currentFight.Used, cardId)

	s.CastCard(cardId, target)

	if len(s.currentFight.Hand) == 0 {
		s.PlayerDrawCard(DrawSize)
	}
}

func (s *Session) PlayerDrawCard(amount int) {
	for i := 0; i < amount; i++ {
		// Shuffle used back in
		if len(s.currentFight.Deck) == 0 && len(s.currentFight.Used) > 0 {
			s.currentFight.Deck = lo.Shuffle(s.currentFight.Used)
			s.currentFight.Used = []string{}
		}

		// If nothing left don't draw
		if len(s.currentFight.Deck) == 0 {
			break
		}

		s.currentFight.Hand = append(s.currentFight.Hand, s.currentFight.Deck[0])
		s.currentFight.Deck = lo.Drop(s.currentFight.Deck, 1)
	}
}

//
// Damage & Heal Function
//

func (s *Session) DealDamage(source string, target string, damage int) int {
	if val, ok := s.actors[target]; ok {
		// TODO: check status effects etc.

		artifacts := s.GetActor(source).Artifacts.ToSlice()
		sort.SliceStable(artifacts, func(i, j int) bool {
			return s.GetArtifactOrder(artifacts[i]) < s.GetArtifactOrder(artifacts[j])
		})

		for _, id := range artifacts {
			// Fetch the instance of the artifact
			artInstance, ok := s.instances[id]
			if !ok {
				continue
			}

			// Check if it's really a artifact instance
			switch artInstance := artInstance.(type) {
			case ArtifactInstance:
				// Fetch the backing definition of the type
				art, ok := s.resources.Artifacts[artInstance.TypeID]
				if !ok {
					continue
				}

				// Call damage calc callback if present
				if onDmg, ok := art.Callbacks[CallbackOnDamageCalc]; ok {
					res, err := onDmg(art.ID, artInstance.GUID, source, target, damage)
					if err != nil {
						// TODO: error handling
						continue
					}

					// Update damage
					if newDamage, ok := res.(float64); ok {
						damage = int(newDamage)
					}
				}
			}

			continue
		}

		if source == PlayerActorID {
			s.Log(LogTypeSuccess, fmt.Sprintf("You hit the enemy for %d damage", damage))
		} else {
			s.Log(LogTypeDanger, fmt.Sprintf("You took %d damage", damage))
		}

		// Negative damage aka heal is not allowed!
		if damage < 0 {
			damage = 0
		}

		val.HP = lo.Clamp(val.HP-damage, 0, val.MaxHP)

		// Remove dead non-player actor
		if target != PlayerActorID && val.HP == 0 {
			s.Log(LogTypeSuccess, fmt.Sprintf("%s died and dropped %d gold!", val.Name, val.Gold))
			s.GetPlayer().Gold += val.Gold
			s.RemoveActor(target)
			s.FinishFight()
		}

		return damage
	}
	return 0
}

func (s *Session) Heal(source string, target string, heal int) int {
	if val, ok := s.actors[target]; ok {
		// TODO: check status effects etc.

		artifacts := s.GetActor(source).Artifacts.ToSlice()
		sort.SliceStable(artifacts, func(i, j int) bool {
			return s.GetArtifactOrder(artifacts[i]) < s.GetArtifactOrder(artifacts[j])
		})

		for _, id := range artifacts {
			// Fetch the instance of the artifact
			artInstance, ok := s.instances[id]
			if !ok {
				continue
			}

			// Check if it's really a artifact instance
			switch artInstance := artInstance.(type) {
			case ArtifactInstance:
				// Fetch the backing definition of the type
				art, ok := s.resources.Artifacts[artInstance.TypeID]
				if !ok {
					continue
				}

				// Call damage calc callback if present
				if onHeal, ok := art.Callbacks[CallbackOnHealCalc]; ok {
					res, err := onHeal(art.ID, artInstance.GUID, source, target, heal)
					if err != nil {
						// TODO: error handling
						continue
					}

					// Update damage
					if newHeal, ok := res.(float64); ok {
						heal = int(newHeal)
					}
				}
			}

			continue
		}

		// Negative heal aka damage is not allowed!
		if heal < 0 {
			heal = 0
		}

		val.HP = lo.Clamp(val.HP+heal, 0, val.MaxHP)
		return heal
	}
	return 0
}

//
// Actor Functions
//

func (s *Session) GetPlayer() *Actor {
	return s.actors[PlayerActorID]
}

func (s *Session) GetActor(id string) *Actor {
	return s.actors[id]
}

func (s *Session) AddActor(actor *Actor) {
	s.actors[actor.ID] = actor
}

func (s *Session) AddActorFromEnemy(id string) string {
	if base, ok := s.resources.Enemies[id]; ok {
		actor := NewActor(NewGuid(id))

		actor.TypeID = id
		actor.Name = base.Name
		actor.Description = base.Description
		actor.HP = base.InitialHP
		actor.MaxHP = base.MaxHP

		s.AddActor(actor)

		if onInit, ok := base.Callbacks[CallbackOnInit]; ok {
			// TODO: error handling
			_, _ = onInit(id, actor.ID)
		}

		return actor.ID
	}

	return ""
}

func (s *Session) RemoveActor(id string) {
	var deleteInstances []string

	for _, val := range s.instances {
		switch val := val.(type) {
		case CardInstance:
			if val.Owner == id {
				deleteInstances = append(deleteInstances, id)
			}
		case ArtifactInstance:
			if val.Owner == id {
				deleteInstances = append(deleteInstances, id)
			}
		}
	}

	// Clear actor owned items
	for _, k := range deleteInstances {
		delete(s.instances, k)
	}

	delete(s.actors, id)
}

func (s *Session) RemoveNonPlayer() {
	var deleteActors []string
	for _, a := range s.actors {
		if a.ID != PlayerActorID {
			deleteActors = append(deleteActors, a.ID)
		}
	}

	for _, k := range deleteActors {
		delete(s.actors, k)
	}
}

func (s *Session) GetOpponentCount(viewpoint string) int {
	switch viewpoint {
	// From the viewpoint of the player we can have multiple enemies.
	case PlayerActorID:
		return len(lo.Filter(lo.Keys(s.actors), func(item string, index int) bool {
			return item != PlayerActorID
		}))
	// From the viewpoint of an enemy we only have the player as enemy.
	default:
		return 1
	}
}

func (s *Session) GetOpponentByIndex(viewpoint string, i int) *Actor {
	switch viewpoint {
	// From the viewpoint of the player we can have multiple enemies.
	case PlayerActorID:
		if len(s.actors) <= 1 {
			return nil
		}

		ids := lo.Keys(s.actors)
		sort.Strings(ids)
		if i < 0 || i >= len(ids) {
			return nil
		}

		return s.actors[ids[i]]
	// From the viewpoint of an enemy we only have the player as enemy.
	default:
		return s.actors[PlayerActorID]
	}
}

//
// Misc Functions
//

func (s *Session) Notification(msg string) {
	<-s.interactions.Notify(msg)
}

func (s *Session) Log(t LogType, msg string) {
	s.Logs = append(s.Logs, LogEntry{
		Time:    time.Now(),
		Type:    t,
		Message: msg,
	})
}

func (s *Session) DebugLog(msg string) {
	//TODO implement me
	panic("implement me")
}
