package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const DotChar = " • "

// General stuff for styling the view
var (
	// Meny (Status 0)
	SubtleStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	CheckboxStyleSelected = lipgloss.NewStyle().MarginLeft(3).Foreground(lipgloss.Color("#FDDD00"))
	CheckboxStyle         = lipgloss.NewStyle().MarginLeft(3).Foreground(lipgloss.Color("#00A29C"))
	GopherHeader          = lipgloss.NewStyle().MarginLeft(5).Foreground(lipgloss.Color("#00ADD8"))
	GopherQuestion        = lipgloss.NewStyle().Foreground(lipgloss.Color("#5DC9E2"))
	DotStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(DotChar)

	// Eneter form (Status 1, 2)
	CursorStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#00A29C"))
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#00A29C"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle.Copy()
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	FocusedButton       = FocusedStyle.Copy().Render("[ Submit ]")
	BlurredButton       = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))

	// Ok and Error styles

	OkStyle1   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FDDD00"))
	OkStyle2   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00A29C"))
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CE3262"))
)
