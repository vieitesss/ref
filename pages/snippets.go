package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vieitesss/ref/scraper"
)

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	HalfUp: key.NewBinding(
		key.WithKeys("half up", "u"),
		key.WithHelp("u", "page up"),
	),
	HalfDown: key.NewBinding(
		key.WithKeys("half down", "d"),
		key.WithHelp("d", "page down"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace", "b"),
		key.WithHelp("󰭜/b", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type SnippetsProps struct {
	section, reference, title string
}

type SnippetsPageMsg SnippetsProps

// Received when the snippets are available
type snippetsMsg string

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	HalfUp   key.Binding
	HalfDown key.Binding
	Back     key.Binding
	Help     key.Binding
	Quit     key.Binding
}

type SnippetsPage struct {
	loading   bool
	spinner   spinner.Model
	keys      keyMap
	help      help.Model
	viewport  viewport.Model
	section   string
	reference string
	title     string
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.HalfUp, k.HalfDown, k.Back, k.Quit, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit, k.Help},
		{k.Up, k.Down, k.HalfUp, k.HalfDown},
	}
}

func NewSnippetsPage(p SnippetsProps) SnippetsPage {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return SnippetsPage{
		keys:      keys,
		loading:   true,
		spinner:   s,
		section:   p.section,
		reference: p.reference,
		title:     p.title,
		help:      help.New(),
	}
}

func (s SnippetsPage) Init() tea.Cmd {
	return tea.RequestWindowSize
}

func getTextCmd(reference, title string) tea.Cmd {
	return func() tea.Msg {
		in := scraper.GetSnippets(reference, title)
		return snippetsMsg(in)
	}
}

func setViewportSize(v *viewport.Model, w, h int) {
	v.SetWidth(w - 4)
	v.SetHeight(h - 4)
}

func (s SnippetsPage) Update(msg tea.Msg) (PageModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width = msg.Width
		height = msg.Height

		s.help.Width = width

		if s.loading {
			v := viewport.New()
			setViewportSize(&v, width, height)
			s.viewport = v
			return s, getTextCmd(s.reference, s.title)
		} else {
			setViewportSize(&s.viewport, width, height)
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Help):
			s.help.ShowAll = !s.help.ShowAll
		case key.Matches(msg, s.keys.Back):
			return s, func() tea.Msg {
				return CheatsheetPageMsg(CheatsheetProps{
					section:   s.section,
					reference: s.reference,
				})
			}
		}

	case snippetsMsg:
		s.loading = false
		r, _ := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dracula"),
		)
		out, _ := r.Render(string(msg))
		s.viewport.SetContent(out)
		return s, tea.RequestWindowSize
	}

	s.viewport, cmd = s.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s SnippetsPage) View() string {

	vp := lipgloss.NewStyle().
		Align(lipgloss.Top).
		Border(lipgloss.RoundedBorder()).
		BorderBottom(false).
		BorderTop(false).
		Render(s.viewport.View())

	var percentage string

	switch {
	case s.viewport.AtBottom():
		percentage = "At bottom"
	case s.viewport.AtTop():
		percentage = "At top"
	default:
		percentage = fmt.Sprintf("%.0f%%", s.viewport.ScrollPercent()*100)
	}

	h := lipgloss.NewStyle().
		Margin(0, 2).
		Align(lipgloss.Bottom).
		Render(s.help.View(s.keys))

	vpWidth := lipgloss.Width(vp)

	// Has always 2 due to the border
	if vpWidth < 3 {
		return lipgloss.JoinVertical(0, vp, h)
	}

	title := fmt.Sprintf("%s - %s", s.reference, s.title)
	topBar := lipgloss.NewStyle().
		Width(vpWidth).
		Render(fmt.Sprintf("%s %s %s%s",
			lipgloss.RoundedBorder().TopLeft,
			title,
			strings.Repeat(lipgloss.RoundedBorder().Top, vpWidth-4-len(title)),
			lipgloss.RoundedBorder().TopRight,
		))

	bottomBar := lipgloss.NewStyle().
		Width(vpWidth).
		Render(fmt.Sprintf("%s%s %s %s",
			lipgloss.RoundedBorder().BottomLeft,
			strings.Repeat(lipgloss.RoundedBorder().Top, vpWidth-4-len(percentage)),
			percentage,
			lipgloss.RoundedBorder().BottomRight,
		))

	return lipgloss.JoinVertical(0, topBar, vp, bottomBar, h)
}
