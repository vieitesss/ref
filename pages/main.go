package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss/v2"
	l "github.com/vieitesss/ref/components/list"
	"github.com/vieitesss/ref/scraper"
)

var (
	width, height = 0, 0

	spinnerStyle = lipgloss.NewStyle().AlignVertical(lipgloss.Top).Margin(1, 0, 0, 1).Foreground(lipgloss.Color("205"))
	docStyle     = lipgloss.NewStyle().AlignVertical(lipgloss.Top).Margin(1, 0, 0, 1)
)

type showSectionsMsg []scraper.Section
type showReferencesMsg []scraper.Reference
type showCheatsheetMsg []string
type showSnippetsMsg string

type MainPage struct {
	page      Page
	list      list.Model
	viewport  viewport.Model
	help      help.Model
	keys      keyMap
	spinner   spinner.Model
	loading   bool
	section   string
	reference string
	refPart   string
}

func NewMainPage() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return MainPage{
		loading: true,
		spinner: s,
	}
}

func (m MainPage) Init() tea.Cmd {
	return tea.Batch(
		m.loadPageCmd(Sections),
		m.spinner.Tick,
	)
}

func (m MainPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.page {
			case Sections:
				if selected := m.list.SelectedItem(); selected != nil {
					m.section = selected.FilterValue()
					return m, m.loadPageCmd(References)
				}
			case References:
				if selected := m.list.SelectedItem(); selected != nil {
					m.reference = selected.FilterValue()
					return m, m.loadPageCmd(Cheatsheet)
				}
			case Cheatsheet:
				if selected := m.list.SelectedItem(); selected != nil {
					m.refPart = selected.FilterValue()
					return m, m.loadPageCmd(Snippets)
				}
			}
		case "backspace":
			switch m.page {
			case References:
				return m, m.loadPageCmd(Sections)
			case Cheatsheet:
				return m, m.loadPageCmd(References)
			case Snippets:
				return m, m.loadPageCmd(Cheatsheet)
			}
		}

	case tea.WindowSizeMsg:
		width = msg.Width
		height = msg.Height
		if len(m.list.Items()) > 0 {
			l.UpdateSize(&m.list, height)
		}

		return m, nil

	case showSectionsMsg:
		m.loading = false
		m.list = l.NewSectionsList(m.title(), msg, width, height)

		return m, nil

	case showReferencesMsg:
		m.loading = false
		m.list = l.NewReferencesList(m.title(), msg, width, height)

		return m, nil

	case showCheatsheetMsg:
		m.loading = false
		m.list = l.NewCheatsheetList(m.title(), msg, width, height)

		return m, nil

	case showSnippetsMsg:
		m.loading = false

		v := viewport.New()
		setViewportSize(&v, width, height)
		m.viewport = v

		m.help = help.New()
		m.keys = keys

		r, _ := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dracula"),
		)
		out, _ := r.Render(string(msg))
		m.viewport.SetContent(out)

		return m, nil
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.page == Snippets {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
		m.help, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainPage) View() string {
	if !m.loading {
		if m.page == Snippets {
			return m.viewSnippets()
		}
		return lipgloss.NewStyle().
			MarginTop(1).
			Render(fmt.Sprintf("%v", m.list.View()))
	}

	return fmt.Sprintf("%s %s", m.spinner.View(), m.loadingText())
}

func (m *MainPage) loadPageCmd(p Page) tea.Cmd {
	m.loading = true
	m.page = p
	return tea.Batch(
		pageConfigs[p].Cmd(*m),
		m.spinner.Tick,
	)
}

func (m MainPage) title() string {
	return pageConfigs[m.page].Title(m)
}

func (m MainPage) loadingText() string {
	return pageConfigs[m.page].LoadingText(m)
}

func setViewportSize(v *viewport.Model, w, h int) {
	v.SetWidth(w - 4)
	v.SetHeight(h - 4)
}

func (m MainPage) viewSnippets() string {
	var sections []string
	vpWidth := m.viewport.Width() + 2

	title := pageConfigs[Snippets].Title(m)
	topBar := lipgloss.NewStyle().
		Width(vpWidth).
		Render(fmt.Sprintf("%s %s %s%s",
			lipgloss.RoundedBorder().TopLeft,
			title,
			strings.Repeat(lipgloss.RoundedBorder().Top, max(vpWidth-4-len(title), 0)),
			lipgloss.RoundedBorder().TopRight,
		))

	sections = append(sections, topBar)

	vp := lipgloss.NewStyle().
		Align(lipgloss.Top).
		Border(lipgloss.RoundedBorder()).
		BorderBottom(false).
		BorderTop(false).
		Render(m.viewport.View())

	sections = append(sections, vp)

	var percentage string

	switch {
	case m.viewport.AtBottom():
		percentage = "At bottom"
	case m.viewport.AtTop():
		percentage = "At top"
	default:
		percentage = fmt.Sprintf("%.0f%%", m.viewport.ScrollPercent()*100)
	}

	bottomBar := lipgloss.NewStyle().
		Width(vpWidth).
		Render(fmt.Sprintf("%s%s %s %s",
			lipgloss.RoundedBorder().BottomLeft,
			strings.Repeat(lipgloss.RoundedBorder().Top, max(vpWidth-4-len(percentage), 0)),
			percentage,
			lipgloss.RoundedBorder().BottomRight,
		))

	sections = append(sections, bottomBar)

	h := lipgloss.NewStyle().
		Margin(0, 2).
		Align(lipgloss.Bottom).
		Render(m.help.View(m.keys))

	sections = append(sections, h)

	return lipgloss.JoinVertical(0, sections...)
}
