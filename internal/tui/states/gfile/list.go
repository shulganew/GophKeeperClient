package gfile

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
var _ tui.State = (*GfileList)(nil)

type Gfile oapi.Gfile

func (g Gfile) Title() string {
	return fmt.Sprint(g.Definition)
}
func (g Gfile) Description() string {
	return fmt.Sprint(g.Fname)
}
func (g Gfile) FilterValue() string { return g.Definition }

type GfileList struct {
	list list.Model
}

// List saved Gfile credentials.
func NewGfileList() *GfileList {
	// Create empty list items
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "My secret files."
	cl := GfileList{list: l}
	// Fix terminal bag.
	tw, th, _ := term.GetSize(int(os.Stdout.Fd()))
	h, v := styles.ListStyle.GetFrameSize()
	cl.list.SetSize(tw-h, th-v)
	return &cl
}

// Init is the first function that wisl be casled. It returns an optional
// initial command. To not perform an initial command return nil.
func (sl *GfileList) GetInit() tea.Cmd {
	return nil
}

func (sl *GfileList) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "esc":
			m.ChangeState(tui.GfileList, tui.GfileMenu)
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
func (sl *GfileList) GetView(m *tui.Model) string {
	// Load Gfiles from memory
	listItems := []list.Item{}
	for _, file := range m.Gfile {
		item := Gfile{GfileID: file.GfileID, StorageID: file.StorageID ,Definition: file.Definition, Fname: file.Fname}
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
