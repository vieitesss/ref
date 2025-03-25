package scraper

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Inline elements
func ParseA(a *html.Node) string {
	var link string

	for _, r := range a.Attr {
		if r.Key == "href" {
			link = r.Val
			break
		}
	}

	text := GetAllText(a)

	return fmt.Sprintf("[%s](%s)", text, link)
}

func parseYel(yel *html.Node) string {
	return fmt.Sprintf("**%s**", GetAllText(yel))
}

func parseStrong(strong *html.Node) string {
	text := parseYel(strong)

	return text
}

func parseCodeInline(code *html.Node, inTable bool) string {
	return fmt.Sprintf("`%s`", GetAllText(code))
}

func parseSpan(span *html.Node) string {
	var text string

	for c := range span.ChildNodes() {
		switch {
		case c.Type == html.TextNode:
			text += c.Data
		case c.Type == html.ElementNode:
			text += parseSpan(c)
		}
	}

	return text
}

func parseCode(lang string, code *html.Node) string {
	return fmt.Sprintf("```%s\n%s\n```", lang, GetAllText(code))
}

// Block elements
func ParsePre(lang string, pre *html.Node) string {
	code := pre.LastChild
	text := fmt.Sprintf("%s\n", parseCode(lang, code))
	return text
}

func ParseH1(h1 *html.Node) string {
	parsed := GetFirstLevelText(h1)
	return fmt.Sprintf("# %s\n", parsed)
}

func ParseH2(h2 *html.Node) string {
	parsed := GetFirstLevelText(h2)
	return fmt.Sprintf("## %s\n", parsed)
}

func ParseH3(h3 *html.Node) string {
	parsed := GetFirstLevelText(h3)
	return fmt.Sprintf("### %s\n", parsed)
}

func ParseH4(h4 *html.Node) string {
	parsed := GetFirstLevelText(h4)
	return fmt.Sprintf("#### %s\n", parsed)
}

func ParseP(p *html.Node) string {
	var text string

	for c := range p.ChildNodes() {
		switch {
		case c.Type == html.TextNode:
			text += c.Data

		case c.Type == html.ElementNode:
			switch c.DataAtom {
			case atom.A:
				text += ParseA(c)

			case atom.Strong:
				text += parseStrong(c)
			}
		}
	}

	return fmt.Sprintf("%s\n", text)
}

func parseTd(lang string, td *html.Node) string {
	var text string

	for c := range td.ChildNodes() {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Code {
				text += parseCodeInline(c, true)
			} else if c.Data == "yel" {
				text += parseYel(c)
			}
		} else if c.Type == html.TextNode {
			text += c.Data
		}
	}

	return text
}

func parseTr(lang string, tr *html.Node) string {
	text := "|"

	for col := range tr.ChildNodes() {
		if col.Type == html.ElementNode {
			if col.DataAtom == atom.Th {
				text += fmt.Sprintf("%s|", GetAllText(col))
			} else if col.DataAtom == atom.Td {
				text += fmt.Sprintf("%s|", parseTd(lang, col))
			}
		}
	}

	text += "\n"

	return text
}

func parseTHeadTBody(lang string, thb *html.Node) string {
	var text string

	for c := range thb.ChildNodes() {
		if c.Type == html.ElementNode && c.DataAtom == atom.Tr {
			text += parseTr(lang, c)
		}
	}

	return text
}

func ParseTable(lang string, table *html.Node) string {
	var text string
	var thead, tbody *html.Node

	for c := range table.ChildNodes() {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Thead {
				thead = c
			} else if c.DataAtom == atom.Tbody {
				tbody = c
			}
		}
	}

	columns := GetTHeadColumns(thead)

	text += parseTHeadTBody(lang, thead)
	text += fmt.Sprintf("%s|\n", strings.Repeat("|---", columns))
	text += parseTHeadTBody(lang, tbody)

	return text
}
