package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vieitesss/ref/pkg/scraper"
)

type Page int

type PageConfig struct {
	Cmd   func(MainPage) tea.Cmd
	Title func(MainPage) string
	View  func(MainPage) string
}

const (
	Sections Page = iota
	References
	Cheatsheet
	Snippets
)

var (
	listStyle = lipgloss.NewStyle().MarginTop(1)
)

var pageConfigs = map[Page]PageConfig{
	// Page that shows the different sections
	Sections: PageConfig{
		Title: func(m MainPage) string {
			return "Sections"
		},
		Cmd: func(m MainPage) tea.Cmd {
			return func() tea.Msg {
				secs, _ := scraper.GetSections()
				return showSectionsMsg(secs)
			}
		},
		View: func(m MainPage) string {
			if m.loading {
				return fmt.Sprintf("%s %s", m.spinner.View(), "Loading sections")
			}

			return listStyle.Render(m.list.View())
		},
	},

	// Page that shows the references available for a specific section
	References: PageConfig{
		Title: func(m MainPage) string {
			return m.section
		},
		Cmd: func(m MainPage) tea.Cmd {
			return func() tea.Msg {
				refs := scraper.GetSectionReferences(m.section)
				return showReferencesMsg(refs)
			}
		},
		View: func(m MainPage) string {
			if m.loading {
				return fmt.Sprintf("%s Loading \"%s\" references", m.spinner.View(), m.section)
			}

			return listStyle.Render(m.list.View())
		},
	},

	// Page that shows the posible snippets to watch by category
	Cheatsheet: PageConfig{
		Title: func(m MainPage) string {
			return m.reference
		},
		Cmd: func(m MainPage) tea.Cmd {
			return func() tea.Msg {
				cheat := scraper.GetCheatSheet(m.reference)
				return showCheatsheetMsg(cheat)
			}
		},
		View: func(m MainPage) string {
			if m.loading {
				return fmt.Sprintf("%s Loading \"%s\" cheatsheet", m.spinner.View(), m.reference)
			}

			return listStyle.Render(m.list.View())
		},
	},

	// Page where the actual reference cheatsheet is shown
	Snippets: PageConfig{
		Title: func(m MainPage) string {
			return fmt.Sprintf("%s - %s", m.reference, m.refPart)
		},
		Cmd: func(m MainPage) tea.Cmd {
			return func() tea.Msg {
				data := scraper.GetSnippets(m.reference, m.refPart)
				return showSnippetsMsg(data)
			}
		},
		View: func(m MainPage) string {
			if m.loading {
				return fmt.Sprintf("%s Loading \"%s - %s\" snippets", m.spinner.View(), m.reference, m.refPart)
			}

			var sections []string

			h := lipgloss.NewStyle().
				Margin(0, 2).
				Align(lipgloss.Bottom).
				Render(m.help.View(m.keys))

			title := fmt.Sprintf("%s - %s", m.reference, m.refPart)
			topBar := lipgloss.NewStyle().
				Width(m.width).
				Render(fmt.Sprintf("%s %s %s%s",
					lipgloss.RoundedBorder().TopLeft,
					title,
					strings.Repeat(lipgloss.RoundedBorder().Top, max(m.width-4-len(title), 0)),
					lipgloss.RoundedBorder().TopRight,
				))

			percentage := m.getScrollPercentage()
			bottomBar := lipgloss.NewStyle().
				Width(m.width).
				Render(fmt.Sprintf("%s%s %s %s",
					lipgloss.RoundedBorder().BottomLeft,
					strings.Repeat(lipgloss.RoundedBorder().Top, max(m.width-4-len(percentage), 0)),
					percentage,
					lipgloss.RoundedBorder().BottomRight,
				))

			vp := lipgloss.NewStyle().
				Align(lipgloss.Top).
				Border(lipgloss.RoundedBorder()).
				BorderBottom(false).
				BorderTop(false).
				Render(m.viewport.View())

			sections = append(sections, topBar, vp, bottomBar, h)

			return lipgloss.JoinVertical(0, sections...)
		},
	},
}
