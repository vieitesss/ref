package references

import (
	"github.com/charmbracelet/bubbles/list"
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

	list       list.Model
	references []scraper.Reference
}

type referencesMsg []scraper.Reference

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return getSectionReferencesCmd(props.Section)
}

func New() *Component {
	return &Component{}
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

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return cmd
}

func (c *Component) Render(int, int) string {
	return docStyle.Render(c.list.View())
}
