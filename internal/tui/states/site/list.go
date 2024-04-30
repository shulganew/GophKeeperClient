package site

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/styles"
	"go.uber.org/zap"
	"golang.org/x/term"
)

// Implemet State.
var _ tui.State = (*SiteList)(nil)

type Site oapi.Site

func (s Site) Title() string {
	return fmt.Sprintf("%s ◉ %s", s.Definition, s.Site)
}
func (s Site) Description() string {
	return fmt.Sprintf("%s ◉ %s", s.Slogin, s.Spw)
}
func (s Site) FilterValue() string { return s.Site }

type SiteList struct {
	list list.Model
}

// SiteList, state 7
// List saved site credentials.
func NewSiteList() *SiteList {
	// Create empty list items
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Sites login and passowrds."
	sl := SiteList{list: l}
	// Fix terminal bag.
	tw, th, _ := term.GetSize(int(os.Stdout.Fd()))
	h, v := styles.ListStyle.GetFrameSize()
	sl.list.SetSize(tw-h, th-v)
	return &sl
}

// Init is the first function that wisl be casled. It returns an optional
// initial command. To not perform an initial command return nil.
func (sl *SiteList) GetInit() tea.Cmd {
	return nil
}

func (sl *SiteList) GetUpdate(m *tui.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.ChangeState(tui.SiteList, tui.MainMenu)
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
func (sl *SiteList) GetView(m *tui.Model) string {
	// Load sites from memory
	listItems := []list.Item{}
	for _, site := range m.Sites {
		item := Site{SiteID: site.SiteID, Definition: site.Definition, Site: site.Site, Slogin: site.Slogin, Spw: site.Spw}
		listItems = append(listItems, item)
	}
	sl.list.SetItems(listItems)

	return styles.ListStyle.Render(sl.list.View())
}
