<p align="center">
  <img align="center" src="./github/header.png" />
</p>

[![Discord](https://img.shields.io/discord/1099310842564059168?label=discord)](https://discord.gg/XpDvfvVuB2)

> Welcome to a world 500 years in the future, ravaged by climate change and nuclear wars. The remaining humans have become few and far between, replaced by mutated and plant-based creatures. In this gonzo-fantasy setting, you find yourself awakening from cryo sleep in an underground facility, long forgotten and alone. With all other cryosleep capsules broken, it's up to you to navigate this strange and dangerous world and uncover the secrets that led to your isolation...

**End of Eden...**
- Is a "Slay the Spire"-like, roguelite deck-builder game running fully in console
- Collect Artifacts that give you new cards or various passive effects
- Clash with strange beings and try to survive as long as possible


# Screenshots

![Screenshot](./github/screenshot.png)
![Screenshot](./github/screenshot_merchant.png)


# How to play

- You can play end_of_eden in
  - Console: ``end_of_eden``
  - Browser: ``end_of_eden_browser``
  - Over SSH: ``end_of_eden_ssh``

## Console

A modern console is required to support all the features like full mouse control. Just start the  ``end_of_eden(.exe)`` executable in your terminal.

### Tested Terminals
| Terminal                                              |   OS    | Status             | Note                                                            |
|-------------------------------------------------------|---------|--------------------|-----------------------------------------------------------------|
| **[terminal](https://github.com/microsoft/terminal)** | windows | :white_check_mark: |                                                                 |
| **cmd**                                               | windows | :warning:          | no mouse motion support, mouse clicks and everything else works |
| **[iterm2](https://iterm2.com/)**                     | osx     | :white_check_mark: |                                                                 |


# Tech

## Lua & Modding

Lua is used to define artifacts, cards, enemies and everything else that is dynamic in the game. This makes End of Eden easily extendable. If you want to create mods or learn more about lua:

- See [Lua Documentation](docs/LUA_DOCS.md)

## Building

- You need golang ``>= 1.20`` installed
- Build binary: ``cd ./cmd/game && go build``
- Run without building binary: ``go run ./cmd/game/``
- **Important:** The games working directory needs to be where the ``./assets`` folder is available!

# Credits

- Thanks to **Huw Millward** for the face data published in [Warsim Generator Toolbox](https://huw2k8.itch.io/warsims-generator-toolbox)
- [Interface Beep Sounds](https://bleeoop.itch.io/interface-bleeps) by **Bleeoop**
- [512 Sound Effect Pack](https://opengameart.org/content/512-sound-effects-8-bit-style) by **Juhani Junkala**
- Music and additional audio work by [synthroton](https://synthroton.bandcamp.com/)

# License

- **Code:** licensed under MIT
- **Assets:** See README.md in corresponding folder
