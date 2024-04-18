package states

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
)

// Implemet State.
var _ tui.State = (*LoginForm)(nil)

// LoginForm, state 1
type LoginForm struct {
	Choices []string
	Choice  int
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (lf *LoginForm) GetInit() tea.Cmd {
	return nil
}

func (lf *LoginForm) GetUpdate(m *tui.Model, msg tea.Msg) (tm tea.Model, tcmd tea.Cmd) {
	// Add header.

	lf.updateChoices(m, msg)
	tm, tcmd = GetDefaulUpdate(m, msg)
	return tm, tcmd
}

// The main view, which just calls the appropriate sub-view
func (lf *LoginForm) GetView(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString(GetHeaderView())
	if m.Quitting {
		s.WriteString("\n  See you later!\n\n")
	}

	s.WriteString(lf.choicesRegister())

	s.WriteString(GetHelpView())
	return s.String()
}

// Method for working with views.
//
// Update loop for the first view where you're choosing a task.
func (lf *LoginForm) updateChoices(m *tui.Model, msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			lf.Choice++
			if lf.Choice > len(lf.Choices)-1 {
				lf.Choice = 0
			}
		case "up":
			lf.Choice--
			if lf.Choice < 0 {
				lf.Choice = len(lf.Choices) - 1
			}
		case "enter":
			m.ChanegeState(tui.LoginForm, tui.NotLoginState)

		}
	}

}

// Choosing menu.
func (lf *LoginForm) choicesRegister() string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(styles.GopherQuestion.Render("It will be login form:"))
	s.WriteString("\n\n")
	for i := 0; i < len(lf.Choices); i++ {
		s.WriteString(checkbox(lf.Choices[i], lf.Choice == i))
		s.WriteString("\n")
	}

	return s.String()
}
