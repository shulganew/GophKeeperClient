package states

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
)

// Implemet State.
var _ tui.State = (*NotLogin)(nil)

// Not login menu, state 0
// Menu for registration and log in.
type NotLogin struct {
	Choices []string
	Choice  int
}

func NewNotLogin() NotLogin {
	return NotLogin{Choices: []string{"Log In", "Sign Up"}}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (nl *NotLogin) GetInit() tea.Cmd {
	return nil
}

// Main update function.
func (nl *NotLogin) GetUpdate(m *tui.Model, msg tea.Msg) (tm tea.Model, tcmd tea.Cmd) {
	// Add header.

	nl.updateChoices(m, msg)
	tm, tcmd = GetDefaulUpdate(m, msg)
	return tm, tcmd
}

// The main view, which just calls the appropriate sub-view
func (nl *NotLogin) GetView(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString(GetHeaderView())
	if m.Quitting {
		s.WriteString("\n  See you later!\n\n")
	}

	s.WriteString(nl.choicesRegister())

	s.WriteString(GetHelpView())
	return s.String()
}

// Method for working with views.
//
// Update loop for the first view where you're choosing a task.
func (nl *NotLogin) updateChoices(m *tui.Model, msg tea.Msg) {
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
			switch nl.Choice {
			// Log in
			case 0:
				m.ChangeState(tui.NotLoginMenu, tui.LoginForm)
				// Sign up
			case 1:
				m.ChangeState(tui.NotLoginMenu, tui.SignUpForm)
			}

		}
	}

}

// Choosing menu.
func (nl *NotLogin) choicesRegister() string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(styles.GopherQuestion.Render("You are not authorized, Log In or Sign Up:"))
	s.WriteString("\n\n")
	for i := 0; i < len(nl.Choices); i++ {
		s.WriteString(Checkbox(nl.Choices[i], nl.Choice == i))
		s.WriteString("\n")
	}

	return s.String()
}
