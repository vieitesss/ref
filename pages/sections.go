package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	l "github.com/vieitesss/ref/components/list"
	"github.com/vieitesss/ref/scraper"
)

// Received when changing to this page
type SectionPageMsg struct{}

// Received when sections data is available
type sectionsMsg []scraper.Section

type SectionPage struct {
	loading  bool
	spinner  spinner.Model
	list     list.Model
	sections []scraper.Section
}

func NewSectionPage() SectionPage {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	c := SectionPage{
		loading: true,
		spinner: s,
	}

	return c
}

func getSectionsCmd() tea.Msg {
	secs := scraper.GetSections()
	return sectionsMsg(secs)
}

func (s SectionPage) Init() tea.Cmd {
	return getSectionsCmd
}

func (s SectionPage) Update(msg tea.Msg) (PageModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case sectionsMsg:
		s.sections = []scraper.Section(msg)
		s.list = l.NewSectionsList("Sections", s.sections, width, height-2)
		s.loading = false

	case tea.WindowSizeMsg:
		if !s.loading {
			l.UpdateSize(&s.list, msg.Width, msg.Height-2)
		}

	case tea.KeyMsg:
		if s.list.SettingFilter() {
			break
		}

		if msg.Key().Code == tea.KeyEnter {
			selected := s.list.SelectedItem().FilterValue()
			return s, func() tea.Msg {
				return ReferencesPageMsg(selected)
			}
		}
	}

	if !s.loading {
		s.list, cmd = s.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	s.spinner, cmd = s.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (c SectionPage) View() string {
	if c.loading {
		content := fmt.Sprintf("%v Getting sections", c.spinner.View())
		return spinnerStyle.Height(height).Width(width).Render(content)
	}

	return docStyle.Width(width).Render(c.list.View())
}
