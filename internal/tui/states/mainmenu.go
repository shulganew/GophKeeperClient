package states

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/backup"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

// Implemet State.
var _ tui.State = (*MainMenu)(nil)

// Main menu for log in users, state 3
type MainMenu struct {
	Choices []string
	Choice  int
}

func NewMainMenu() MainMenu {
	return MainMenu{Choices: []string{"Sites logins/pw", "Credit cards", "Secret text", "Sectret bin data", "Help", "Logout"}}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (mm *MainMenu) GetInit() tea.Cmd {
	return nil
}

// Main update function.
func (mm *MainMenu) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
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
			// Site's login and passes.
			case 0:
				m.ChangeState(tui.MainMenu, tui.LoginMenu)
				return m, nil
				// Credit cards.
			case 1:
				m.ChangeState(tui.MainMenu, tui.MainMenu)
				return m, nil
				// TODO
				// Logout - delte userf from model, clean backup data, go to login menu
			case 5:
				err := backup.CleanData()
				if err != nil {
					zap.S().Errorln("Error clean user's tmp file: ", err)
				}
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
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
// Choosing menu.
func (mm *MainMenu) choicesRegister(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(styles.GopherQuestion.Render("Hello, ", m.User.Login, ", choose your secters:"))
	s.WriteString("\n\n")
	for i := 0; i < len(mm.Choices); i++ {
		s.WriteString(Checkbox(mm.Choices[i], mm.Choice == i))
		s.WriteString("\n")
		// Add Help and logout Separator
		if len(mm.Choices)-3 == i {
			s.WriteString("\n")
		}
	}
	return s.String()
}
