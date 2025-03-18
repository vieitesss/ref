package section

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/vieitesss/ref/scraper"
)

type Props struct{}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	loading  bool
	spinner  spinner.Model
	list     list.Model
	sections []scraper.Section
}

type sectionsMsg []scraper.Section

var (
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
	spinnerStyle = lipgloss.NewStyle().Margin(1, 0, 0, 1).Foreground(lipgloss.Color("205"))
)

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return getSectionsCmd
}

func New() *Component {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	c := &Component{
		loading:  true,
		spinner:  s,
	}

	return c
}

func getSectionsCmd() tea.Msg {
	secs := scraper.GetSections()
	return sectionsMsg(secs)
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case sectionsMsg:
		c.sections = []scraper.Section(msg)
		c.list = NewList(c.sections)
		c.loading = false

	case tea.WindowSizeMsg:
		UpdateSize(&c.list, msg.Width, msg.Height)

	case tea.KeyMsg:
		if c.list.SettingFilter() {
			break
		}

		if msg.Type == tea.KeyEnter {
			selected := c.list.SelectedItem().FilterValue()
			reactea.SetCurrentRoute(fmt.Sprintf("%s/references", selected))
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
		return fmt.Sprintf("%v Getting sections", c.spinner.View())
	}

	return docStyle.Render(c.list.View())
}
