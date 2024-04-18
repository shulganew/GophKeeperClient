package states

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
)

// Set of faunctions with default behavior, contanins blank menu items.

func GetDefaulInit() tea.Cmd {
	return nil
}

// DefaulUpdate has work with main keys across all app.
func GetDefaulUpdate(m tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func GetHeaderView() string {
	return styles.GopherHeader.Render(fmt.Sprintf("GopherKeeper client, build version: 1.0.0, pid %d \n\n", os.Getpid()))

}
func GetHelpView(string) string {
	s := strings.Builder{}
	s.WriteString("\n\n")
	s.WriteString(styles.SubtleStyle.Render("up/down: select"))
	s.WriteString(styles.DotStyle)
	s.WriteString(styles.SubtleStyle.Render("enter: choose"))
	s.WriteString(styles.DotStyle)
	s.WriteString(styles.SubtleStyle.Render("esc: quit"))
	return s.String()
}
