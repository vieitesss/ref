package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/vieitesss/ref/scraper"
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

var routes []string
var current_route string

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

func defaultDelegate() list.DefaultDelegate {
	del := list.NewDefaultDelegate()
	del.ShowDescription = false
	del.SetSpacing(0)

	return del
}

func newList(items []list.Item, title string) list.Model {
	l := list.New(items, defaultDelegate(), WIDTH, min(len(items)+7, HEIGHT))
	l.KeyMap.ShowFullHelp.SetEnabled(false)
	l.Title = title

	return l
}

func NewSectionsList(title string, elems []scraper.Section) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: string(e)})
	}

	current_route = "default"

	return newList(items, title)
}

func customBindings() []key.Binding {
	return []key.Binding{BackKey}
}

func NewReferencesList(title string, elems []scraper.Reference) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: string(e)})
	}

	l := newList(items, title)
	l.AdditionalShortHelpKeys = customBindings

	return l
}

func NewCheatsheetList(title string, elems []string) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: e})
	}

	l := newList(items, title)
	l.AdditionalShortHelpKeys = customBindings

	return l
}

func UpdateSize(l *list.Model, width, height int) {
	l.SetSize(width, min(len(l.Items())+7, height))
}
