package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

var references []string

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("quickref.me"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnHTML("h2.font-medium", func(e *colly.HTMLElement) {
		as := e.DOM.NextAllFiltered("div + div.grid").First().ChildrenFiltered("a").EachIter()

		for _, a := range as {
			val, _ := a.Attr("href")
			references = append(references, val)
		}
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println(references)
	})

	c.Visit("https://quickref.me")
}
