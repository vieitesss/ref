package cheatsheet

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

const WIDTH = 20
const HEIGHT = 50

type item struct {
	name string
}

var BackKey = key.NewBinding(
	key.WithKeys("backspace", "b"),
	key.WithHelp("󰭜/b", "back"),
)

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

func customBindings() []key.Binding {
	return []key.Binding{BackKey}
}

func NewList(ref string, elems []string) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: string(e)})
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	l := list.New(items, delegate, WIDTH, min(len(items)+7, HEIGHT))
	l.Title = fmt.Sprintf("%s", ref)
	l.AdditionalShortHelpKeys = customBindings
	l.KeyMap.ShowFullHelp.SetEnabled(false)

	return l
}

func UpdateSize(l *list.Model, width, height int) {
	l.SetSize(width, min(len(l.Items())+7, height))
}
