package states

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
)

// Implemet State.
var _ tui.State = (*MainMenu)(nil)

// Main menu for log in users, state 3
type MainMenu struct {
	Choices []string
	Choice  int
}

func NewMainMenu() MainMenu {
	return MainMenu{Choices: []string{"Sites logins/pw", "Credit cards", "Secret text", "Sectret bin data", "Logout"}}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (mm *MainMenu) GetInit() tea.Cmd {
	return nil
}

// Main update function.
func (mm *MainMenu) GetUpdate(m *tui.Model, msg tea.Msg) (tm tea.Model, tcmd tea.Cmd) {
	// Add header.

	mm.updateChoices(m, msg)
	tm, tcmd = GetDefaulUpdate(m, msg)
	return tm, tcmd
}

// The main view, which just calls the appropriate sub-view
func (mm *MainMenu) GetView(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString(GetHeaderView())
	if m.Quitting {
		s.WriteString("\n  See you later!\n\n")
	}

	s.WriteString(mm.choicesRegister(m))

	s.WriteString(GetHelpView())
	return s.String()
}

// Method for working with views.
//
// Update loop for the first view where you're choosing a task.
func (mm *MainMenu) updateChoices(m *tui.Model, msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			mm.Choice++
			if mm.Choice > len(mm.Choices)-1 {
				mm.Choice = 0
			}
		case "up":
			mm.Choice--
			if mm.Choice < 0 {
				mm.Choice = len(mm.Choices) - 1
			}
		case "enter":
			switch mm.Choice {
			// Log in
			case 0:
				m.ChangeState(tui.NotLoginState, tui.LoginForm)
				// Sign up
			case 1:
				m.ChangeState(tui.NotLoginState, tui.SignUpForm)
			}

		}
	}

}

// Choosing menu.
func (mm *MainMenu) choicesRegister(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(styles.GopherQuestion.Render("Hello, ", m.User.Login, ", choose your secters:"))
	s.WriteString("\n\n")
	for i := 0; i < len(mm.Choices); i++ {
		s.WriteString(Checkbox(mm.Choices[i], mm.Choice == i))
		s.WriteString("\n")
	}

	return s.String()
}
