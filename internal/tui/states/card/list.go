package card

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
	"golang.org/x/term"
)

// Implemet State.
var _ tui.State = (*CardList)(nil)

type Card oapi.Card

func (c Card) Title() string {
	return fmt.Sprintf("%s ◉ %s", c.Definition, c.Ccn)
}
func (c Card) Description() string {
	return fmt.Sprintf("%s ◉ %s ◉ %s", c.Hld, c.Exp, c.Cvv)
}
func (c Card) FilterValue() string { return c.Definition }

type CardList struct {
	list list.Model
}

// CardList, state 10
// List saved site credentials.
func NewCardList() *CardList {
	// Create empty list items
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "My bank cards. <ctrl+u> - update, <ctrl+d> - delete"
	cl := CardList{list: l}
	// Fix terminal bag.
	tw, th, _ := term.GetSize(int(os.Stdout.Fd()))
	h, v := styles.ListStyle.GetFrameSize()
	cl.list.SetSize(tw-h, th-v)
	return &cl
}

// Init is the first function that wisl be casled. It returns an optional
// initial command. To not perform an initial command return nil.
func (sl *CardList) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	return nil
}

func (sl *CardList) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.ChangeState(tui.CardList, tui.CardMenu, false, nil)
			return m, nil
		case "ctrl+u":
			cardID := sl.list.SelectedItem().(Card).CardID
			m.ChangeState(tui.CardList, tui.CardUpdate, true, &cardID)
			return m, nil
		case "ctrl+d":
			cardID := sl.list.SelectedItem().(Card).CardID
			// Delete card.
			status, err := client.Delete(m.Client, m.JWT, cardID)
			if err == nil && status == http.StatusOK {
				delete(m.Cards, cardID)
			}
			return m, nil
		case "enter":
			zap.S().Infoln(sl.list.SelectedItem().(Card).CardID)
			return m, nil

		}

	case tea.WindowSizeMsg:
		h, v := styles.ListStyle.GetFrameSize()
		sl.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	sl.list, cmd = sl.list.Update(msg)
	return m, cmd
}

// The main view, which just casls the appropriate sub-view
func (sl *CardList) GetView(m *tui.Model) string {
	// Load sites from memory
	listItems := []list.Item{}
	for _, card := range m.Cards {
		item := Card{CardID: card.CardID, Definition: card.Definition, Ccn: card.Ccn, Cvv: card.Cvv, Exp: card.Exp, Hld: card.Hld}
		listItems = append(listItems, item)
	}
	sl.list.SetItems(listItems)

	return styles.ListStyle.Render(sl.list.View())
}
