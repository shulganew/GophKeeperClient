package loginpw

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
var _ tui.State = (*AddLogin)(nil)

const InputsSitesLogin = 3

// AddLogin, state 2
// Inputs: login, email, pw, pw (check corret input)
type AddLogin struct {
	focusIndex  int
	Inputs      []textinput.Model
	ansver      bool  // Add info message if servier send answer.
	IsRegOk     bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.
}

func NewAddLogin() AddLogin {
	rf := AddLogin{
		Inputs: make([]textinput.Model, InputsSitesLogin),
	}

	var t textinput.Model
	for i := range rf.Inputs {
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
		rf.Inputs[i] = t
	}
	return rf
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (rf *AddLogin) GetInit() tea.Cmd {
	return textinput.Blink
}

func (rf *AddLogin) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			rf.cleanform()
			if m.IsUserLogedIn {
				m.ChangeState(tui.AddLogin, tui.NotLoginMenu)
				return m, nil
			}
			m.ChangeState(tui.AddLogin, tui.MainMenu)
			return m, nil
		case "insert":
			// Hide or show password.
			if rf.Inputs[2].EchoMode == textinput.EchoPassword {
				rf.Inputs[2].EchoMode = textinput.EchoNormal
			} else {
				rf.Inputs[2].EchoMode = textinput.EchoPassword
			}
			return m, nil
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			// Clean shown errors in menu.
			rf.ansver = false
			s := msg.String()
			// If user registration done, enter for continue...
			if rf.IsRegOk {
				rf.cleanform()
				m.ChangeState(tui.SignUpForm, tui.MainMenu)
				return m, nil
			}
			// Submit button pressed!
			if s == "enter" && rf.focusIndex == len(rf.Inputs) {

				zap.S().Infof("Text inputs %s  %s", rf.Inputs[0].Value(), rf.Inputs[1].Value(), rf.Inputs[2].Value())
				user, status, err := client.UserReg(m.Conf, rf.Inputs[0].Value(), rf.Inputs[1].Value(), rf.Inputs[2].Value())
				rf.ansver = true
				rf.ansverCode = status
				rf.ansverError = err
				if status == http.StatusOK {
					rf.IsRegOk = true
					m.User = user
					return m, nil
				}
				zap.S().Infof("Text inputs %d | %w", status, err)

				return m, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				rf.focusIndex--
			} else {
				rf.focusIndex++
			}

			if rf.focusIndex > len(rf.Inputs) {
				rf.focusIndex = 0
			} else if rf.focusIndex < 0 {
				rf.focusIndex = len(rf.Inputs)
			}

			cmds := make([]tea.Cmd, len(rf.Inputs))
			for i := 0; i <= len(rf.Inputs)-1; i++ {
				if i == rf.focusIndex {
					// Set focused state
					cmds[i] = rf.Inputs[i].Focus()
					rf.Inputs[i].PromptStyle = styles.FocusedStyle
					rf.Inputs[i].TextStyle = styles.FocusedStyle
					continue
				}
				// Remove focused state
				rf.Inputs[i].Blur()
				rf.Inputs[i].PromptStyle = styles.NoStyle
				rf.Inputs[i].TextStyle = styles.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := rf.updateInputs(msg)

	return m, cmd
}

// The main view, which just calls the appropriate sub-view
func (rf *AddLogin) GetView(m *tui.Model) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.GopherQuestion.Render("Registration form:\n"))
	b.WriteString("\n")
	for i := range rf.Inputs {
		b.WriteString(rf.Inputs[i].View())
		if i < len(rf.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.BlurredButton
	if rf.focusIndex == len(rf.Inputs) {
		button = &styles.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	if rf.ansver {
		if rf.ansverCode == http.StatusOK {
			b.WriteString(styles.OkStyle1.Render("User registerd successful: ", m.User.Login))
			b.WriteString("\n\n")
			b.WriteString(styles.GopherQuestion.Render("Press <Enter> to continue... "))
			b.WriteString("\n\n")
		} else {
			b.WriteString(styles.ErrorStyle.Render("Server ansver with code: ", fmt.Sprint(rf.ansverCode)))
			b.WriteString("\n\n")
		}
		if rf.ansverError != nil {
			b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", rf.ansverError.Error())))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("<Insert> - show or hide password, <Esc> - back to menu."))

	return b.String()
}

// Help functions
func (rf *AddLogin) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(rf.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range rf.Inputs {
		rf.Inputs[i], cmds[i] = rf.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Reset all inputs and form errors.
func (rf *AddLogin) cleanform() {
	rf.ansver = false
	rf.IsRegOk = false
	rf.ansverCode = 0
	rf.ansverError = nil
}
