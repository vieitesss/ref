package pages

import (
	_ "embed"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/vieitesss/ref/scraper"
)

//go:embed file.md
var in string

type SnippetsProps struct {
	section, reference, title string
}

type SnippetsPageMsg SnippetsProps

// Received when the snippets are available
type snippetsMsg string

type SnippetsPage struct {
	loading                   bool
	spinner                   spinner.Model
	viewport                  viewport.Model
	section, reference, title string
}

func NewSnippetsPage(p SnippetsProps) SnippetsPage {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return SnippetsPage{
		loading:   true,
		spinner:   s,
		section:   p.section,
		reference: p.reference,
		title:     p.title,
	}
}

func (s SnippetsPage) Init() tea.Cmd {
	return tea.WindowSize()
}

func getText(reference, title string) tea.Cmd {
	return func() tea.Msg {
		in := scraper.GetSnippets(reference, title)
		return snippetsMsg(in)
	}
}

func (s SnippetsPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width = msg.Width
		height = msg.Height

		if s.loading {
			v := viewport.New(width, height)
			s.viewport = v
			return s, getText(s.reference, s.title)
		} else {
			s.viewport.Height = height
			s.viewport.Width = width
		}

	case snippetsMsg:
		s.loading = false
		out, _ := glamour.Render(string(msg), "dark")
		s.viewport.SetContent(out)
		return s, nil
	}

	s.viewport, cmd = s.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s SnippetsPage) View() string {
	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Render(s.viewport.View())
}
