package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	l "github.com/vieitesss/ref/components/list"
	"github.com/vieitesss/ref/scraper"
)

type CheatsheetProps struct {
	section, reference string
}

type CheatsheetPageMsg CheatsheetProps
type titleMsg []string

type CheatsheetPage struct {
	loading            bool
	section, reference string
	spinner            spinner.Model
	list               list.Model
	titles             []string
}

func NewCheatsheetPage(p CheatsheetProps) CheatsheetPage {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return CheatsheetPage{
		section:   p.section,
		reference: p.reference,
		loading:   true,
		spinner:   s,
	}
}

func (c CheatsheetPage) Init() tea.Cmd {
	return func() tea.Msg {
		titles := scraper.GetCheatSheet(c.reference)
		return titleMsg(titles)
	}
}

func (c CheatsheetPage) Update(msg tea.Msg) (PageModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case titleMsg:
		c.titles = []string(msg)
		c.list = l.NewCheatsheetList(c.reference, c.titles, width, height-2)
		c.loading = false

	case tea.WindowSizeMsg:
		if !c.loading {
			l.UpdateSize(&c.list, width, height-2)
		}

	case tea.KeyMsg:
		if c.list.SettingFilter() {
			break
		}

		if key.Matches(msg, l.BackKey) {
			return c, func() tea.Msg {
				return ReferencesPageMsg(c.section)
			}
		}

		if msg.Key().Code == tea.KeyEnter {
			selected := c.list.SelectedItem().FilterValue()
			return c, func() tea.Msg {
				return SnippetsPageMsg(SnippetsProps{
					section:   c.section,
					reference: c.reference,
					title:     selected,
				})
			}
		}
	}

	if !c.loading {
		c.list, cmd = c.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	c.spinner, cmd = c.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c CheatsheetPage) View() string {
	if c.loading {
		content := fmt.Sprintf("%v Getting %s cheatsheet", c.spinner.View(), c.reference)
		return spinnerStyle.Height(height).Width(width).Render(content)
	}

	return docStyle.Width(width).Render(c.list.View())
}
