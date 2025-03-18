package references

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/vieitesss/ref/scraper"
)

type Props struct {
	Section string
}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	loading    bool
	spinner    spinner.Model
	list       list.Model
	references []scraper.Reference
}

type referencesMsg []scraper.Reference

var (
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
	spinnerStyle = lipgloss.NewStyle().Margin(1, 0, 0, 1).Foreground(lipgloss.Color("205"))
)

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return getSectionReferencesCmd(props.Section)
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

func getSectionReferencesCmd(section string) tea.Cmd {
	return func() tea.Msg {
		refs := scraper.GetSectionReferences(section)
		return referencesMsg(refs)
	}
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case referencesMsg:
		c.references = []scraper.Reference(msg)
		c.list = NewList(c.Props().Section, c.references)
		c.loading = false
	case tea.WindowSizeMsg:
		UpdateSize(&c.list, msg.Width, msg.Height)
		// case tea.KeyMsg:
		// 	if msg.Type == tea.KeyEnter {
		// 		// Lifted state power! Woohooo
		// 		reactea.SetCurrentRoute("/displayname")
		//
		// 		return nil
		// 	}
	}

	var cmd_list tea.Cmd
	var cmd_spinner tea.Cmd
	c.list, cmd_list = c.list.Update(msg)
	c.spinner, cmd_spinner = c.spinner.Update(msg)

	return tea.Batch(cmd_list, cmd_spinner)
}

func (c *Component) Render(int, int) string {
	if c.loading {
		return fmt.Sprintf("%v Getting %s references", c.spinner.View(), c.Props().Section)
	}

	return docStyle.Render(c.list.View())
}
