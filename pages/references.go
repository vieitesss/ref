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

// Received when changing to this page
type ReferencesPageMsg string

// Received when references data is available
type referencesMsg []scraper.Reference

type ReferencesPage struct {
	loading    bool
	spinner    spinner.Model
	section    string
	list       list.Model
	references []scraper.Reference
}

func (r ReferencesPage) Init() tea.Cmd {
	return func() tea.Msg {
		refs := scraper.GetSectionReferences(r.section)
		return referencesMsg(refs)
	}
}

func NewReferencesPage(section string) ReferencesPage {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return ReferencesPage{
		loading: true,
		spinner: s,
		section: section,
	}
}

func (r ReferencesPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case referencesMsg:
		r.references = []scraper.Reference(msg)
		r.list = l.NewReferencesList(r.section, r.references, width, height-2)
		r.loading = false

	case tea.WindowSizeMsg:
		if !r.loading {
			l.UpdateSize(&r.list, msg.Width, msg.Height-2)
		}

	case tea.KeyMsg:
		if r.list.SettingFilter() {
			break
		}

		if key.Matches(msg, l.BackKey) {
			return r, func() tea.Msg {
				return SectionPageMsg{}
			}
		}

		if msg.Type == tea.KeyEnter {
			selected := r.list.SelectedItem().FilterValue()
			return r, func() tea.Msg {
				return CheatsheetPageMsg(CheatsheetProps{
					section:   r.section,
					reference: selected,
				})
			}
		}
	}

	if !r.loading {
		r.list, cmd = r.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	r.spinner, cmd = r.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
}

func (r ReferencesPage) View() string {
	if r.loading {
		content := fmt.Sprintf("%v Getting %s references", r.spinner.View(), r.section)
		return spinnerStyle.Height(height).Width(width).Render(content)
	}

	return docStyle.Width(width).Render(r.list.View())
}
