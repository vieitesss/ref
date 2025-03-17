package section

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/vieitesss/ref/scraper"
)

type Props struct {}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	textinput textinput.Model
	list list.Model
	sections []scraper.Section
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return c.textinput.Focus()
}

func New() *Component {
	sections := scraper.GetSections()

	c := &Component{
		list: NewList(sections),
		textinput: textinput.New(),
		sections: sections,
	}

	return c
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		c.list.SetSize(msg.Width-h, msg.Height-v)
	// case tea.KeyMsg:
	// 	if msg.Type == tea.KeyEnter {
	// 		// Lifted state power! Woohooo
	// 		reactea.SetCurrentRoute("/displayname")
	//
	// 		return nil
	// 	}
	}

	var ti_cmd tea.Cmd
	var list_cmd tea.Cmd
	c.textinput, ti_cmd = c.textinput.Update(msg)
	c.list, list_cmd = c.list.Update(msg)
	return tea.Batch(ti_cmd, list_cmd)
}

func (c *Component) Render(int, int) string {
	return docStyle.Render(c.list.View())
}
