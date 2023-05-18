package style

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Colors

const (
	BaseRedHex       = "#ef233c"
	BaseRedDarkerHex = "#9b2226"
)

var (
	BaseRed        = lipgloss.Color(BaseRedHex)
	BaseRedDarker  = lipgloss.Color(BaseRedDarkerHex)
	BaseRedDarkest = lipgloss.Color("#2b1b1c")
	BaseWhite      = lipgloss.Color("#ffffff")
	BaseGray       = lipgloss.Color("#cccccc")
	BaseGrayDarker = lipgloss.Color("#363636")
	BaseYellow     = lipgloss.Color("#ffd966")
	BaseGreen      = lipgloss.Color("#80ed99")

	TableStyle = func() table.Styles {
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(BaseGrayDarker).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color(BaseWhite)).
			Background(lipgloss.Color(BaseRedDarker)).
			Bold(false)
		return s
	}()
)

// Styles

var (
	BoldStyle      = lipgloss.NewStyle().Bold(true)
	BaseText       = lipgloss.NewStyle().Foreground(BaseWhite)
	RedText        = lipgloss.NewStyle().Foreground(BaseRed)
	GreenText      = lipgloss.NewStyle().Foreground(BaseGreen)
	GrayText       = lipgloss.NewStyle().Foreground(BaseGray)
	GrayTextDarker = lipgloss.NewStyle().Foreground(BaseGrayDarker)
	RedDarkerText  = lipgloss.NewStyle().Foreground(BaseRedDarker)
	TitleStyle     = lipgloss.NewStyle().Margin(1, 2).Foreground(BaseRed)
	HeaderStyle    = lipgloss.NewStyle().Margin(1, 2).Padding(0, 2).Background(BaseRedDarker)
	ListStyle      = lipgloss.NewStyle().Margin(1, 2)
)
