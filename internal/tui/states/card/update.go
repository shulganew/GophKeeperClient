package card

import (
	"errors"
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
var _ tui.State = (*CardUpdate)(nil)

// CardUpdate, state 11
// Form for credit card adding.
type CardUpdate struct {
	inputs      []textinput.Model // model for card inputs.
	focused     int
	updateID    string
	cardErr     error // Card validation
	ansver      bool  // Add info message if servier send answer.
	IsAddOk     bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.
}

// For send error Command in tui.

func NewCardUpdate() *CardUpdate {
	var inputs []textinput.Model = make([]textinput.Model, 5)
	inputs[def] = textinput.New()
	inputs[def].Placeholder = "My bank"
	inputs[def].Focus()
	inputs[def].CharLimit = 40
	inputs[def].Width = 30
	inputs[def].Prompt = ""

	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4505 **** **** 1234"
	inputs[ccn].CharLimit = 20
	inputs[ccn].Width = 30
	inputs[ccn].Prompt = ""
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "MM/YY "
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Prompt = ""
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 3
	inputs[cvv].Width = 5
	inputs[cvv].Prompt = ""
	inputs[cvv].Validate = cvvValidator

	inputs[hld] = textinput.New()
	inputs[hld].Placeholder = "Card Holder"
	inputs[hld].CharLimit = 30
	inputs[hld].Width = 30
	inputs[hld].Prompt = ""

	return &CardUpdate{
		inputs:  inputs,
		focused: 0,
		cardErr: nil,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (ca *CardUpdate) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	if updateID != nil {
		ca.updateID = *updateID
	} else {
		ca.ansver = true
		ca.ansverError = errors.New("can't find update id")
	}
	// Init fields.
	card := m.Cards[ca.updateID]
	ca.inputs[def].SetValue(card.Definition)
	ca.inputs[def].Focus()
	ca.inputs[ccn].SetValue(card.Ccn)
	ca.inputs[exp].SetValue(card.Exp)
	ca.inputs[hld].SetValue(card.Hld)
	return textinput.Blink
}

func (ca *CardUpdate) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd = make([]tea.Cmd, len(ca.inputs))

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type.String() {
		case "ctrl+c", "esc":
			ca.cleanform()
			m.ChangeState(tui.CardUpdate, tui.CardMenu, false, nil)
			return m, nil
		case "enter":
			// If user adding done success, enter for continue...
			if ca.IsAddOk {
				ca.cleanform()
				m.ChangeState(tui.CardUpdate, tui.CardList, false, nil)
				return m, nil
			}
			if ca.focused == len(ca.inputs)-1 {
				zap.S().Infof("Text inputs %s  %s", ca.inputs[0].Value(), ca.inputs[1].Value(), ca.inputs[2].Value(), ca.inputs[3].Value(), ca.inputs[4].Value())
				status, err := client.CardsUpdate(m.Client, m.Conf, m.JWT, ca.updateID, ca.inputs[0].Value(), ca.inputs[1].Value(), ca.inputs[2].Value(), ca.inputs[3].Value(), ca.inputs[4].Value())
				ca.ansver = true
				ca.ansverCode = status
				ca.ansverError = err
				if status == http.StatusOK {
					ca.IsAddOk = true
					return m, nil
				}
				zap.S().Infof("Text inputs %d | %w", status, err)

				return m, nil
			}
			ca.nextInput()

		case "shift+tab", "up":
			ca.prevInput()
		case "tab", "down":
			ca.nextInput()
		}
		for i := range ca.inputs {
			ca.inputs[i].Blur()
		}
		ca.inputs[ca.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		ca.cardErr = msg
		return m, nil
	}

	for i := range ca.inputs {
		ca.inputs[i], cmds[i] = ca.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

// The main view, which just calls the appropriate sub-view
func (ca *CardUpdate) GetView(m *tui.Model) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf(
		`
 %s
 %s

 %s
 %s

 %s  %s
 %s  %s
 
 %s
 %s

 %s
`,
		styles.CardAdd.Width(30).Render("Card Number"),
		ca.inputs[def].View(),
		styles.CardAdd.Width(30).Render("Card Number"),
		ca.inputs[ccn].View(),
		styles.CardAdd.Width(6).Render("EXP"),
		styles.CardAdd.Width(6).Render("CVV"),
		ca.inputs[exp].View(),
		ca.inputs[cvv].View(),
		styles.CardAdd.Width(14).Render("First and Last Name"),
		ca.inputs[hld].View(),
		styles.GopherHeader.Render("Continue ->"),
	) + "\n")

	// Client answer checking.
	if ca.ansver {
		if ca.ansverCode == http.StatusCreated {
			b.WriteString(styles.OkStyle1.Render("Debit card add successful: ", m.User.Login))
			b.WriteString("\n\n")
			b.WriteString(styles.CardAdd.Render("Press <Enter> to continue... "))
			b.WriteString("\n\n")
		} else {
			b.WriteString(styles.ErrorStyle.Render("Server ansver with code: ", fmt.Sprint(ca.ansverCode)))
			b.WriteString("\n\n")
		}
		if ca.ansverError != nil {
			b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", ca.ansverError.Error())))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("<tab> - next input, <shift+tab> - previous, <Esc> - back to menu."))
	str := b.String()
	b.Reset()
	return str
}

// Reset all inputs and form errors.
func (ca *CardUpdate) cleanform() {
	ca.ansver = false
	ca.IsAddOk = false
	ca.ansverCode = 0
	ca.ansverError = nil
}

// Validators and help func

// nextInput focuses the next input field
func (ca *CardUpdate) nextInput() {
	ca.focused = (ca.focused + 1) % len(ca.inputs)
}

// prevInput focuses the previous input field
func (ca *CardUpdate) prevInput() {
	ca.focused--
	// Wrap around
	if ca.focused < 0 {
		ca.focused = len(ca.inputs) - 1
	}
}
