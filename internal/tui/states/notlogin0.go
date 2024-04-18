package states

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

const Question = "You are not authorized, Log In or Sign Up:"

// Implemet State.
var _ tui.State = (*NotLogin)(nil)

type NotLogin struct {
	Choices []string
	Choice  int
	Chosen  bool
}

func (nl *NotLogin) GetInit() tea.Cmd {
	return nil
}

func (nl *NotLogin) GetUpdate(m tui.Model, msg tea.Msg) (tm tea.Model, tcmd tea.Cmd) {
	// Add header.

	nl.updateChoices(msg)
	tm, tcmd = GetDefaulUpdate(m, msg)
	return tm, tcmd
}

func (nl *NotLogin) GetView(m tui.Model) string {
	zap.S().Infoln("Current choice", nl.Choice)

	s := GetHeaderView()
	if m.Quitting {
		return s + "\n  See you later!\n\n"
	}
	if !nl.Chosen {
		return s + nl.choicesRegister()
	} else {
		// Chosen!
		return s + "\n  Chosen!Chosen!!\n\n"
	}

}

// Update loop for the first view where you're choosing a task.
func (nl *NotLogin) updateChoices(msg tea.Msg) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			nl.Choice++
			if nl.Choice > len(nl.Choices)-1 {
				nl.Choice = 0
			}
		case "up":
			nl.Choice--

			if nl.Choice < 0 {
				nl.Choice = len(nl.Choices) - 1
			}
		case "enter":
			nl.Chosen = true

		}
	}

}

// Choosing menu.
func (nl *NotLogin) choicesRegister() string {
	zap.S().Infoln("Update choice", nl.Choice, " ", len(nl.Choices))
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(styles.GopherQuestion.Render(Question))
	s.WriteString("\n\n")
	for i := 0; i < len(nl.Choices); i++ {
		s.WriteString(checkbox(nl.Choices[i], nl.Choice == i))
		s.WriteString("\n")
	}

	return s.String()
}

func checkbox(label string, checked bool) string {
	if checked {
		return styles.CheckboxStyleSelected.Render("[x] " + label)
	}
	return styles.CheckboxStyle.Render(fmt.Sprintf("[ ] %s", label))
}
