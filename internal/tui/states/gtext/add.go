package gtext

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

type errMsg error

// Implemet State.
var _ tui.State = (*GtextAdd)(nil)

// GtextAdd, state
type GtextAdd struct {
	textarea    textarea.Model
	err         error
	ansver      bool  // Add info message if servier send answer.
	IsAddOk     bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.

}

func NewGtextAdd() *GtextAdd {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.Focus()

	return &GtextAdd{
		textarea: ti,
		err:      nil,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (ga *GtextAdd) GetInit() tea.Cmd {
	return textinput.Blink
}

func (ga *GtextAdd) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type.String() {
		case "enter":
			if ga.IsAddOk {
				ga.cleanform()
				m.ChangeState(tui.GtextAdd, tui.GtextMenu)
				return m, nil
			}
		case "ctrl+d":
			ga.textarea.Reset()
		case "ctrl+s":
			if ga.IsAddOk {
				ga.cleanform()
				m.ChangeState(tui.GtextAdd, tui.GtextMenu)
				return m, nil
			}

			zap.S().Infof("Text  %s", ga.textarea.Value())
			_, status, err := client.GtextAdd(m.Client, m.Conf, m.JWT, ga.textarea.Value())
			ga.ansver = true
			ga.ansverCode = status
			ga.ansverError = err
			if status == http.StatusCreated {
				ga.IsAddOk = true
				return m, nil
			}
			zap.S().Infof("Text inputs %d | %w", status, err)

			return m, nil

		case "ctrl+c", "esc":
			ga.cleanform()
			m.ChangeState(tui.GtextAdd, tui.GtextMenu)
			return m, nil
		}

		// We handle errors just like any other message
	case errMsg:
		ga.err = msg
		return m, nil
	}

	ga.textarea, cmd = ga.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// The main view, which just calls the appropriate sub-view
func (ga *GtextAdd) GetView(m *tui.Model) string {
	b := strings.Builder{}
	b.WriteString("Enter your secret note:\n\n")
	b.WriteString(ga.textarea.View())
	// Client answer checking.
	if ga.ansver {
		if ga.ansverCode == http.StatusCreated {
			b.WriteString("\n\n")
			b.WriteString(styles.OkStyle1.Render("Note added!"))
			b.WriteString("\n\n")
			b.WriteString(styles.HelpStyle.Render("Press <Enter> to continue... "))
			b.WriteString("\n\n")
		} else {
			b.WriteString("\n\n")
			b.WriteString(styles.ErrorStyle.Render("Server ansver with code: ", fmt.Sprint(ga.ansverCode)))
			b.WriteString("\n\n")
		}
		if ga.ansverError != nil {
			b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", ga.ansverError.Error())))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(styles.SubtleStyle.Render("<ctrl+s> save text"))
	b.WriteString(styles.DotStyle)
	b.WriteString(styles.SubtleStyle.Render("<ctrl+d> clear form"))
	b.WriteString(styles.DotStyle)
	b.WriteString(styles.SubtleStyle.Render("<Esc>: quit"))
	str := b.String()
	b.Reset()
	return str
}

// Reset all inputs and form errors.
func (ga *GtextAdd) cleanform() {
	ga.ansver = false
	ga.IsAddOk = false
	ga.ansverCode = 0
	ga.ansverError = nil
}
