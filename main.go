package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Section struct {
	Name string
	References []string
}

func NewSection(name string, references []string) Section {
	return Section{
		Name: name,
		References: references,
	}
}

var sections []Section

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("quickref.me"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnHTML("h2.font-medium", func(e *colly.HTMLElement) {
		as := e.DOM.NextAllFiltered("div + div.grid").First().ChildrenFiltered("a").EachIter()

		var references []string
		for _, a := range as {
			val, _ := a.Attr("href")
			references = append(references, val)
		}

		sections = append(sections, NewSection(e.Text, references))
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println(sections)
	})

	c.Visit("https://quickref.me")
}
