package game

import (
	"errors"
	"fmt"
	"github.com/BigJk/project_gonzo/debug"
	"github.com/BigJk/project_gonzo/gen"
	"github.com/BigJk/project_gonzo/gen/faces"
	"github.com/BigJk/project_gonzo/util"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"math/rand"
	"sort"
	"time"
)

type GameState string

const (
	GameStateFight    = GameState("FIGHT")
	GameStateMerchant = GameState("MERCHANT")
	GameStateEvent    = GameState("EVENT")
	GameStateRandom   = GameState("RANDOM")
	GameStateGameOver = GameState("GAME_OVER")
)

const (
	PointsPerRound = 3
	DrawSize       = 3
)

// FightState represents the current state of the fight in regard to the
// deck of the player.
type FightState struct {
	Round         int
	Description   string
	CurrentPoints int
	Deck          []string
	Hand          []string
	Used          []string
	Exhausted     []string
}

type MerchantState struct {
	Face      string
	Text      string
	Cards     []string
	Artifacts []string
}

// Session represents the state inside a game session.
type Session struct {
	log       *log.Logger
	luaState  *lua.LState
	resources *ResourcesManager

	state         GameState
	actors        map[string]Actor
	instances     map[string]any
	stagesCleared int
	currentEvent  *Event
	currentFight  FightState
	merchant      MerchantState
	eventHistory  []string

	stateCheckpoints []StateCheckpoint
	closer           []func() error

	Logs []LogEntry
}

func NewSession(options ...func(s *Session)) *Session {
	session := &Session{
		log:   log.New(io.Discard, "", 0),
		state: GameStateEvent,
		actors: map[string]Actor{
			PlayerActorID: NewActor(PlayerActorID),
		},
		instances:     map[string]any{},
		stagesCleared: 0,
	}

	for i := range options {
		if options[i] == nil {
			continue
		}
		options[i](session)
	}

	session.luaState = SessionAdapter(session)
	session.resources = NewResourcesManager(session.luaState, session.log)
	session.SetEvent("START")

	session.log.Println("Session started!")

	session.UpdatePlayer(func(actor *Actor) bool {
		actor.HP = 80
		actor.MaxHP = 80
		actor.Gold = 50 + rand.Intn(50)
		return true
	})

	return session
}

func WithDebugEnabled(bind string) func(s *Session) {
	return func(s *Session) {
		s.closer = append(s.closer, debug.Expose(bind, s.luaState, s.log))
	}
}

func WithLogging(logger *log.Logger) func(s *Session) {
	return func(s *Session) {
		s.log = logger
	}
}

func (s *Session) Close() {
	for i := range s.closer {
		if err := s.closer[i](); err != nil {
			s.log.Println("Close error:", err)
		}
	}
	s.luaState.Close()
}

//
// Checkpoints
//

func (s *Session) MarkState() StateCheckpointMarker {
	return StateCheckpointMarker{checkpoints: s.stateCheckpoints}
}

func (s *Session) PushState(events map[StateEvent]any) {
	savedState := *s

	// Only have the current session have the state checkpoints
	savedState.stateCheckpoints = make([]StateCheckpoint, 0)
	savedState.actors = lo.MapValues(util.CopyMap(savedState.actors), func(actor Actor, key string) Actor {
		return actor.Clone()
	})
	savedState.instances = util.CopyMap(savedState.instances)

	s.stateCheckpoints = append(s.stateCheckpoints, StateCheckpoint{
		Session: &savedState,
		Events:  events,
	})
}

func (s *Session) GetFormerState(index int) *Session {
	if index == 0 {
		return s
	}

	index = len(s.stateCheckpoints) + index
	if index >= len(s.stateCheckpoints) {
		return nil
	}

	return s.stateCheckpoints[index].Session
}

func (s *Session) FindLastState(event StateEvent) (*Session, map[StateEvent]any) {
	for i := len(s.stateCheckpoints); i >= 0; i-- {
		if s.stateCheckpoints[i].Events != nil {
			if _, ok := s.stateCheckpoints[i].Events[event]; ok {
				return s.stateCheckpoints[i].Session, s.stateCheckpoints[i].Events
			}
		}
	}
	return nil, nil
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
	case GameStateRandom:
		s.LetTellerDecide()
	case GameStateMerchant:
		s.SetupMerchant()
	}
}

