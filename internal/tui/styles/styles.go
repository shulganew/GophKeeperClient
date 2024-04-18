package styles

import "github.com/charmbracelet/lipgloss"

const DotChar = " • "

// General stuff for styling the view
var (
	SubtleStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	CheckboxStyleSelected = lipgloss.NewStyle().MarginLeft(3).Foreground(lipgloss.Color("#FDDD00"))
	CheckboxStyle         = lipgloss.NewStyle().MarginLeft(3).Foreground(lipgloss.Color("#00A29C"))
	GopherHeader          = lipgloss.NewStyle().MarginLeft(5).Foreground(lipgloss.Color("#00ADD8"))
	GopherQuestion        = lipgloss.NewStyle().Foreground(lipgloss.Color("#5DC9E2"))
	DotStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(DotChar)
)
