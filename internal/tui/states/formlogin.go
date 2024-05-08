package states

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/backup"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

const inputsLogin = 2

// Implemet State.
var _ tui.State = (*LoginForm)(nil)

// LoginForm, state 1
type LoginForm struct {
	focusIndex  int
	inputs      []textinput.Model
	ansver      bool  // Add info message if servier send answer.
	isLogInOk   bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.
}

func NewLoginForm() *LoginForm {
	lf := LoginForm{
		inputs: make([]textinput.Model, inputsLogin),
	}

	var t textinput.Model
	for i := range lf.inputs {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Login"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
			t.SetValue("igor")
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.SetValue("123")
		}
		lf.inputs[i] = t
	}
	return &lf
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (lf *LoginForm) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	return textinput.Blink
}

func (lf *LoginForm) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			lf.cleanform()
			if m.IsUserLogedIn {
				m.ChangeState(tui.LoginForm, tui.MainMenu, false, nil)
				return m, nil
			}
			m.ChangeState(tui.LoginForm, tui.NotLoginMenu, false, nil)
			return m, nil

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			// Clean shown errors in menu.
			lf.ansver = false
			s := msg.String()
			// Loged in, exit.
			if lf.isLogInOk {
				lf.cleanform()
				m.ChangeState(tui.LoginForm, tui.MainMenu, false, nil)
				return m, nil
			}
			// Submit button pressed!
			if s == "enter" && lf.focusIndex == len(lf.inputs) {
				zap.S().Infof("Text inputs %s  %s", lf.inputs[0].Value(), lf.inputs[1].Value())
				user, jwt, status, err := client.UserLogin(m.Client, lf.inputs[0].Value(), lf.inputs[1].Value())
				lf.ansver = true
				lf.ansverCode = status
				lf.ansverError = err
				if status == http.StatusOK {
					lf.isLogInOk = true
					m.User = user
					m.JWT = jwt
					// Backup curent user.
					err = backup.SaveUser(*user, m.JWT)
					if err != nil {
						zap.S().Errorln("Can't save user: ", err)
					}
					return m, nil

				}
				zap.S().Infof("Text inputs %d | %w", status, err)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				lf.focusIndex--
			} else {
				lf.focusIndex++
			}

			if lf.focusIndex > len(lf.inputs) {
				lf.focusIndex = 0
			} else if lf.focusIndex < 0 {
				lf.focusIndex = len(lf.inputs)
			}

			cmds := make([]tea.Cmd, len(lf.inputs))
			for i := 0; i <= len(lf.inputs)-1; i++ {
				if i == lf.focusIndex {
					// Set focused state
					cmds[i] = lf.inputs[i].Focus()
					lf.inputs[i].PromptStyle = styles.FocusedStyle
					lf.inputs[i].TextStyle = styles.FocusedStyle
					continue
				}
				// Remove focused state
				lf.inputs[i].Blur()
				lf.inputs[i].PromptStyle = styles.NoStyle
				lf.inputs[i].TextStyle = styles.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := lf.updateInputs(msg)
	return m, cmd
}

// The main view, which just calls the appropriate sub-view
func (lf *LoginForm) GetView(m *tui.Model) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.GopherQuestion.Render("Log in form:\n"))
	b.WriteString("\n")
	for i := range lf.inputs {
		b.WriteString(lf.inputs[i].View())
		if i < len(lf.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.BlurredButton
	if lf.focusIndex == len(lf.inputs) {
		button = &styles.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	if lf.ansver {
		if lf.ansverCode == http.StatusOK {
			b.WriteString(styles.OkStyle1.Render("User loged in successful: ", m.User.Login))
			b.WriteString("\n\n")
			b.WriteString(styles.GopherQuestion.Render("Press <Enter> to continue... "))
			b.WriteString("\n\n")
		} else {
			if lf.ansverCode == http.StatusUnauthorized {
				b.WriteString(styles.ErrorStyle.Render("Login or password not correct."))
				b.WriteString("\n")
			}

			b.WriteString(styles.ErrorStyle.Render("Server ansver with code: ", fmt.Sprint(lf.ansverCode)))
			b.WriteString("\n\n")
		}
		if lf.ansverError != nil {
			b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", lf.ansverError.Error())))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("<Esc> - back to menu."))

	str := b.String()
	b.Reset()
	return str
}

// Help functions.
func (lf *LoginForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(lf.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range lf.inputs {
		lf.inputs[i], cmds[i] = lf.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Reset all inputs and form errors.
func (lf *LoginForm) cleanform() {
	lf.ansver = false
	lf.isLogInOk = false
	lf.ansverCode = 0
	lf.ansverError = nil
}