func (s *Session) SetEvent(id string) {
	s.currentEvent = s.resources.Events[id]

	if s.currentEvent != nil {
		s.eventHistory = append(s.eventHistory, id)
		_, _ = s.resources.Events[id].OnEnter.Call(CreateContext("type_id", id))
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
	s.currentFight.Round = 0

	s.PlayerDrawCard(DrawSize)
	s.TriggerOnPlayerTurn()
}

func (s *Session) GetFight() FightState {
	return s.currentFight
}

func (s *Session) GetStagesCleared() int {
	return s.stagesCleared
}

func (s *Session) FinishPlayerTurn() {
	// Enemies are allowed to act.
	for k, v := range s.actors {
		if k == PlayerActorID || v.IsNone() {
			continue
		}

		if enemy, ok := s.resources.Enemies[v.TypeID]; ok {
			if _, err := enemy.Callbacks[CallbackOnTurn].Call(CreateContext("type_id", v.TypeID, "guid", k, "round", s.currentFight.Round)); err != nil {
				s.log.Printf("Error from Callback:CallbackOnTurn type=%s %s\n", v.TypeID, err.Error())
			}
		}
	}

	// Turn over so we remove all dead status effects.
	var removeStatus []string

	instanceKeys := lo.Keys(s.instances)
	for _, guid := range instanceKeys {
		switch instance := s.instances[guid].(type) {
		case StatusEffectInstance:
			se := s.resources.StatusEffects[instance.TypeID]

			// If it can decay we reduce rounds.
			if se.Decay != DecayNone {
				instance.RoundsLeft -= 1
				s.instances[guid] = instance
			}

			if _, err := s.GetStatusEffect(guid).Callbacks[CallbackOnTurn].Call(CreateContext("type_id", instance.TypeID, "guid", guid, "owner", instance.Owner, "round", s.currentFight.Round, "stacks", instance.Stacks)); err != nil {
				s.log.Printf("Error from Callback:CallbackOnTurn type=%s %s\n", instance.TypeID, err.Error())
			}

			switch se.Decay {
			// Decay stacks by one and re-set rounds if stacks left.
			case DecayOne:
				if instance.Stacks <= 0 && instance.RoundsLeft <= 0 {
					removeStatus = append(removeStatus, guid)
				} else {
					instance.Stacks -= 1
					instance.RoundsLeft = se.Rounds
					s.instances[guid] = instance
				}
			// Remove all.
			case DecayAll:
				if instance.RoundsLeft <= 0 {
					removeStatus = append(removeStatus, guid)
				}
			}
		}
	}

	for i := range removeStatus {
		s.RemoveStatusEffect(removeStatus[i])
	}

	// Advance to new Round
	s.currentFight.CurrentPoints = PointsPerRound
	s.currentFight.Round += 1
	s.currentFight.Used = append(s.currentFight.Used, s.currentFight.Hand...)
	s.currentFight.Hand = []string{}

	s.PlayerDrawCard(DrawSize)
	s.TriggerOnPlayerTurn()
}

func (s *Session) FinishFight() bool {
	if s.GetOpponentCount(PlayerActorID) == 0 {
		s.currentFight.Description = ""
		s.stagesCleared += 1

		// If an event is already set we switch to it
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

func (s *Session) GetFightRound() int {
	return s.currentFight.Round
}

func (s *Session) HadEvent(id string) bool {
	return lo.Contains(s.eventHistory, id)
}

func (s *Session) GetEventHistory() []string {
	return s.eventHistory
}

//
// Merchant
//

func (s *Session) SetupMerchant() {
	s.merchant.Artifacts = nil
	s.merchant.Cards = nil
	s.merchant.Face = faces.Global.GenRand()
	s.merchant.Text = gen.GetRandom("merchant_lines")

	for i := 0; i < 3; i++ {
		s.AddMerchantArtifact()
		s.AddMerchantCard()
	}
}

func (s *Session) LeaveMerchant() {
	s.SetGameState(GameStateRandom)
}

func (s *Session) GetMerchant() MerchantState {
	return s.merchant
}

func (s *Session) GetMerchantGoldMax() int {
	return 150 + s.stagesCleared*30
}

func (s *Session) AddMerchantArtifact() {
	possible := lo.Filter(lo.Values(s.resources.Artifacts), func(item *Artifact, index int) bool {
		return item.Price >= 0 && item.Price < s.GetMerchantGoldMax()
	})

	if len(possible) > 0 {
		s.merchant.Artifacts = append(s.merchant.Artifacts, lo.Shuffle(possible)[0].ID)
	}
}

func (s *Session) AddMerchantCard() {
	possible := lo.Filter(lo.Values(s.resources.Cards), func(item *Card, index int) bool {
		return item.Price >= 0 && item.Price < s.GetMerchantGoldMax()
	})

	if len(possible) > 0 {
		s.merchant.Cards = append(s.merchant.Cards, lo.Shuffle(possible)[0].ID)
	}
}

func (s *Session) PlayerBuyCard(t string) bool {
	if !lo.Contains(s.merchant.Cards, t) {
		return false
	}

	card, _ := s.GetCard(t)

	if s.GetPlayer().Gold < card.Price {
		return false
	}

	s.UpdatePlayer(func(actor *Actor) bool {
		actor.Gold -= card.Price
		return true
	})

	firstFound := false
	s.merchant.Cards = lo.Filter(s.merchant.Cards, func(item string, index int) bool {
		if firstFound {
			return true
		}

		isType := item == t
		if isType {
			firstFound = true
			return false
		}

		return true
	})
	s.GiveCard(card.ID, PlayerActorID)
	return true
}

func (s *Session) PlayerBuyArtifact(t string) bool {
	if !lo.Contains(s.merchant.Artifacts, t) {
		return false
	}

	art, _ := s.GetArtifact(t)

	if s.GetPlayer().Gold < art.Price {
		return false
	}

	s.UpdatePlayer(func(actor *Actor) bool {
		actor.Gold -= art.Price
		return true
	})

	firstFound := false
	s.merchant.Artifacts = lo.Filter(s.merchant.Artifacts, func(item string, index int) bool {
		if firstFound {
			return true
		}

		isType := item == t
		if isType {
			firstFound = true
			return false
		}

		return true
	})
	s.GiveArtifact(art.ID, PlayerActorID)
	return true
}

//
// StoryTeller
//

func (s *Session) ActiveTeller() *StoryTeller {
	teller := lo.Filter(lo.Values(s.resources.StoryTeller), func(teller *StoryTeller, index int) bool {
		res, err := teller.Active(CreateContext("type_id", teller.ID))
		if err != nil {
			s.log.Printf("Error from Callback:Active type=%s %s\n", teller.ID, err.Error())
			return false
		}
		if val, ok := res.(float64); ok {
			return val > 0
		}
		return false
	})

	if len(teller) == 0 {
		s.log.Printf("No active teller found!")
		return nil
	}

	slices.SortFunc(teller, func(a, b *StoryTeller) bool {
		aOrder, _ := a.Active(CreateContext("type_id", a.ID))
		bOrder, _ := b.Active(CreateContext("type_id", b.ID))

		return aOrder.(float64) > bOrder.(float64)
	})

	return teller[0]
}

func (s *Session) LetTellerDecide() {
	active := s.ActiveTeller()

	if active == nil {
		s.log.Printf("No active teller found! Can't decide")
		return
	}

	res, err := active.Decide(CreateContext("type_id", active.ID))
	if err != nil {
		s.log.Printf("Error from Callback:Decide type=%s %s\n", active.ID, err.Error())
		return
	}

	if val, ok := res.(string); ok {
		s.SetGameState(GameState(val))
	} else {
		s.log.Printf("Error from Callback:Decide type=%s %s\n", active.ID, "return wasn't a game state")
	}
}

//
// Instances
//

func (s *Session) GetInstance(guid string) any {
	return s.instances[guid]
}

func (s *Session) TraverseArtifactsStatus(guids []string, artifact func(instance ArtifactInstance, artifact *Artifact), status func(instance StatusEffectInstance, statusEffect *StatusEffect)) {
	sort.SliceStable(guids, func(i, j int) bool {
		oa := util.Max(s.GetArtifactOrder(guids[i]), s.GetStatusEffectOrder(guids[i]))
		ob := util.Max(s.GetArtifactOrder(guids[j]), s.GetStatusEffectOrder(guids[j]))
		return oa > ob
	})

	for _, id := range guids {
		instance, ok := s.instances[id]
		if !ok {
			continue
		}

		switch instance := instance.(type) {
		case ArtifactInstance:
			// Fetch the backing definition of the type
			art, ok := s.resources.Artifacts[instance.TypeID]
			if !ok {
				continue
			}

			artifact(instance, art)
		case StatusEffectInstance:
			// Fetch the backing definition of the type
			se, ok := s.resources.StatusEffects[instance.TypeID]
			if !ok {
				continue
			}

			status(instance, se)
		}
	}
}

//
// Status Effect Functions
//

func (s *Session) GetStatusEffectOrder(guid string) int {
	// Try as type id
	if e, ok := s.resources.StatusEffects[guid]; ok {
		return e.Order
	}

	instance, ok := s.instances[guid]
	if !ok {
		return 0
	}
	switch instance := instance.(type) {
	case StatusEffectInstance:
		if e, ok := s.resources.StatusEffects[instance.TypeID]; ok {
			return e.Order
		}
	}
	return 0
}

func (s *Session) GetStatusEffect(guid string) *StatusEffect {
	// Try as type id
	if e, ok := s.resources.StatusEffects[guid]; ok {
		return e
	}

	instance, ok := s.instances[guid]
	if !ok {
		return nil
	}
	switch instance := instance.(type) {
	case StatusEffectInstance:
		if e, ok := s.resources.StatusEffects[instance.TypeID]; ok {
			return e
		}
	}
	return nil
}

func (s *Session) GetStatusEffectInstance(guid string) StatusEffectInstance {
	return s.instances[guid].(StatusEffectInstance)
}

func (s *Session) GiveStatusEffect(typeId string, owner string, stacks int) string {
	if len(owner) == 0 {
		s.log.Println("Error: trying to give status effect without owner!")
		return ""
	}

	status := s.resources.StatusEffects[typeId]

	// TODO: This should always be either 0 or 1 len, so the logic down below is a bit meh.
	same := lo.Filter(s.actors[owner].StatusEffects.ToSlice(), func(guid string, index int) bool {
		instance := s.instances[guid].(StatusEffectInstance)
		return instance.TypeID == typeId
	})

	// If it can't stack we delete all existing instances
	if !status.CanStack {
		for i := range same {
			s.RemoveStatusEffect(same[i])
		}
	} else if len(same) > 0 {
		// Increase stack and re-set rounds left
		for i := range same {
			instance := s.instances[same[i]].(StatusEffectInstance)
			instance.Stacks += stacks
			instance.RoundsLeft = status.Rounds
			s.instances[same[i]] = instance

			if _, err := status.Callbacks[CallbackOnStatusStack].Call(CreateContext("type_id", typeId, "guid", same[i], "owner", owner, "stacks", instance.Stacks)); err != nil {
				s.log.Printf("Error from Callback:CallbackOnStatusStack type=%s %s\n", instance.TypeID, err.Error())
			}

			return instance.GUID
		}
	}

	instance := StatusEffectInstance{
		TypeID:     typeId,
		GUID:       NewGuid("STATUS"),
		Owner:      owner,
		RoundsLeft: status.Rounds,
		Stacks:     stacks,
	}
	s.instances[instance.GUID] = instance
	s.actors[owner].StatusEffects.Add(instance.GUID)

	// Call OnStatusAdd callback for the new instance
	_, _ = status.Callbacks[CallbackOnStatusAdd].Call(CreateContext("type_id", typeId, "guid", instance.GUID))

	return instance.GUID
}

func (s *Session) RemoveStatusEffect(guid string) {
	instance := s.instances[guid].(StatusEffectInstance)
	if _, err := s.resources.StatusEffects[instance.TypeID].Callbacks[CallbackOnStatusRemove].Call(CreateContext("type_id", instance.TypeID, "guid", guid, "owner", instance.Owner)); err != nil {
		s.log.Printf("Error from Callback:CallbackOnStatusRemove type=%s %s\n", instance.TypeID, err.Error())
	}
	if actor, ok := s.actors[instance.Owner]; ok {
		actor.StatusEffects.Remove(instance.GUID)
	}
	delete(s.instances, guid)
}

func (s *Session) GetActorStatusEffects(guid string) []string {
	return s.actors[guid].StatusEffects.ToSlice()
}

func (s *Session) AddStatusEffectStacks(guid string, stacks int) {
	instance := s.instances[guid].(StatusEffectInstance)
	instance.Stacks += stacks
	if instance.Stacks <= 0 {
		s.RemoveStatusEffect(guid)
	} else {
		s.instances[guid] = instance
	}
}

func (s *Session) SetStatusEffectStacks(guid string, stacks int) {
	instance := s.instances[guid].(StatusEffectInstance)
	instance.Stacks = stacks
	if instance.Stacks <= 0 {
		s.RemoveStatusEffect(guid)
	} else {
		s.instances[guid] = instance
	}
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

func (s *Session) GetRandomArtifactType(maxPrice int) string {
	possible := lo.Filter(lo.Values(s.resources.Artifacts), func(item *Artifact, index int) bool {
		return item.Price < maxPrice
	})
	if len(possible) == 0 {
		return ""
	}
	return lo.Shuffle(possible)[0].ID
}

func (s *Session) GetArtifacts(owner string) []string {
	guids := s.actors[owner].Artifacts.ToSlice()
	sort.Strings(guids)
	return guids
}

func (s *Session) GetArtifact(guid string) (*Artifact, ArtifactInstance) {
	// check if guid is actually typeId
	if val, ok := s.resources.Artifacts[guid]; ok {
		return val, ArtifactInstance{}
	}

	artInstance, ok := s.instances[guid]
	if !ok {
		return nil, ArtifactInstance{}
	}
	switch artInstance := artInstance.(type) {
	case ArtifactInstance:
		if art, ok := s.resources.Artifacts[artInstance.TypeID]; ok {
			return art, artInstance
		}
	}
	return nil, ArtifactInstance{}
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
	if _, err := s.resources.Artifacts[typeId].Callbacks[CallbackOnPickUp].Call(CreateContext("type_id", typeId, "guid", instance.GUID, "owner", owner)); err != nil {
		s.log.Printf("Error from Callback:CallbackOnPickUp type=%s %s\n", instance.TypeID, err.Error())
	}

	return instance.GUID
}

func (s *Session) RemoveArtifact(guid string) {
	instance := s.instances[guid].(ArtifactInstance)
	if _, err := s.resources.Artifacts[instance.TypeID].Callbacks[CallbackOnRemove].Call(CreateContext("type_id", instance.TypeID, "guid", guid, "owner", instance.Owner)); err != nil {
		s.log.Printf("Error from Callback:CallbackOnRemove type=%s %s\n", instance.TypeID, err.Error())
	}
	s.actors[instance.Owner].Artifacts.Remove(instance.GUID)
	delete(s.instances, guid)
}

//
// Card Functions
//

func (s *Session) GetCard(guid string) (*Card, CardInstance) {
	// check if guid is actually typeId
	if val, ok := s.resources.Cards[guid]; ok {
		return val, CardInstance{}
	}

	cardInstance, ok := s.instances[guid]
	if !ok {
		return nil, CardInstance{}
	}
	switch cardInstance := cardInstance.(type) {
	case CardInstance:
		if card, ok := s.resources.Cards[cardInstance.TypeID]; ok {
			return card, cardInstance
		}
	}
	return nil, CardInstance{}
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

func (s *Session) CastCard(guid string, target string) bool {
	if card, instance := s.GetCard(guid); card != nil {
		res, err := card.Callbacks[CallbackOnCast].Call(CreateContext("type_id", card.ID, "guid", guid, "caster", instance.Owner, "target", target, "level", instance.Level))
		if err != nil {
			s.log.Printf("Error from Callback:CallbackOnCast type=%s %s\n", instance.TypeID, err.Error())
		}
		if val, ok := res.(bool); ok {
			return val
		}
	}
	return false
}

func (s *Session) GetCards(owner string) []string {
	guids := s.actors[owner].Cards.ToSlice()
	sort.Strings(guids)
	return guids
}

func (s *Session) GetCardState(guid string) string {
	card, instance := s.GetCard(guid)
	if card == nil {
		return ""
	}

	if card.State == nil {
		return card.Description
	}

	res, err := card.State.Call(CreateContext("type_id", card.ID, "guid", guid, "level", instance.Level, "owner", instance.Owner))
	if err != nil {
		s.log.Printf("Error from Callback:State type=%s %s\n", instance.TypeID, err.Error())
	}

	if res == nil {
		return card.Description
	}
	return res.(string)
}

func (s *Session) PlayerCastHand(i int, target string) error {
	if i >= len(s.currentFight.Hand) {
		return errors.New("hand empty")
	}

	cardId := s.currentFight.Hand[i]

	// Only cast a card if points are available and subtract them.
	if card, _ := s.GetCard(cardId); card != nil {
		if s.currentFight.CurrentPoints < card.PointCost {
			return errors.New("not enough points")
		}

		s.currentFight.CurrentPoints -= card.PointCost
	} else {
		return errors.New("card not exists")
	}

	// Remove from hand.
	s.currentFight.Hand = lo.Reject(s.currentFight.Hand, func(item string, index int) bool {
		return index == i
	})

	// Cast and exhaust if needed.
	if s.CastCard(cardId, target) {
		s.currentFight.Exhausted = append(s.currentFight.Exhausted, cardId)
	} else {
		s.currentFight.Used = append(s.currentFight.Used, cardId)
	}

	return nil
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

func (s *Session) DealDamage(source string, target string, damage int, flat bool) int {
	if val, ok := s.actors[target]; ok {
		guids := lo.Flatten([][]string{
			s.GetActor(source).Artifacts.ToSlice(),
			s.GetActor(target).Artifacts.ToSlice(),
			s.GetActor(target).StatusEffects.ToSlice(),
			s.GetActor(source).StatusEffects.ToSlice(),
		})

		if !flat {
			s.TraverseArtifactsStatus(guids,
				func(instance ArtifactInstance, art *Artifact) {
					res, err := art.Callbacks[CallbackOnDamageCalc].Call(CreateContext("type_id", art.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "damage", damage))
					if err != nil {
						s.log.Printf("Error from Callback:CallbackOnDamageCalc type=%s %s\n", instance.TypeID, err.Error())
					} else if res != nil {
						if newDamage, ok := res.(float64); ok {
							damage = int(newDamage)
						}
					}
				},
				func(instance StatusEffectInstance, se *StatusEffect) {
					res, err := se.Callbacks[CallbackOnDamageCalc].Call(CreateContext("type_id", se.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "stacks", instance.Stacks, "damage", damage))
					if err != nil {
						s.log.Printf("Error from Callback:CallbackOnDamageCalc type=%s %s\n", instance.TypeID, err.Error())
					} else if res != nil {
						if newDamage, ok := res.(float64); ok {
							damage = int(newDamage)
						}
					}
				},
			)
		}

		if source == PlayerActorID {
			s.Log(LogTypeSuccess, fmt.Sprintf("You hit the enemy for %d damage", damage))
		} else if target == PlayerActorID {
			s.Log(LogTypeDanger, fmt.Sprintf("You took %d damage", damage))
		} else {
			s.Log(LogTypeSuccess, fmt.Sprintf("%s took %d damage", val.Name, damage))
		}

		// Negative damage aka heal is not allowed!
		if damage < 0 {
			damage = 0
		}

		// Trigger OnDamage callbacks
		s.TraverseArtifactsStatus(guids,
			func(instance ArtifactInstance, art *Artifact) {
				_, err := art.Callbacks[CallbackOnDamage].Call(CreateContext("type_id", art.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "damage", damage))
				if err != nil {
					s.log.Printf("Error from Callback:CallbackOnDamage type=%s %s\n", instance.TypeID, err.Error())
				}
			},
			func(instance StatusEffectInstance, se *StatusEffect) {
				_, err := se.Callbacks[CallbackOnDamage].Call(CreateContext("type_id", se.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "stacks", instance.Stacks, "damage", damage))
				if err != nil {
					s.log.Printf("Error from Callback:CallbackOnDamage type=%s %s\n", instance.TypeID, err.Error())
				}
			},
		)

		// Re-fetch actor in case the OnDamage callback triggered some kind of damage or healing.
		val = s.actors[target]

		hpLeft := lo.Clamp(val.HP-damage, 0, val.MaxHP)

		// Remove dead non-player actor
		if target != PlayerActorID && hpLeft == 0 {
			s.PushState(map[StateEvent]any{
				StateEventDeath: StateEventDeathData{
					Source: source,
					Target: target,
					Damage: damage,
				},
			})
			s.Log(LogTypeSuccess, fmt.Sprintf("%s died and dropped %d gold!", val.Name, val.Gold))
			s.UpdatePlayer(func(actor *Actor) bool {
				if val.Gold > 0 {
					actor.Gold += val.Gold
					s.PushState(map[StateEvent]any{
						StateEventMoney: StateEventMoneyData{
							Target: target,
							Money:  val.Gold,
						},
					})
				}
				return true
			})
			s.TraverseArtifactsStatus(guids,
				func(instance ArtifactInstance, art *Artifact) {
					_, err := art.Callbacks[CallbackOnActorDie].Call(CreateContext("type_id", art.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "damage", damage))
					if err != nil {
						s.log.Printf("Error from Callback:CallbackOnDamage type=%s %s\n", instance.TypeID, err.Error())
					}
				},
				func(instance StatusEffectInstance, se *StatusEffect) {
					_, err := se.Callbacks[CallbackOnActorDie].Call(CreateContext("type_id", se.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "stacks", instance.Stacks, "damage", damage))
					if err != nil {
						s.log.Printf("Error from Callback:CallbackOnDamage type=%s %s\n", instance.TypeID, err.Error())
					}
				},
			)
			s.RemoveActor(target)
			s.FinishFight()
		} else {
			s.PushState(map[StateEvent]any{
				StateEventDamage: StateEventDamageData{
					Source: source,
					Target: target,
					Damage: damage,
				},
			})
			s.UpdateActor(target, func(actor *Actor) bool {
				actor.HP = hpLeft
				return true
			})
			if target == PlayerActorID {
				if s.GetPlayer().HP == 0 {
					s.SetGameState(GameStateGameOver)
				}
			}
		}

		return damage
	}
	return 0
}

func (s *Session) DealDamageMulti(source string, targets []string, damage int, flat bool) []int {
	return lo.Map(targets, func(guid string, index int) int {
		return s.DealDamage(source, guid, damage, flat)
	})
}

func (s *Session) Heal(source string, target string, heal int, flat bool) int {
	if val, ok := s.actors[target]; ok {
		if !flat {
			s.TraverseArtifactsStatus(lo.Flatten([][]string{
				s.GetActor(source).Artifacts.ToSlice(),
				s.GetActor(target).StatusEffects.ToSlice(),
				s.GetActor(source).StatusEffects.ToSlice(),
			}),
				func(instance ArtifactInstance, art *Artifact) {
					res, err := art.Callbacks[CallbackOnHealCalc].Call(CreateContext("type_id", art.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "heal", heal))
					if err != nil {
						s.log.Printf("Error from Callback:CallbackOnDamageCalc type=%s %s\n", instance.TypeID, err.Error())
					} else if res != nil {
						if newHeal, ok := res.(float64); ok {
							heal = int(newHeal)
						}
					}
				},
				func(instance StatusEffectInstance, se *StatusEffect) {
					res, err := se.Callbacks[CallbackOnHealCalc].Call(CreateContext("type_id", se.ID, "guid", instance.GUID, "source", source, "target", target, "owner", instance.Owner, "stacks", instance.Stacks, "heal", heal))
					if err != nil {
						s.log.Printf("Error from Callback:CallbackOnDamageCalc type=%s %s\n", instance.TypeID, err.Error())
					} else if res != nil {
						if newHeal, ok := res.(float64); ok {
							heal = int(newHeal)
						}
					}
				},
			)
		}

		if target == PlayerActorID {
			s.Log(LogTypeDanger, fmt.Sprintf("You healed %d damage", heal))
		} else {
			s.Log(LogTypeSuccess, fmt.Sprintf("%s healed %d damage", val.Name, heal))
		}

		// Negative heal aka damage is not allowed!
		if heal < 0 {
			heal = 0
		}

		s.UpdateActor(target, func(actor *Actor) bool {
			actor.HP = lo.Clamp(val.HP+heal, 0, val.MaxHP)
			return true
		})

		return heal
	}
	return 0
}

//
// Actor Functions
//

func (s *Session) GetPlayer() Actor {
	return s.actors[PlayerActorID]
}

func (s *Session) UpdatePlayer(update func(actor *Actor) bool) {
	s.UpdateActor(PlayerActorID, update)
}

func (s *Session) GetActor(id string) Actor {
	return s.actors[id]
}

func (s *Session) UpdateActor(id string, update func(actor *Actor) bool) {
	actor := s.GetActor(id)
	if update(&actor) {
		s.actors[id] = actor
	}
}

func (s *Session) ActorAddMaxHP(id string, val int) {
	s.UpdateActor(id, func(actor *Actor) bool {
		actor.MaxHP += val
		return true
	})
}

func (s *Session) AddActor(actor Actor) {
	s.actors[actor.GUID] = actor
}

func (s *Session) AddActorFromEnemy(id string) string {
	if base, ok := s.resources.Enemies[id]; ok {
		actor := NewActor(NewGuid(id))

		actor.TypeID = id
		actor.Name = base.Name
		actor.Description = base.Description
		actor.HP = base.InitialHP
		actor.MaxHP = base.MaxHP

		// Its important we add the actor before any callbacks so that it's instance is available
		// to add cards etc. to!
		s.AddActor(actor)

		if _, err := base.Callbacks[CallbackOnInit].Call(CreateContext("type_id", id, "guid", actor.GUID)); err != nil {
			s.log.Printf("Error from Callback:CallbackOnDamageCalc type=%s %s\n", actor.TypeID, err.Error())
		}

		return actor.GUID
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
		if a.GUID != PlayerActorID {
			deleteActors = append(deleteActors, a.GUID)
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

func (s *Session) GetOpponentByIndex(viewpoint string, i int) Actor {
	switch viewpoint {
	// From the viewpoint of the player we can have multiple enemies.
	case PlayerActorID:
		if len(s.actors) <= 1 {
			return Actor{}
		}

		ids := lo.Filter(lo.Keys(s.actors), func(guid string, index int) bool {
			return guid != PlayerActorID
		})
		sort.Strings(ids)
		if i < 0 || i >= len(ids) {
			return Actor{}
		}

		return s.actors[ids[i]]
	// From the viewpoint of an enemy we only have the player as enemy.
	default:
		return s.actors[PlayerActorID]
	}
}

func (s *Session) GetOpponents(viewpoint string) []Actor {
	return lo.Map(s.GetOpponentGUIDs(viewpoint), func(guid string, index int) Actor {
		return s.actors[guid]
	})
}

func (s *Session) GetOpponentGUIDs(viewpoint string) []string {
	switch viewpoint {
	// From the viewpoint of the player we can have multiple enemies.
	case PlayerActorID:
		guids := lo.Filter(lo.Keys(s.actors), func(guid string, index int) bool {
			return guid != PlayerActorID
		})
		sort.Strings(guids)
		return guids
	// From the viewpoint of an enemy we only have the player as enemy.
	default:
		return []string{PlayerActorID}
	}
}

func (s *Session) GetEnemy(typeId string) *Enemy {
	return s.resources.Enemies[typeId]
}

//
// Gold
//

func (s *Session) GivePlayerGold(amount int) {
	s.UpdatePlayer(func(actor *Actor) bool {
		actor.Gold += amount
		return true
	})
}

//
// Misc Callback
//

func (s *Session) TriggerOnPlayerTurn() {
	s.TraverseArtifactsStatus(lo.Flatten([][]string{
		s.GetPlayer().Artifacts.ToSlice(),
		s.GetPlayer().StatusEffects.ToSlice(),
	}),
		func(instance ArtifactInstance, artifact *Artifact) {
			if _, err := artifact.Callbacks[CallbackOnPlayerTurn].Call(CreateContext("type_id", artifact.ID, "guid", instance.GUID, "owner", instance.Owner, "round", s.GetFightRound())); err != nil {
				s.log.Printf("Error from Callback:CallbackOnPlayerTurn type=%s %s\n", instance.TypeID, err.Error())
			}
		},
		func(instance StatusEffectInstance, statusEffect *StatusEffect) {
			if _, err := statusEffect.Callbacks[CallbackOnPlayerTurn].Call(CreateContext("type_id", statusEffect.ID, "guid", instance.GUID, "owner", instance.Owner, "round", s.GetFightRound(), "stacks", instance.Stacks)); err != nil {
				s.log.Printf("Error from Callback:CallbackOnPlayerTurn type=%s %s\n", instance.TypeID, err.Error())
			}
		},
	)

	lo.ForEach(s.GetOpponents(PlayerActorID), func(actor Actor, index int) {
		if enemy := s.GetEnemy(actor.TypeID); enemy != nil {
			if _, err := enemy.Callbacks[CallbackOnPlayerTurn].Call(CreateContext("type_id", enemy.ID, "guid", actor.GUID, "round", s.GetFightRound())); err != nil {
				s.log.Printf("Error from Callback:CallbackOnPlayerTurn type=%s %s\n", enemy.ID, err.Error())
			}
		}
	})
}

//
// Misc Functions
//

func (s *Session) Log(t LogType, msg string) {
	s.Logs = append(s.Logs, LogEntry{
		Time:    time.Now(),
		Type:    t,
		Message: msg,
	})
}
