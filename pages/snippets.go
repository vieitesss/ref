package pages

import "github.com/charmbracelet/bubbles/v2/key"

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

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	HalfUp   key.Binding
	HalfDown key.Binding
	Back     key.Binding
	Help     key.Binding
	Quit     key.Binding
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
