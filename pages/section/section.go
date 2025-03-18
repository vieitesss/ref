package section

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/vieitesss/ref/scraper"
)

type Props struct{}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	list     list.Model
	sections []scraper.Section
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return nil
}

func New() *Component {
	sections := scraper.GetSections()

	c := &Component{
		list:     NewList(sections),
		sections: sections,
	}

	return c
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		UpdateSize(&c.list, msg.Width, msg.Height)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			selected := c.list.SelectedItem().FilterValue()
			reactea.SetCurrentRoute(fmt.Sprintf("%s/references", selected))

			return nil
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return cmd
}

func (c *Component) Render(int, int) string {
	return docStyle.Render(c.list.View())
}
