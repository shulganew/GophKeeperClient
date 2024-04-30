package gtext

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
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
	l.Title = "My secret notes."
	cl := GtextList{list: l}
	// Fix terminal bag.
	tw, th, _ := term.GetSize(int(os.Stdout.Fd()))
	h, v := styles.ListStyle.GetFrameSize()
	cl.list.SetSize(tw-h, th-v)
	return &cl
}

// Init is the first function that wisl be casled. It returns an optional
// initial command. To not perform an initial command return nil.
func (sl *GtextList) GetInit() tea.Cmd {
	return nil
}

func (sl *GtextList) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "esc":
			m.ChangeState(tui.GtextList, tui.GtextMenu)
			return m, nil
		case "enter":
			zap.S().Infoln(sl.list.SelectedItem())
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
