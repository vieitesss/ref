package scraper

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

	for c := range h.ChildNodes() {
		if c.Type == html.TextNode {
			text += c.Data
		}
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

func parseCodeInline(code *html.Node) string {
	return fmt.Sprintf("`%s`", getAllText(code))
}

func parseCode(lang string, code *html.Node) string {
	return fmt.Sprintf("```%s\n%s\n```", lang, getAllText(code))
}

func ParsePre(lang string, pre *html.Node) string {
	code := pre.LastChild
	text := fmt.Sprintf("%s\n", parseCode(lang, code))
	return text
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

	for c := range p.ChildNodes() {
		if c.Type == html.ElementNode && c.DataAtom == atom.A {
			text += ParseA(c)
		} else if c.Type == html.TextNode {
			text += c.Data
		}
	}

	return fmt.Sprintf("%s\n", text)
}

func getTHeadColumns(thead *html.Node) int {
	var tr *html.Node
	var columns int

	for c := range thead.ChildNodes() {
		if c.Type == html.ElementNode && c.DataAtom == atom.Tr {
			tr = c
			break
		}
	}

	for th := range tr.ChildNodes() {
		if th.Type == html.ElementNode && th.DataAtom == atom.Th {
			columns += 1
		}
	}

	return columns
}

func parseTd(lang string, td *html.Node) string {
	var text string

	for c := range td.ChildNodes() {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Code {
				text += parseCodeInline(c)
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
				text += fmt.Sprintf("%s|", getAllText(col))
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

	columns := getTHeadColumns(thead)

	text += parseTHeadTBody(lang, thead)
	text += fmt.Sprintf("%s|\n", strings.Repeat("|---", columns))
	text += parseTHeadTBody(lang, tbody)

	return text
}
