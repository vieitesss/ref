package cheatsheet

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/vieitesss/ref/scraper"
	l "github.com/vieitesss/ref/components/list"
)

type Props struct {
	Section   string
	Reference string
}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	loading bool
	spinner spinner.Model
	list    list.Model
	titles  []string
}

type titleMsg []string

var (
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
	spinnerStyle = lipgloss.NewStyle().Margin(1, 0, 0, 1).Foreground(lipgloss.Color("205"))
)

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return getReferenceCheatsheetTitlesCmd(props.Reference)
}

func New() *Component {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return &Component{
		loading: true,
		spinner: s,
	}
}

func getReferenceCheatsheetTitlesCmd(ref string) tea.Cmd {
	return func() tea.Msg {
		titles := scraper.GetCheatSheet(ref)
		return titleMsg(titles)
	}
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case titleMsg:
		c.titles = []string(msg)
		c.list = l.NewCheatsheetList(c.Props().Reference, c.titles)
		c.loading = false

	case tea.WindowSizeMsg:
		l.UpdateSize(&c.list, msg.Width, msg.Height)

	case tea.KeyMsg:
		if key.Matches(msg, l.BackKey) {
			reactea.SetCurrentRoute(fmt.Sprintf("%s/references", c.Props().Section))

			return nil
		}
	}

	var cmd_list tea.Cmd
	var cmd_spinner tea.Cmd
	c.list, cmd_list = c.list.Update(msg)
	c.spinner, cmd_spinner = c.spinner.Update(msg)

	return tea.Batch(cmd_list, cmd_spinner)
}

func (c *Component) Render(int, int) string {
	if c.loading {
		return fmt.Sprintf("%v Getting %s cheatsheet", c.spinner.View(), c.Props().Section)
	}

	return docStyle.Render(c.list.View())
}
