# Executables

# Game

- The base game running in console

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



# Docs

- Generates the LUA documentation of the game and prints a markdown formatted text to stdout.
- ``go run ./cmd/docs > LUA_API_DOCS.md``

# Fuzzy Tester

```
End Of Eden :: Fuzzy Tester
The fuzzy tester hits a game session with a random number of operations and tries to trigger a panic.

  -mods string
        mods to load and test, separated by ',' (e.g. mod1,mod2,mod3)
  -n int
        number of goroutines
  -seed int
        random seed
  -timeout duration
        length of testing (default 1m0s)
```

# Environment Variables

- ``EOE_NO_PROTECT=1``: Disables lua safety and kills the program if a lua error is encountered. Good for debugging.
- ``EOE_DEBUG=1``: Enables the debugging api access if a game is started. This needs ``ttyd`` and ``wscat`` installed for full usage.