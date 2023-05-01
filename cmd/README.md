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

# Game Browser

- The base game but running in browser
- ``ttyd`` is needed as dependency, either in your path or in the sub-folder of the working directory

### ttyd search paths

```
ttyd
./ttyd/ttyd.win32.exe
./ttyd/ttyd.aarch64
./ttyd/ttyd.arm
./ttyd/ttyd.armhf
./ttyd/ttyd.i686
./ttyd/ttyd.mips
./ttyd/ttyd.mips64
./ttyd/ttyd.mips64el
./ttyd/ttyd.mipsel
./ttyd/ttyd.s390x
./ttyd/ttyd.x86_64
```

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

# Docs

- Generates the LUA documentation of the game and prints a markdown formatted text to stdout.
- ``go run ./cmd/docs > LUA_API_DOCS.md``

# Fuzzy Tester

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

# Environment Variables

- ``EOE_NO_PROTECT=1``: Disables lua safety and kills the program if a lua error is encountered. Good for debugging.
- ``EOE_DEBUG=1``: Enables the debugging api access if a game is started. This needs ``ttyd`` and ``wscat`` installed for full usage.