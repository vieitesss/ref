package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vieitesss/ref/pkg/scraper"
)

var spinnerStyle = lipgloss.NewStyle().AlignVertical(lipgloss.Top).Margin(1, 0, 0, 1).Foreground(lipgloss.Color("205"))

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
	width     int
	height    int
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
	case tea.KeyPressMsg:
		if key.Matches(msg, m.keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
			m.updateViewportSize()

			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		cmd = m.handleKeyMsg(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		return m.handleResize(msg)

	case showSectionsMsg, showReferencesMsg, showCheatsheetMsg:
		m.loading = false
		m.setList(m.title(), msg, m.width, m.height)

		return m, nil

	case showSnippetsMsg:
		m.loading = false
		m.setViewport(msg)

		return m, nil
	}

	updates := m.updateComponents(msg)
	cmds = append(cmds, updates)

	return m, tea.Batch(cmds...)
}

func (m MainPage) View() string {
	return pageConfigs[m.page].View(m)
}

func (m *MainPage) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit
	case "enter":
		return m.handleEnter()
	case "backspace":
		return m.handleBackspace()
	}

	return nil
}

func (m *MainPage) handleEnter() tea.Cmd {
	selected := m.list.SelectedItem()

	if selected == nil {
		return nil
	}

	switch m.page {
	case Sections:
		m.section = selected.FilterValue()
		return m.loadPageCmd(References)
	case References:
		m.reference = selected.FilterValue()
		return m.loadPageCmd(Cheatsheet)
	case Cheatsheet:
		m.refPart = selected.FilterValue()
		return m.loadPageCmd(Snippets)
	}

	return nil
}

func (m *MainPage) handleBackspace() tea.Cmd {
	switch m.page {
	case References:
		return m.loadPageCmd(Sections)
	case Cheatsheet:
		return m.loadPageCmd(References)
	case Snippets:
		return m.loadPageCmd(Cheatsheet)
	}

	return nil
}

func (m MainPage) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	switch {
	case m.page == Snippets:
		m.updateViewportSize()
	case len(m.list.Items()) > 0:
		UpdateSize(&m.list, m.height)
	}

	return m, nil
}

func (m *MainPage) updateComponents(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch {
	case m.loading:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case m.page == Snippets:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
		m.help, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *MainPage) setList(title string, items tea.Msg, w, h int) {
	switch items := items.(type) {
	case showReferencesMsg:
		m.list = NewReferencesList(title, items, w, h)
	case showSectionsMsg:
		m.list = NewSectionsList(title, items, w, h)
	case showCheatsheetMsg:
		m.list = NewCheatsheetList(title, items, w, h)
	}
}

func (m *MainPage) setViewport(msg showSnippetsMsg) {
	m.viewport = viewport.New()
	m.viewport.FillHeight = false
	m.updateViewportSize()

	m.help = help.New()
	m.keys = keys

	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dracula"),
	)
	out, _ := r.Render(string(msg))
	m.viewport.SetContent(out)
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

func (m *MainPage) updateViewportSize() {
	// width = (terminal - borders) width
	m.viewport.SetWidth(m.width - 2)

	// height = (terminal - help - top and bottom bars) height
	actualHeight := m.height - lipgloss.Height(m.help.View(m.keys)) - 2
	m.viewport.SetHeight(actualHeight)
}

func (m MainPage) getScrollPercentage() string {
	var percentage string

	switch {
	case m.viewport.AtBottom():
		percentage = "At bottom"
	case m.viewport.AtTop():
		percentage = "At top"
	default:
		percentage = fmt.Sprintf("%.0f%%", m.viewport.ScrollPercent()*100)
	}

	return percentage
}
