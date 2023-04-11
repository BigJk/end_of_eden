package menus

import "github.com/charmbracelet/lipgloss"

// Colors

var (
	BaseRed        = lipgloss.Color("#ef233c")
	BaseRedDarker  = lipgloss.Color("#9b2226")
	BaseRedDarkest = lipgloss.Color("#2b1b1c")
	BaseWhite      = lipgloss.Color("#ffffff")
	BaseGray       = lipgloss.Color("#cccccc")
	BaseGrayDarker = lipgloss.Color("#363636")
	BaseYellow     = lipgloss.Color("#ffd966")
	BaseGreen      = lipgloss.Color("#80ed99")
)

// Styles

var (
	BoldStyle   = lipgloss.NewStyle().Bold(true)
	BaseText    = lipgloss.NewStyle().Foreground(BaseWhite)
	RedText     = lipgloss.NewStyle().Foreground(BaseRed)
	TitleStyle  = lipgloss.NewStyle().Margin(1, 2).Foreground(BaseRed)
	HeaderStyle = lipgloss.NewStyle().Margin(1, 2).Padding(0, 2).Background(BaseRedDarker)
	ListStyle   = lipgloss.NewStyle().Margin(1, 2)
)
