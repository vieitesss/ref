package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/vieitesss/ref/app"
)

func main() {
	program := reactea.NewProgram(app.New(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
