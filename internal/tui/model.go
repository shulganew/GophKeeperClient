package tui

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/entities"
)

type Result struct {
	Coosen int
}

// Menu for not login User with Log in and Sign up choices
const NotLoginState = 0

// Users login form.
const LoginForm = 1

// User registration form.
const SignUpForm = 2

// Interface for all states selection.
type State interface {
	GetInit() tea.Cmd
	GetUpdate(*Model, tea.Msg) (tea.Model, tea.Cmd)
	GetView(*Model) string
}

type Model struct {
	Conf          config.Config
	User          entities.User // Store user after login or register.
	IsUserSet     bool          // Quick check users registration.
	Quitting      bool
	CurrentState  int
	PreviousState int
	States        []State
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m Model) Init() tea.Cmd {
	return m.States[m.CurrentState].GetInit()
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.States[m.CurrentState].GetUpdate(&m, msg)

}

// The main view, which just calls the appropriate sub-view
func (m Model) View() string {
	return m.States[m.CurrentState].GetView(&m)
}

func (m *Model) ChanegeState(current, next int) {
	m.CurrentState = next
	m.PreviousState = current
}
