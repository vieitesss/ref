package scraper

import (
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	// "golang.org/x/net/html/atom"
)

func getAllText(h *html.Node) string {
	var text string

	for c := range h.Descendants() {
		if c.Type == html.TextNode {
			text += c.Data
		}
	}

	return text
}

func getFirstLevelText(h *html.Node) string {
	var text string

	child := h.FirstChild

	for {
		if child.Type == html.TextNode {
			text += child.Data
		}

		if child.NextSibling == nil {
			break
		}

		child = child.NextSibling
	}

	return text
}

func ParseA(a *html.Node) string {
	var link string

	for _, r := range a.Attr {
		if r.Key == "href" {
			link = r.Val
			break
		}
	}

	text := getAllText(a)

	return fmt.Sprintf("[%s](%s)", text, link)
}

func ParsePre(lang string, pre *html.Node) string {
	code := pre.LastChild
	text := getAllText(code)

	return fmt.Sprintf("```%s\n%s\n```\n", lang, text)
}


func ParseH1(h1 *html.Node) string {
	parsed := getFirstLevelText(h1)
	return fmt.Sprintf("# %s\n", parsed)
}

func ParseH2(h2 *html.Node) string {
	parsed := getFirstLevelText(h2)
	return fmt.Sprintf("## %s\n", parsed)
}

func ParseH3(h3 *html.Node) string {
	parsed := getFirstLevelText(h3)
	return fmt.Sprintf("### %s\n", parsed)
}

func ParseH4(h4 *html.Node) string {
	parsed := getFirstLevelText(h4)
	return fmt.Sprintf("#### %s\n", parsed)
}

func ParseP(p *html.Node) string {
	var text string

	c := p.FirstChild
	for {
		if c.Type == html.ElementNode && c.DataAtom == atom.A {
			text += ParseA(c)
		} else if c.Type == html.TextNode {
			text += c.Data
		}

		if c.NextSibling == nil {
			break
		}

		c = c.NextSibling
	}

	return fmt.Sprintf("%s\n", text)
}
