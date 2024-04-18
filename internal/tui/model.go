package tui

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Result struct {
	Coosen int
}

const NotLoginSt = 0

// Interface for all states selection.
type State interface {
	GetInit() tea.Cmd
	GetUpdate(Model, tea.Msg) (tea.Model, tea.Cmd)
	GetView(Model) string
}

type Model struct {
	Quitting     bool
	CurrentState int
	States       []State
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m Model) Init() tea.Cmd {
	return m.States[0].GetInit()
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.States[0].GetUpdate(m, msg)

}

// The main view, which just calls the appropriate sub-view
func (m Model) View() string {
	return m.States[0].GetView(m)
}