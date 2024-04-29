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
const RegisterForm = 2

// Mani menu for loged in users.
const MainMenu = 3

// Menu for site's logins and passwords.
const SiteMenu = 4

// List site's logins and passwords.
const SiteList = 5

// Add site's logins and passwords.
const SiteAdd = 6

// TODO
const SiteUpdate = 7

// Card menu.
const CardMenu = 8

// Add credit card.
const CardAdd = 9

// List credit card.
const CardList = 10

// List credit card. //TODO
const CardUpdate = 11

// Goph text menu.
const GtextMenu = 12

// Add credit card.
const GtextAdd = 13

// List credit card.
const GtextList = 14

// Goph file menu.
const GfileMenu = 15

// Add credit card.
const GfileAdd = 16

// List credit card.
const GfileList = 17

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
	Sites         []oapi.Site // Memory storage of Site data.
	Cards         []oapi.Card // Memory storage of Cards data.
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
	case CardList:
		cards, status, err := client.CardsList(m.Conf, m.JWT)
		if err != nil {
			zap.S().Errorln("Can't loading user's cards data: ", err)
			break
		}
		if status != http.StatusOK {
			zap.S().Errorln(errors.New(fmt.Sprintln("Get wrong status: ", status)))
			break
		}
		m.SetCards(cards)
		zap.S().Infoln("Set sites from server: ", len(cards))
	}
	

}

// Set size, used for interface conformance save.
func (m *Model) SetSites(sites []oapi.Site) {
	m.Sites = sites
}

// Set size, used for interface conformance save.
func (m *Model) SetCards(cards []oapi.Card) {
	m.Cards = cards
}
