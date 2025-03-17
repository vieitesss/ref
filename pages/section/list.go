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

func NewList(sections []scraper.Section) list.Model {
	var items []list.Item
	for _, s := range sections {
		items = append(items, item{name: s.Name})
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	l := list.New(items, delegate, HEIGHT, len(items) * 2 + 7)
	l.Title = "Sections"

	return l
}
