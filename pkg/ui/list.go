package ui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/vieitesss/ref/pkg/scraper"
)

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

func newList(items []list.Item, title string, width, height int) list.Model {
	listSize := len(items)
	l := list.New(items, defaultDelegate(), width, min(getHeightForList(listSize), height))
	l.KeyMap.ShowFullHelp.SetEnabled(false)
	l.Title = title

	return l
}

func NewSectionsList(title string, elems []scraper.Section, width, height int) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: string(e)})
	}

	return newList(items, title, width, height)
}

func customBindings() []key.Binding {
	return []key.Binding{BackKey}
}

func NewReferencesList(title string, elems []scraper.Reference, width, height int) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: string(e)})
	}

	l := newList(items, title, width, height)
	l.AdditionalShortHelpKeys = customBindings

	return l
}

func NewCheatsheetList(title string, elems []string, width, height int) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: e})
	}

	l := newList(items, title, width, height)
	l.AdditionalShortHelpKeys = customBindings

	return l
}

func UpdateSize(l *list.Model, height int) {
	listSize := len(l.Items())
	h := min(getHeightForList(listSize), height)
	l.SetHeight(h)
}

func getHeightForList(listSize int) int {
	return listSize + 8
}
