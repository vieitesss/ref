package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Section string
type Reference string

var c *colly.Collector

func Scrapper() *colly.Collector {
	if c != nil {
		return c
	}

	c = colly.NewCollector(colly.AllowedDomains("quickref.me"))
	c.AllowURLRevisit = true

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	return c
}

func GetSections() []Section {
	var sections []Section

	Scrapper().OnHTML("h2.font-medium", func(e *colly.HTMLElement) {
		sections = append(sections, Section(e.Text))
	})

	Scrapper().Visit("https://quickref.me")

	return sections
}

func GetSectionReferences(section string) []Reference {
	var references []Reference

	Scrapper().OnHTML("h2.font-medium", func(e *colly.HTMLElement) {
		if section != e.Text {
			return
		}

		as := e.DOM.NextAllFiltered("div + div.grid").First().ChildrenFiltered("a").EachIter()

		for _, a := range as {
			val, _ := a.Attr("href")
			references = append(references, Reference(strings.Replace(val, "/", "", 1)))
		}
	})

	Scrapper().Visit("https://quickref.me")

	return references
}

func GetReference(reference string) []string {
	var refTitles []string

	Scrapper().OnHTML(".h2-wrap", func(e *colly.HTMLElement) {
		title := e.DOM.ChildrenFiltered("h2").Text()
		refTitles = append(refTitles, title)
	})

	Scrapper().Visit(fmt.Sprintf("https://quickref.me/%s", reference))

	return refTitles
}
