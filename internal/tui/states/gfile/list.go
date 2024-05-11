package gfile

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
func NewGfileList(path string) *GfileList {
	// Create empty list items
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = fmt.Sprintf("<Insert> - download selected file to folder:  %s. <ctrl+d> - del", path)
	cl := GfileList{list: l}
	// Fix terminal bag.
	tw, th, _ := term.GetSize(int(os.Stdout.Fd()))
	h, v := styles.ListStyle.GetFrameSize()
	cl.list.SetSize(tw-h, th-v)
	return &cl
}

// Init is the first function that wisl be casled. It returns an optional
// initial command. To not perform an initial command return nil.
func (sl *GfileList) GetInit(m *tui.Model, updateID *string) tea.Cmd {
	return nil
}

func (sl *GfileList) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "esc":
			sl.list.ResetSelected()
			m.ChangeState(tui.GfileList, tui.GfileMenu, false, nil)
			return m, nil
		case "enter":
			return m, nil
		case "ctrl+d":
			fileID := sl.list.SelectedItem().(Gfile).GfileID
			// Delete file.
			status, err := client.DeleteFile(m.Client, m.JWT, fileID)
			if err == nil && status == http.StatusOK {
				delete(m.Gfile, fileID)
			}
		case "insert":
			zap.S().Infoln(client.GfileGet(m.Client, m.Conf, m.JWT, sl.list.SelectedItem().(Gfile).GfileID, sl.list.SelectedItem().(Gfile).Fname))
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

// The main view, which just casls the appropriate sub-view.
func (sl *GfileList) GetView(m *tui.Model) string {
	// Load Gfiles from memory
	listItems := []list.Item{}
	for _, file := range m.Gfile {
		item := Gfile{GfileID: file.GfileID, Definition: file.Definition, Fname: file.Fname}
		listItems = append(listItems, item)
	}
	sl.list.SetItems(listItems)
	return styles.ListStyle.Render(sl.list.View())
}
