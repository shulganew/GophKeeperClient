package gtext

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"golang.org/x/term"
)

// Implemet State.
var _ tui.State = (*GtextList)(nil)

type Gtext oapi.Gtext

func (g Gtext) Title() string {
	return fmt.Sprint(g.Definition)
}
func (g Gtext) Description() string {
	return fmt.Sprint(g.Note)
}
func (g Gtext) FilterValue() string { return g.Definition }

type GtextList struct {
	list list.Model
}

// List saved gtext credentials.
func NewGtextList() *GtextList {
	// Create empty list items
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "My secret notes. <ctrl+u> - update, <ctrl+d> - delete"
	cl := GtextList{list: l}
	// Fix terminal bag.
	tw, th, _ := term.GetSize(int(os.Stdout.Fd()))
	h, v := styles.ListStyle.GetFrameSize()
	cl.list.SetSize(tw-h, th-v)
	return &cl
}

// Init is the first function that wisl be casled. It returns an optional
// initial command. To not perform an initial command return nil.
func (sl *GtextList) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	return nil
}

func (sl *GtextList) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "esc":
			m.ChangeState(tui.GtextList, tui.GtextMenu, false, nil)
			return m, nil
		case "ctrl+u":
			gtextID := sl.list.SelectedItem().(Gtext).GtextID
			m.ChangeState(tui.GtextList, tui.GtextUpdate, true, &gtextID)
			return m, nil
		case "ctrl+d":
			if len(sl.list.Items()) == 0 {
				return m, nil
			}
			gtextID := sl.list.SelectedItem().(Gtext).GtextID
			// Delete site.
			status, err := client.DeleteAny(m.Client, m.JWT, gtextID)
			if err == nil && status == http.StatusOK {
				delete(m.Gtext, gtextID)
			}
			return m, nil
		case "enter":
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
func (sl *GtextList) GetView(m *tui.Model) string {
	// Load gtexts from memory
	listItems := []list.Item{}
	for _, text := range m.Gtext {
		item := Gtext{GtextID: text.GtextID, Definition: text.Definition, Note: getNoteSecond(&text.Note)}
		listItems = append(listItems, item)
	}
	sl.list.SetItems(listItems)

	return styles.ListStyle.Render(sl.list.View())
}

// Return secord row from list text.
func getNoteSecond(text *string) string {
	scanner := bufio.NewScanner(strings.NewReader(*text))
	// Skip first sentence.
	scanner.Scan()
	if scanner.Scan() {
		return scanner.Text()
	}

	return "No header note."
}
