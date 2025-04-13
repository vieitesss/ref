package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/ref/scraper"
)

type Page int

type PageConfig struct {
	Title, LoadingText func(MainPage) string
	Cmd                func(MainPage) tea.Cmd
}

const (
	Sections Page = iota
	References
	Cheatsheet
	Snippets
)

var pageConfigs = map[Page]PageConfig{
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
		LoadingText: func(m MainPage) string {
			return "Loading sections"
		},
	},
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
		LoadingText: func(m MainPage) string {
			return fmt.Sprintf("Loading \"%s\" references", m.section)
		},
	},
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
		LoadingText: func(m MainPage) string {
			return fmt.Sprintf("Loading \"%s\" cheatsheet", m.reference)
		},
	},
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
		LoadingText: func(m MainPage) string {
			return fmt.Sprintf("Loading \"%s - %s\" snippets", m.reference, m.refPart)
		},
	},
}
