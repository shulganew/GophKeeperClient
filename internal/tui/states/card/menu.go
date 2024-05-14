package card

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/states"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
)

// Implemet State.
var _ tui.State = (*CardMenu)(nil)

// Main site's login and password administration, state 4
type CardMenu struct {
	Choices []string
	Choice  int
}

func NewCardMenu() *CardMenu {
	return &CardMenu{Choices: []string{"Add NEW card", "List cards"}}
}

// Init is the first function that will be called. It returns an optional
// initial colmand. To not perform an initial colmand return nil.
func (cm *CardMenu) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	return nil
}

// Main update function.
func (cm *CardMenu) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.ChangeState(tui.CardMenu, tui.MainMenu, false, nil)
		case "down":
			cm.Choice++
			if cm.Choice > len(cm.Choices)-1 {
				cm.Choice = 0
			}
			return m, nil
		case "up":
			cm.Choice--
			if cm.Choice < 0 {
				cm.Choice = len(cm.Choices) - 1
			}
			return m, nil
		case "enter":
			switch cm.Choice {
			// List logins/pw.
			case 0:
				m.ChangeState(tui.CardMenu, tui.CardAdd, false, nil)
			// Add NEW.
			case 1:
				m.ChangeState(tui.CardMenu, tui.CardList, false, nil)
			}
			return m, nil
		}
	}
	return m, nil
}

// The main view, which just calls the appropriate sub-view
func (cm *CardMenu) GetView(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString(styles.GopherHeader.Render(fmt.Sprintf("GopherKeeper client, build version: 1.0.0, pid %d \n\n", os.Getpid())))

	s.WriteString(cm.choicesRegister(m))

	s.WriteString(states.GetHelpView())
	return s.String()
}

// Method for working with views/
//
// Choosing menu.
func (cm *CardMenu) choicesRegister(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(styles.GopherQuestion.Render(m.User.Login, ", yours debet cards here:"))
	s.WriteString("\n\n")
	for i := 0; i < len(cm.Choices); i++ {
		s.WriteString(states.Checkbox(cm.Choices[i], cm.Choice == i))
		s.WriteString("\n")
	}
	return s.String()
}
