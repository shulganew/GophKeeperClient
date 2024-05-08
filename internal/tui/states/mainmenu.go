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
	Choices      []string
	Choice       int
	BuildVersion string
	BuilDate     string
}

func NewMainMenu(buildVersion, buildData string) *MainMenu {
	return &MainMenu{BuildVersion: buildVersion, BuilDate: buildData, Choices: []string{"Sites logins/pw", "Credit cards", "Secret text", "Sectret bin data", "Logout"}}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (mm *MainMenu) GetInit(m *tui.Model, updateID *string) tea.Cmd {
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
				m.ChangeState(tui.MainMenu, tui.SiteMenu, false, nil)
				return m, nil
				// Credit cards.
			case 1:
				m.ChangeState(tui.MainMenu, tui.CardMenu, false, nil)
				return m, nil
				// Goph text.
			case 2:
				m.ChangeState(tui.MainMenu, tui.GtextMenu, false, nil)
				return m, nil
				// Goph files.
			case 3:
				m.ChangeState(tui.MainMenu, tui.GfileMenu, false, nil)
				return m, nil

			case 4:
				err := backup.CleanData()
				if err != nil {
					zap.S().Errorln("Error clean user's tmp file: ", err)
				}
				m.ChangeState(tui.MainMenu, tui.NotLoginMenu, false, nil)
				return m, nil
			}
		}
	}

	return m, nil
}

// The main view, which just calls the appropriate sub-view
func (mm *MainMenu) GetView(m *tui.Model) string {
	s := strings.Builder{}
	s.WriteString(GetHeaderView(mm.BuildVersion, mm.BuilDate))
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
	b := strings.Builder{}
	b.WriteString("\n")
	b.WriteString(styles.GopherQuestion.Render("Hello, ", m.User.Login, ", choose your secters:"))
	b.WriteString("\n\n")
	for i := 0; i < len(mm.Choices); i++ {
		b.WriteString(Checkbox(mm.Choices[i], mm.Choice == i))
		b.WriteString("\n")
		// Add Help and logout Separator
		if len(mm.Choices)-2 == i {
			b.WriteString("\n")
		}
	}
	str := b.String()
	b.Reset()
	return str
}
