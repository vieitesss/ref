package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/ref/pages"
)

func main() {
	program := tea.NewProgram(pages.NewMainPage(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
