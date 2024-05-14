package gfile

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/states"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
)

// Implemet State.
var _ tui.State = (*GfileMenu)(nil)

// Main gtext menu
type GfileMenu struct {
	Choices      []string
	Choice       int
	BuildVersion string
	BuilDate     string
}

func NewGfileMenu(buildVersion, buildData string) *GfileMenu {
	return &GfileMenu{BuildVersion: buildVersion, BuilDate: buildData, Choices: []string{"Upload secret file data", "List and download"}}
}

// Init is the first function that will be called. It returns an optional
// initial colmand. To not perform an initial colmand return nil.
func (lm *GfileMenu) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	return nil
}

// Main update function.
func (lm *GfileMenu) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.ChangeState(tui.GfileMenu, tui.MainMenu, false, nil)
			return m, nil
		case "down":
			lm.Choice++
			if lm.Choice > len(lm.Choices)-1 {
				lm.Choice = 0
			}
			return m, nil
		case "up":
			lm.Choice--
			if lm.Choice < 0 {
				lm.Choice = len(lm.Choices) - 1
			}
			return m, nil
		case "enter":
			switch lm.Choice {
			// List logins/pw.
			case 0:
				m.ChangeState(tui.GfileMenu, tui.GfileAdd, false, nil)
			// Add NEW.
			case 1:
				m.ChangeState(tui.GfileMenu, tui.GfileList, false, nil)
			}
			return m, nil
		}
	}
	return m, nil
}

// The main view, which just calls the appropriate sub-view
func (lm *GfileMenu) GetView(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString(states.GetHeaderView(lm.BuildVersion, lm.BuilDate))
	s.WriteString(lm.choicesRegister(m))
	s.WriteString(states.GetHelpView())
	return s.String()
}

// Method for working with views/
//
// Choosing menu.
func (lm *GfileMenu) choicesRegister(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(styles.GopherQuestion.Render(m.User.Login, ", uplod files to secret storage or downlod from it:"))
	s.WriteString("\n\n")
	for i := 0; i < len(lm.Choices); i++ {
		s.WriteString(states.Checkbox(lm.Choices[i], lm.Choice == i))
		s.WriteString("\n")
	}
	str := s.String()
	s.Reset()
	return str
}
