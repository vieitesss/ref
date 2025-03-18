package section

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/vieitesss/ref/scraper"
)

const HEIGHT = 20

type item struct {
	name string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

func NewList(elems []scraper.Section) list.Model {
	var items []list.Item
	for _, e := range elems {
		items = append(items, item{name: string(e)})
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	l := list.New(items, delegate, HEIGHT, len(items) + 7)
	l.KeyMap.ShowFullHelp.SetEnabled(false)
	l.KeyMap.ShowFullHelp.Unbind()
	l.Title = "Sections"

	return l
}

func UpdateSize(l *list.Model, width, height int) {
	l.SetSize(width, min(len(l.Items()) + 7, height))
}
