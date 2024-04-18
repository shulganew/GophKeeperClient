package states

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

// Implemet State.
var _ tui.State = (*RegisterForm)(nil)

// RegisterForm, state 2
// Inputs: login, email, pw, pw (check corret input)
type RegisterForm struct {
	focusIndex  int
	Inputs      []textinput.Model
	cursorMode  cursor.Mode
	ansver      bool  // Add info message if servier send answer.
	IsRegOk     bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.
}

func NewRegisterForm() RegisterForm {
	rf := RegisterForm{
		Inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range rf.Inputs {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
			t.SetValue("igor")
		case 1:
			t.Placeholder = "e-mail"
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
func (rf *RegisterForm) GetInit() tea.Cmd {
	return textinput.Blink
}

func (rf *RegisterForm) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.ChanegeState(tui.LoginForm, tui.NotLoginState)
			return m, nil

		// Change cursor mode
		case "ctrl+r":
			rf.cursorMode++
			if rf.cursorMode > cursor.CursorHide {
				rf.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(rf.Inputs))
			for i := range rf.Inputs {
				cmds[i] = rf.Inputs[i].Cursor.SetMode(rf.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			// Submit button pressed!
			if s == "enter" && rf.focusIndex == len(rf.Inputs) {
				// If user registration done, enter for continue...
				if rf.IsRegOk {
					m.ChanegeState(tui.SignUpForm, tui.NotLoginState)
				}
				zap.S().Infof("Text inputs %s  %s", rf.Inputs[0].Value(), rf.Inputs[1].Value(), rf.Inputs[2].Value())
				user, status, err := client.UserReg(m.Conf, rf.Inputs[0].Value(), rf.Inputs[1].Value(), rf.Inputs[2].Value())
				rf.ansver = true
				rf.ansverCode = status
				rf.ansverError = err
				if status == http.StatusOK {
					rf.IsRegOk = true
					m.User = *user

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
func (rf *RegisterForm) GetView(m *tui.Model) string {
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
			b.WriteString(styles.GopherQuestion.Render("User registerd successful: ", m.User.Login))
			b.WriteString("\n\n")
			b.WriteString(styles.GopherQuestion.Render("Press <Enter> to continue... "))
			b.WriteString("\n\n")
		} else {
			b.WriteString(styles.GopherQuestion.Render("Server ansver with code: ", fmt.Sprint(rf.ansverCode)))
			b.WriteString("\n\n")
		}

	}
	if rf.ansverError != nil {
		b.WriteString(styles.GopherQuestion.Render(fmt.Sprintf("Error: %s", rf.ansverError.Error())))
		b.WriteString("\n")
	}
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("cursor mode is "))
	b.WriteString(styles.CursorModeHelpStyle.Render(rf.cursorMode.String()))
	b.WriteString(styles.HelpStyle.Render(" (<ctrl+r> to change style), <Esc> - back to menu."))

	return b.String()
}

// Help functions
func (rf *RegisterForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(rf.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range rf.Inputs {
		rf.Inputs[i], cmds[i] = rf.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
