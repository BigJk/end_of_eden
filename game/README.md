# Game

This is the core implementation of end_of_eden. This doesn't contain any ui code, as the base game logic and ui is fully seperated in this project.

# Session

The **Session** type is the core of a running game session. It exposes a multitude of methods that are able to mutate the game state. By only exposing methods, calling them should mutate the session from one valid game state to another valid game state. The session will (with a few exceptions) only return copies of data and no pointer so that the consumer can't mutate the state without going through the right method call.

At the moment access to a **Session** is expected to be single-threaded!

```go
session := game.NewSession(/* options */)

// interact with the session
```

# Types

- **Artifact:** Base definition for an artifact
- **Card:** Base definition for a card
- **StatusEffect:** Base definition for a status effect
- **Event:** Base definition for an event
- **Enemy:** Base definition for an enemy
- **StoryTeller:** Base definition for an story teller
- **ArtifactInstance:** Instance of an artifact that is owned by an **Actor** and references its base definition via the ``TypeID``
- **CardInstance:** Instance of a card that is owned by an **Actor** and references its base definition via the ``TypeID``
- **StatusEffectInstance:** Instance of a status effect that is owned by an **Actor** and references its base definition via the ``TypeID``
- **Actor:** Instance of either the player or an enemy. If it's an enemy the ``TypeID`` references its base definition


# Checkpointing

Although the game code doesn't contain any ui there is a mechanism that helps ui and animations to react to events in the game. It is possible to create checkpoints of the game state, execute one or more operations that alter the game state and then get a diff of what happened.

In case we want to see if the enemy turn resulted in damage to the player we could use checkpoints:

```go
// Store state checkpoint before operation
before := session.MarkState()

// Finish the player turn, which means that enemies are allowed to act
session.FinishPlayerTurn()

// Check if any damage was triggered
damages := before.DiffEvent(m.Session, game.StateEventDamage)
for i := range damages {
	// Do something with the damage data -> queue animation, play audio
}
```

This also makes it really easy to keep track of everything that happened in a fight, so that a re-cap screen can be shown. We just need to store the state at the beginning of the fight and diff it when it ends. 