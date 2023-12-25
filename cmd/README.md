# Executables

# Game

```
End Of Eden :: Game

  -artifacts string
        test artifacts
  -audio
        disable audio (default true)
  -cards string
        test cards
  -enemies string
        test enemies
  -game_state string
        test game state
  -help
        show help
```

# Game Win

- The base game but running in a window

# Game SSH

```
End Of Eden :: SSH Server
Each SSH session creates it's own game session. Modding and audio not supported.

  -bind string
        ip and port to bind to (default ":8273")
  -help
        show help
  -max_inst int
        maximum of game instances (default 10)
  -timeout int
        ssh idle timeout
```

# Environment Variables

- ``EOE_NO_PROTECT=1``: Disables lua safety and kills the program if a lua error is encountered. Good for debugging.
- ``EOE_DEBUG=1``: Enables the debugging api access if a game is started.

## Internal Tools

These tools are found in ``cmd/internal`` and are not meant to be used by the player. These are used to generate documentation and test the game. Can be useful for modders.

# Docs

- Generates the LUA documentation of the game and prints a markdown formatted text to stdout.
- ``go run ./cmd/internal/docs > LUA_API_DOCS.md``

# Tester

The tester is used to test all artifacts, cards and status effects. This can be used to check if the game is still working after a change. Can be embedded in CI and is also useful for modders.

- Runs the ``test`` function of all artifacts, cards and status effects
- Reports the results
- ``go run ./cmd/internal/tester -mods=mod1,mod2,mod3``
- Will exit with a non-zero exit code if any test fails

```
End Of Eden :: Tester
The tester tests all artifacts, cards and status effects based on their test function.

  -help
        show help
  -mods string
        mods to load (e.g. 'my-mod,test-mod,another-mod')

```

# Fuzzy Tester

The fuzzy tester is used to test the game for panics. It will run a game session with a random number of operations and try to trigger a panic. This is useful to find bugs that are not found by the normal tester.

```
End Of Eden :: Fuzzy Tester
The fuzzy tester hits a game session with a random number of operations and tries to trigger a panic.

  -help
        show help
  -mods string
        mods to load and test, separated by ',' (e.g. mod1,mod2,mod3)
  -n int
        number of goroutines (default 1)
  -seed int
        random seed
  -timeout duration
        length of testing (default 1m0s)
```