package gfile

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
var _ tui.State = (*FileAdd)(nil)

const InputsGfile = 2

// FileAdd
// Inputs: login, email, pw, pw (check corret input)
type FileAdd struct {
	focusIndex  int
	Inputs      []textinput.Model
	ansver      bool  // Add info message if servier send answer.
	IsAddOk     bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.
}

func NewFileAdd() *FileAdd {
	rf := FileAdd{
		Inputs: make([]textinput.Model, InputsGfile),
	}

	var t textinput.Model
	for i := range rf.Inputs {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Description"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
			t.SetValue("My sectet file")
		case 1:
			t.Placeholder = "/home/myfile.txt"
			t.PromptStyle = styles.NoStyle
			t.TextStyle = styles.NoStyle
			t.SetValue("/home/igor/gfile.txt")

		}
		rf.Inputs[i] = t
	}
	return &rf
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (rf *FileAdd) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	return textinput.Blink
}

func (rf *FileAdd) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			rf.cleanform()
			if m.IsUserLogedIn {
				m.ChangeState(tui.GfileAdd, tui.NotLoginMenu, false, nil)
				return m, nil
			}
			m.ChangeState(tui.GfileAdd, tui.GfileMenu, false, nil)
			return m, nil

		case "tab", "shift+tab", "enter", "up", "down":
			// Clean shown errors in menu.
			rf.ansver = false
			s := msg.String()
			// If user adding done success, enter for continue...
			if rf.IsAddOk {
				rf.cleanform()
				m.ChangeState(tui.GfileAdd, tui.GfileMenu, false, nil)
				return m, nil
			}
			// Submit button pressed!
			if s == "enter" && rf.focusIndex == len(rf.Inputs) {
				zap.S().Infof("Text inputs %s  %s", rf.Inputs[0].Value(), rf.Inputs[1].Value())

				status, err := client.FileAdd(m.Client, m.Conf, m.JWT, rf.Inputs[0].Value(), rf.Inputs[1].Value())
				rf.ansver = true
				rf.ansverCode = status
				rf.ansverError = err
				if status == http.StatusCreated {
					rf.IsAddOk = true
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
func (rf *FileAdd) GetView(m *tui.Model) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.GopherQuestion.Render("Add new file to secret storage:\n"))
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
		if rf.ansverCode == http.StatusCreated {
			b.WriteString(styles.OkStyle1.Render("File successfuly added: ", m.User.Login))
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
	b.WriteString(styles.HelpStyle.Render("<Esc> - back to menu."))

	str := b.String()
	b.Reset()
	return str
}

// Help functions
func (rf *FileAdd) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(rf.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range rf.Inputs {
		rf.Inputs[i], cmds[i] = rf.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Reset all inputs and form errors.
func (rf *FileAdd) cleanform() {
	rf.ansver = false
	rf.IsAddOk = false
	rf.ansverCode = 0
	rf.ansverError = nil
}
