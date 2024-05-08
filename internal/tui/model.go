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

// Update Site
const SiteUpdate = 7

// Card menu.
const CardMenu = 8

// Add credit card.
const CardAdd = 9

// List credit card.
const CardList = 10

// Update cards
const CardUpdate = 11

// Goph text menu.
const GtextMenu = 12

// Add Goph text.
const GtextAdd = 13

// List Goph text.
const GtextList = 14

// List Goph text.
const GtextUpdate = 15

// Goph file menu.
const GfileMenu = 16

// Add file.
const GfileAdd = 17

// List files.
const GfileList = 18

// Interface for all states selection.
type State interface {
	GetInit(m *Model, updateID *string) tea.Cmd
	GetUpdate(*Model, tea.Msg) (tea.Model, tea.Cmd)
	GetView(*Model) string
}

type Model struct {
	Conf          config.Config
	Client        *oapi.Client  //client for request.
	User          *oapi.NewUser // Store user after login or register.
	JWT           string        // Store user current token.
	IsUserLogedIn bool          // Quick check users registration.
	Quitting      bool
	CurrentState  int
	PreviousState int
	States        []State
	Sites         map[string]oapi.Site  // Memory storage of Site data. SiteID - key
	Cards         map[string]oapi.Card  // Memory storage of Cards data. cardID - key
	Gtext         map[string]oapi.Gtext // Memory storage of Text data. gtextID - key
	Gfile         map[string]oapi.Gfile // Memory storage of Files metadata. fileID - key
	IsUpdate      bool                  // Use for mark update states during state switching
	UpdateID      *string               // updateID = siteID or cardID or gtextID depends on update.
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m Model) Init() tea.Cmd {
	return m.States[m.CurrentState].GetInit(&m, nil)
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.States[m.CurrentState].GetUpdate(&m, msg)

}

// The main view, which just calls the appropriate sub-view
func (m Model) View() string {
	return m.States[m.CurrentState].GetView(&m)
}

// State switcher. If state moves to update, it sent bool value and update string id. updateID = siteID or cardID or gtextID depends on update.
func (m *Model) ChangeState(current, next int, isUpdate bool, updateID *string) {
	m.CurrentState = next
	m.PreviousState = current

	// Init current state

	// Check update parameters for values.
	if isUpdate && updateID != nil {
		m.IsUpdate = true
		m.UpdateID = updateID
		// Init updateID and secret type
		m.States[m.CurrentState].GetInit(m, updateID)
	} else {
		m.IsUpdate = false
		m.UpdateID = nil
	}

	// Preloading data to memory model.
	switch m.CurrentState {
	case SiteList:
		sites, status, err := client.SiteList(m.Client, m.Conf, m.JWT)
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
		cards, status, err := client.CardsList(m.Client, m.JWT)
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
	case GtextList:
		gtext, status, err := client.GtextList(m.Client, m.Conf, m.JWT)
		if err != nil {
			zap.S().Errorln("Can't loading user's gtext data: ", err)
			break
		}
		if status != http.StatusOK {
			zap.S().Errorln(errors.New(fmt.Sprintln("Get wrong status: ", status)))
			break
		}
		m.SetGtext(gtext)
		zap.S().Infoln("Set sites from server: ", len(gtext))
	case GfileList:
		gfile, status, err := client.GfileList(m.Client, m.Conf, m.JWT)
		if err != nil {
			zap.S().Errorln("Can't loading user's gtext data: ", err)
			break
		}
		if status != http.StatusOK {
			zap.S().Errorln(errors.New(fmt.Sprintln("Get wrong status: ", status)))
			break
		}
		m.SetGfile(gfile)
		zap.S().Infoln("Set sites from server: ", len(gfile))
	}
}

// Set size, used for interface conformance save.
func (m *Model) SetSites(sites map[string]oapi.Site) {
	m.Sites = sites
}

// Set size, used for interface conformance save.
func (m *Model) SetCards(cards map[string]oapi.Card) {
	m.Cards = cards
}

// Set size, used for interface conformance save.
func (m *Model) SetGtext(gtexts map[string]oapi.Gtext) {
	m.Gtext = gtexts
}

// Set size, used for interface conformance save.
func (m *Model) SetGfile(gfiles map[string]oapi.Gfile) {
	m.Gfile = gfiles
}
