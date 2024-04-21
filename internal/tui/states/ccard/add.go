package ccard

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
)

// Implemet State.
var _ tui.State = (*CardAdd)(nil)

// CardAdd, state 11
// Form for credit card adding.
type CardAdd struct {
	inputs      []textinput.Model // model for card inputs.
	focused     int
	cardErr     error // Card validation
	ansver      bool  // Add info message if servier send answer.
	IsAddOk     bool  // Successful registration.
	ansverCode  int   // Servier answer code.
	ansverError error // Servier answer error.

}

// For send error Command in tui.
type (
	errMsg error
)

const (
	ccn = iota
	exp
	cvv
	hld
)

func NewCardAdd() CardAdd {
	var inputs []textinput.Model = make([]textinput.Model, 4)
	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4505 **** **** 1234"
	inputs[ccn].Focus()
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


	return CardAdd{
		inputs:  inputs,
		focused: 0,
		cardErr: nil,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (ca *CardAdd) GetInit() tea.Cmd {
	return textinput.Blink
}

func (ca *CardAdd) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd = make([]tea.Cmd, len(ca.inputs))

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type.String() {
		case "ctrl+c", "esc":
			ca.cleanform()
			m.ChangeState(tui.CcardAdd, tui.MainMenu)
			return m, nil
		case "enter":
			if ca.focused == len(ca.inputs)-1 {
				zap.S().Infof("Text inputs %s  %s", ca.inputs[0].Value(), ca.inputs[1].Value(), ca.inputs[2].Value(), ca.inputs[3].Value())
				//status, err := client.SiteAdd(m.Conf, *m.User, ca.inputs[0].Value(), ca.inputs[1].Value(), ca.inputs[2].Value())
				status := http.StatusInternalServerError
				err := errors.New("test error")
				ca.ansver = true
				ca.ansverCode = status
				ca.ansverError = err
				if status == http.StatusCreated {
					ca.IsAddOk = true
					return m, nil
				}
				zap.S().Infof("Text inputs %d | %w", status, err)

				return m, nil
			}
			ca.nextInput()

		case "shift+tab":
			ca.prevInput()
		case "tab":
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
func (ca *CardAdd) GetView(m *tui.Model) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`
 %s
 %s

 %s  %s
 %s  %s
 
 %s
 %s

 %s
`,
		styles.GopherQuestion.Width(30).Render("Card Number"),
		ca.inputs[ccn].View(),
		styles.GopherQuestion.Width(6).Render("EXP"),
		styles.GopherQuestion.Width(6).Render("CVV"),
		ca.inputs[exp].View(),
		ca.inputs[cvv].View(),
		styles.GopherQuestion.Width(14).Render("First and Last Name"),
		ca.inputs[hld].View(),
		styles.GopherHeader.Render("Continue ->"),
	) + "\n")

	// Client answer checking.
	if ca.ansver {
		if ca.ansverCode == http.StatusCreated {
			b.WriteString(styles.OkStyle1.Render("Debit card add successful: ", m.User.Login))
			b.WriteString("\n\n")
			b.WriteString(styles.GopherQuestion.Render("Press <Enter> to continue... "))
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

	return b.String()
}

// Reset all inputs and form errors.
func (ca *CardAdd) cleanform() {
	ca.ansver = false
	ca.IsAddOk = false
	ca.ansverCode = 0
	ca.ansverError = nil
}

// Validators and help func

// nextInput focuses the next input field
func (ca *CardAdd) nextInput() {
	ca.focused = (ca.focused + 1) % len(ca.inputs)
}

// prevInput focuses the previous input field
func (ca *CardAdd) prevInput() {
	ca.focused--
	// Wrap around
	if ca.focused < 0 {
		ca.focused = len(ca.inputs) - 1
	}
}

// Validator functions to ensure valid input
func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}
