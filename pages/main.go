package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	width, height = 0, 0

	spinnerStyle = lipgloss.NewStyle().AlignVertical(lipgloss.Top).Margin(1, 0, 0, 1).Foreground(lipgloss.Color("205"))
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
)

type MainPage struct {
	currentPage tea.Model
}

func NewMainPage() MainPage {
	return MainPage{
		currentPage: NewSectionPage(),
	}
}

func (m MainPage) Init() tea.Cmd {
	return m.currentPage.Init()
}

func (m MainPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	newPage := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width = msg.Width
		height = msg.Height
	case SectionPageMsg, ReferencesPageMsg, CheatsheetPageMsg:
		newPage = true
	}

	switch msg := msg.(type) {
	case SectionPageMsg:
		m.currentPage = NewSectionPage()
	case ReferencesPageMsg:
		m.currentPage = NewReferencesPage(string(msg))
	case CheatsheetPageMsg:
		m.currentPage = NewCheatsheetPage(CheatsheetProps(msg))
	}

	if !newPage {
		m.currentPage, cmd = m.currentPage.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		cmds = append(cmds, m.currentPage.Init())
	}

	return m, tea.Batch(cmds...)
}

func (m MainPage) View() string {
	return m.currentPage.View()
}
