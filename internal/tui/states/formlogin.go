package states

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

// Implemet State.
var _ tui.State = (*LoginForm)(nil)

// LoginForm, state 1
type LoginForm struct {
	focusIndex int
	Inputs     []textinput.Model
	cursorMode cursor.Mode
}

func NewLoginForm() LoginForm {
	lf := LoginForm{
		Inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range lf.Inputs {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		lf.Inputs[i] = t
	}

	return lf
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (lf *LoginForm) GetInit() tea.Cmd {
	return textinput.Blink
}

func (lf *LoginForm) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.ChanegeState(tui.LoginForm, tui.NotLoginState)
			return m, nil

		// Change cursor mode
		case "ctrl+r":
			lf.cursorMode++
			if lf.cursorMode > cursor.CursorHide {
				lf.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(lf.Inputs))
			for i := range lf.Inputs {
				cmds[i] = lf.Inputs[i].Cursor.SetMode(lf.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && lf.focusIndex == len(lf.Inputs) {
				zap.S().Infof("Text inputs %s  %s", lf.Inputs[0].Value(), lf.Inputs[1].Value())

				return m, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				lf.focusIndex--
			} else {
				lf.focusIndex++
			}

			if lf.focusIndex > len(lf.Inputs) {
				lf.focusIndex = 0
			} else if lf.focusIndex < 0 {
				lf.focusIndex = len(lf.Inputs)
			}

			cmds := make([]tea.Cmd, len(lf.Inputs))
			for i := 0; i <= len(lf.Inputs)-1; i++ {
				if i == lf.focusIndex {
					// Set focused state
					cmds[i] = lf.Inputs[i].Focus()
					lf.Inputs[i].PromptStyle = styles.FocusedStyle
					lf.Inputs[i].TextStyle = styles.FocusedStyle
					continue
				}
				// Remove focused state
				lf.Inputs[i].Blur()
				lf.Inputs[i].PromptStyle = styles.NoStyle
				lf.Inputs[i].TextStyle = styles.NoStyle
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

	for i := range lf.Inputs {
		b.WriteString(lf.Inputs[i].View())
		if i < len(lf.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.BlurredButton
	if lf.focusIndex == len(lf.Inputs) {
		button = &styles.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(styles.HelpStyle.Render("cursor mode is "))
	b.WriteString(styles.CursorModeHelpStyle.Render(lf.cursorMode.String()))
	b.WriteString(styles.HelpStyle.Render(" (ctrl+r to change style), esc - back to menu."))

	return b.String()
}

// Help functions
func (lf *LoginForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(lf.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range lf.Inputs {
		lf.Inputs[i], cmds[i] = lf.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
