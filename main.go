package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("quickref.me"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		fmt.Printf("h2: %v\n", e.Text)
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("scraped!")
	})

	c.Visit("https://quickref.me")
}
