package tui

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"errors"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

type Result struct {
	Coosen int
}

// Menu for not login User with Log in and Sign up choices
const NotLoginMenu = 0

// Users login form.
const LoginForm = 1

// User registration form.
const SignUpForm = 2

// Mani menu for loged in users.
const MainMenu = 3

// Menu for site's logins and passwords.
const LoginMenu = 4

// List site's logins and passwords.
const SiteList = 5

// Add site's logins and passwords.
const SiteAdd = 6

// Edit site's logins and passwords.
const SiteEdit = 7

// TODO Gand site's logins and passwords to othes users.
const SiteGrand = 8

// TODO Rezerved for Igor's ideas.
const Reserved = 9

// Add credit card.
const CcardAdd = 10

// TODO Add text data to system.
const TextAdd = 15

// TODO Add binary data to system.
const TextBin = 21

// Interface for all states selection.
type State interface {
	GetInit() tea.Cmd
	GetUpdate(*Model, tea.Msg) (tea.Model, tea.Cmd)
	GetView(*Model) string
}

type Model struct {
	Conf          config.Config
	User          *oapi.NewUser // Store user after login or register.
	JWT           string        // Store user current token.
	IsUserLogedIn bool          // Quick check users registration.
	Quitting      bool
	CurrentState  int
	PreviousState int
	States        []State
	Sites         []oapi.Site
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

// State switcher.
func (m *Model) ChangeState(current, next int) {
	m.CurrentState = next
	m.PreviousState = current

	// Preloading data to memory model.
	switch m.CurrentState {
	case SiteList:
		sites, status, err := client.SiteList(m.Conf, m.JWT)
		if err != nil {
			zap.S().Errorln("Can't loading user's site data: ", err)
			break
		}
		if status != http.StatusOK {
			zap.S().Errorln(errors.New(fmt.Sprintln("Get wrong status: ", status)))
			break
		}
		m.SetSites(sites)
		zap.S().Infoln("Set sites from server: ", len(sites))
	}
}

// Set size, used for interface conformance save.
func (m *Model) SetSites(sites []oapi.Site) {
	m.Sites = sites
}
