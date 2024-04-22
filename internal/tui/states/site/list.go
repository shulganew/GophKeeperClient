package site

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

// Implemet State.
var _ tui.State = (*ListLogin)(nil)

const InputsSitesList = 3

// ListLogin, state 2
// Inputs: login, email, pw, pw (check corret input)
type ListLogin struct {
	focusIndex  int
	Inputs      []textinput.Model
	ansver      bool  // Add info message if servier send answer.
	IsRegOk     bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.
}

func NewListLogin() ListLogin {
	ll := ListLogin{
		Inputs: make([]textinput.Model, InputsSitesLogin),
	}

	var t textinput.Model
	for i := range ll.Inputs {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "https://mysite.ru"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
			t.SetValue("https://mysite.ru")
		case 1:
			t.Placeholder = "login"
			t.PromptStyle = styles.NoStyle
			t.TextStyle = styles.NoStyle
			t.SetValue("scaevol@yandex.ru")
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.SetValue("123")
		}
		ll.Inputs[i] = t
	}
	return ll
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (ll *ListLogin) GetInit() tea.Cmd {
	return textinput.Blink
}

func (ll *ListLogin) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			ll.cleanform()
			if m.IsUserLogedIn {
				m.ChangeState(tui.ListLogin, tui.NotLoginMenu)
				return m, nil
			}
			m.ChangeState(tui.ListLogin, tui.MainMenu)
			return m, nil
		case "insert":
			// Hide or show password.
			if ll.Inputs[2].EchoMode == textinput.EchoPassword {
				ll.Inputs[2].EchoMode = textinput.EchoNormal
			} else {
				ll.Inputs[2].EchoMode = textinput.EchoPassword
			}
			return m, nil
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			// Clean shown errors in menu.
			ll.ansver = false
			s := msg.String()
			// If user registration done, enter for continue...
			if ll.IsRegOk {
				ll.cleanform()
				m.ChangeState(tui.SignUpForm, tui.MainMenu)
				return m, nil
			}
			// Submit button pressed!
			if s == "enter" && ll.focusIndex == len(ll.Inputs) {

				zap.S().Infof("Text inputs %s  %s", ll.Inputs[0].Value(), ll.Inputs[1].Value(), ll.Inputs[2].Value())
				user, status, err := client.UserReg(m.Conf, ll.Inputs[0].Value(), ll.Inputs[1].Value(), ll.Inputs[2].Value())
				ll.ansver = true
				ll.ansverCode = status
				ll.ansverError = err
				if status == http.StatusOK {
					ll.IsRegOk = true
					m.User = user
					return m, nil
				}
				zap.S().Infof("Text inputs %d | %w", status, err)

				return m, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				ll.focusIndex--
			} else {
				ll.focusIndex++
			}

			if ll.focusIndex > len(ll.Inputs) {
				ll.focusIndex = 0
			} else if ll.focusIndex < 0 {
				ll.focusIndex = len(ll.Inputs)
			}

			cmds := make([]tea.Cmd, len(ll.Inputs))
			for i := 0; i <= len(ll.Inputs)-1; i++ {
				if i == ll.focusIndex {
					// Set focused state
					cmds[i] = ll.Inputs[i].Focus()
					ll.Inputs[i].PromptStyle = styles.FocusedStyle
					ll.Inputs[i].TextStyle = styles.FocusedStyle
					continue
				}
				// Remove focused state
				ll.Inputs[i].Blur()
				ll.Inputs[i].PromptStyle = styles.NoStyle
				ll.Inputs[i].TextStyle = styles.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := ll.updateInputs(msg)

	return m, cmd
}

// The main view, which just calls the appropriate sub-view
func (ll *ListLogin) GetView(m *tui.Model) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.GopherQuestion.Render("Add new site URL, login and password:\n"))
	b.WriteString("\n")
	for i := range ll.Inputs {
		b.WriteString(ll.Inputs[i].View())
		if i < len(ll.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.BlurredButton
	if ll.focusIndex == len(ll.Inputs) {
		button = &styles.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	if ll.ansver {
		if ll.ansverCode == http.StatusOK {
			b.WriteString(styles.OkStyle1.Render("User registerd successful: ", m.User.Login))
			b.WriteString("\n\n")
			b.WriteString(styles.GopherQuestion.Render("Press <Enter> to continue... "))
			b.WriteString("\n\n")
		} else {
			b.WriteString(styles.ErrorStyle.Render("Server ansver with code: ", fmt.Sprint(ll.ansverCode)))
			b.WriteString("\n\n")
		}
		if ll.ansverError != nil {
			b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", ll.ansverError.Error())))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("<Insert> - show or hide password, <Esc> - back to menu."))

	return b.String()
}

// Help functions
func (ll *ListLogin) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(ll.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range ll.Inputs {
		ll.Inputs[i], cmds[i] = ll.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Reset all inputs and form errors.
func (ll *ListLogin) cleanform() {
	ll.ansver = false
	ll.IsRegOk = false
	ll.ansverCode = 0
	ll.ansverError = nil
}
