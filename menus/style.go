package menus

import "github.com/charmbracelet/lipgloss"

// Colors

var BaseRed = lipgloss.Color("#ef233c")
var BaseRedDarker = lipgloss.Color("#9b2226")
var BaseRedDarkest = lipgloss.Color("#2b1b1c")

var BaseWhite = lipgloss.Color("#ffffff")

var BaseGray = lipgloss.Color("#cccccc")
var BaseGrayDarker = lipgloss.Color("#363636")

var BaseYellow = lipgloss.Color("#ffd966")

var BaseGreen = lipgloss.Color("#80ed99")

// Styles

var BoldStyle = lipgloss.NewStyle().Bold(true)
var BaseText = lipgloss.NewStyle().Foreground(BaseWhite)
var TitleStyle = lipgloss.NewStyle().Margin(1, 2).Foreground(BaseRed)
var HeaderStyle = lipgloss.NewStyle().Margin(1, 2).Padding(0, 2).Background(BaseRedDarker)
var ListStyle = lipgloss.NewStyle().Margin(1, 2)
