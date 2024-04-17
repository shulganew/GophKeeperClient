package tuimodel

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const dotChar = " • "

// General stuff for styling the view
var (
	keywordStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyleSelected = lipgloss.NewStyle().MarginLeft(4).Foreground(lipgloss.Color("#FDDD00"))
	checkboxStyle         = lipgloss.NewStyle().MarginLeft(4).Foreground(lipgloss.Color("#00A29C"))
	gopherHeader          = lipgloss.NewStyle().MarginLeft(6).Foreground(lipgloss.Color("#00ADD8"))
	gopherQuestion        = lipgloss.NewStyle().Foreground(lipgloss.Color("#5DC9E2"))
	mainStyle             = lipgloss.NewStyle().MarginLeft(2)
	dotStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
)

type (
	frameMsg struct{}
)

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type Model struct {
	Question string
	Choices  []string
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Quitting bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.

	return updateChoices(msg, m)

}

// The main view, which just calls the appropriate sub-view
func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(gopherHeader.Render(fmt.Sprintf("GopherKeeper client, version: %d \n\n", m.Choice)))

	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s.WriteString(choicesRegister(m))
	} else {
		s.WriteString(chosenView(m))
	}
	return s.String()
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.Choice++
			if m.Choice > len(m.Choices)-1 {
				m.Choice = 0
			}
		case "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = len(m.Choices) - 1
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}
	}

	return m, nil
}

// Sub-views

// The first view, where you're choosing a task
func choicesRegister(m Model) string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(gopherQuestion.Render(m.Question))
	s.WriteString("\n\n")
	for i := 0; i < len(m.Choices); i++ {
		s.WriteString(checkbox(m.Choices[i], m.Choice == i))
		s.WriteString("\n")
	}
	s.WriteString("\n\n")
	s.WriteString(subtleStyle.Render("up/down: select"))
	s.WriteString(dotStyle)
	s.WriteString(subtleStyle.Render("enter: choose"))
	s.WriteString(dotStyle)
	s.WriteString(subtleStyle.Render("esc: quit"))

	return s.String()
}

// The second view, after a task has been chosen
func chosenView(m Model) string {
	return fmt.Sprintf("You choose: %d, is: %s", m.Choice, m.Choices[m.Choice])
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyleSelected.Render("[x] " + label)
	}
	return checkboxStyle.Render(fmt.Sprintf("[ ] %s", label))
}

// Utils
