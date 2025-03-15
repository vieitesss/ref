package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

var c *colly.Collector

type Section struct {
	Name       string
	References []string
}

func NewSection(name string, references []string) Section {
	return Section{
		Name:       name,
		References: references,
	}
}

func Scrapper() *colly.Collector {
	if c != nil {
		return c
	}

	c = colly.NewCollector(colly.AllowedDomains("quickref.me"))

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	return c
}

func GetSections() []Section {
	var sections []Section
	Scrapper().OnHTML("h2.font-medium", func(e *colly.HTMLElement) {
		as := e.DOM.NextAllFiltered("div + div.grid").First().ChildrenFiltered("a").EachIter()

		var references []string
		for _, a := range as {
			val, _ := a.Attr("href")
			references = append(references, val)
		}

		sections = append(sections, NewSection(e.Text, references))
	})

	Scrapper().Visit("https://quickref.me")

	return sections
}

func GetSectionRef(section string) []string {
	var refTitles []string

	Scrapper().OnHTML(".h2-wrap", func(e *colly.HTMLElement) {
		title := e.DOM.ChildrenFiltered("h2").Text()
		refTitles = append(refTitles, title)
	})

	Scrapper().Visit(fmt.Sprintf("https://quickref.me%s", section))

	return refTitles
}

func main() {
	sections := GetSections()
	titles := GetSectionRef("/bash")
	fmt.Println(sections)
	fmt.Println(titles)
}
