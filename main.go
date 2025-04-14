package main

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/ref/pkg/ui"
)

func main() {
	program := tea.NewProgram(
		ui.NewMainPage(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
