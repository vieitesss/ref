package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/ref/app"

	"github.com/londek/reactea"
)

func main() {
	program := reactea.NewProgram(app.New(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
