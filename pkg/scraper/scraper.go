package scraper

import (
	"fmt"
	"strings"

	colly "github.com/gocolly/colly/v2"
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

func GetSections() ([]Section, error) {
	var sections []Section
	var err error

	Scrapper().OnHTML("h2.font-medium", func(e *colly.HTMLElement) {
		sections = append(sections, Section(e.Text))
	})

	Scrapper().OnError(func (r *colly.Response, e error) {
		err = e
	})

	Scrapper().Visit("https://quickref.me")

	return sections, err
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

func GetCheatSheet(reference string) []string {
	var cheatTitles []string

	Scrapper().OnHTML(".h2-wrap", func(e *colly.HTMLElement) {
		hash_title := e.DOM.ChildrenFiltered("h2").Text()
		title := strings.Replace(hash_title, "#", "", 1)
		cheatTitles = append(cheatTitles, title)
	})

	Scrapper().Visit(fmt.Sprintf("https://quickref.me/%s", reference))

	return cheatTitles
}

func GetSnippets(reference, title string) string {
	var sb strings.Builder
	var text string

	Scrapper().OnHTML("h2", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Text, title) {
			return
		}

		sb.Write([]byte(ParseH2(e.DOM.Get(0))))

		list := e.DOM.SiblingsFiltered(".h3-wrap-list").First()
		h3s := list.Find(".h3-wrap").EachIter()

		for _, h := range h3s {
			// Section title
			h3 := h.ChildrenFiltered("h3").Get(0)
			sb.Write([]byte(fmt.Sprintf("%s\n", ParseH3(h3))))

			// Section content
			sec := h.ChildrenFiltered("div.section").First()
			for _, c := range sec.Children().EachIter() {
				node := c.Get(0)
				switch node.Data {
				case "pre":
					sb.Write([]byte(ParsePre(reference, node)))
				case "h4":
					sb.Write([]byte(ParseH4(node)))
				case "p":
					sb.Write([]byte(ParseP(node)))
				case "table":
					sb.Write([]byte(ParseTable(reference, node)))
				case "ul":
					sb.Write([]byte(ParseUl(reference, node)))
				}

				sb.Write([]byte("\n"))
			}
		}
	})

	Scrapper().Visit(fmt.Sprintf("https://quickref.me/%s", reference))

	text = AfterParse(sb.String())

	return text
}
