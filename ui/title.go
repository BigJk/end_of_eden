package ui

import "github.com/BigJk/end_of_eden/system/localization"

// Title is the title of the game.
var Title = `▄▄▄ . ▐ ▄ ·▄▄▄▄            ·▄▄▄    ▄▄▄ .·▄▄▄▄  ▄▄▄ . ▐ ▄
▀▄.▀·•█▌▐███▪ ██     ▪     ▐▄▄·    ▀▄.▀·██▪ ██ ▀▄.▀·•█▌▐█
▐▀▀▪▄▐█▐▐▌▐█· ▐█▌     ▄█▀▄ ██▪     ▐▀▀▪▄▐█· ▐█▌▐▀▀▪▄▐█▐▐▌
▐█▄▄▌██▐█▌██. ██     ▐█▌.▐▌██▌.    ▐█▄▄▌██. ██ ▐█▄▄▌██▐█▌
 ▀▀▀ ▀▀ █▪▀▀▀▀▀•      ▀█▄▀▪▀▀▀      ▀▀▀ ▀▀▀▀▀•  ▀▀▀ ▀▀ █▪`

// AboutText returns the localized about text.
func AboutText() string {
	return localization.G("ui.about_text", "Welcome to a world 500 years in the future, ravaged by climate change and nuclear wars. The remaining humans have become few and far between, replaced by mutated and plant-based creatures. In this gonzo-fantasy setting, you find yourself awakening from cryo sleep in an underground facility, long forgotten and alone. With all other cryosleep capsules broken, it's up to you to navigate this strange and dangerous world and uncover the secrets that led to your isolation...")
}
