package scraper

import (
	"fmt"
	"regexp"
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
	var sb strings.Builder

	for c := range span.ChildNodes() {
		switch {
		case c.Type == html.TextNode:
			sb.Write([]byte(c.Data))
		case c.Type == html.ElementNode:
			sb.Write([]byte(parseSpan(c)))
		}
	}

	return sb.String()
}

func parseCode(lang string, code *html.Node) string {
	text := GetAllText(code)

	if match, _ := regexp.Match(`^\s*\n*$`, []byte(text)); match {
		return "\n"
	}

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
	var sb strings.Builder

	for c := range p.ChildNodes() {
		switch {
		case c.Type == html.TextNode:
			sb.Write([]byte(c.Data))

		case c.Type == html.ElementNode:
			switch c.DataAtom {
			case atom.A:
				sb.Write([]byte(ParseA(c)))

			case atom.Strong:
				sb.Write([]byte(parseStrong(c)))
			}
		}
	}

	return fmt.Sprintf("%s\n", sb.String())
}

func parseTd(lang string, td *html.Node) string {
	var sb strings.Builder

	for c := range td.ChildNodes() {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Code {
				sb.Write([]byte(parseCodeInline(c, true)))
			} else if c.Data == "yel" {
				sb.Write([]byte(parseYel(c)))
			}
		} else if c.Type == html.TextNode {
			sb.Write([]byte(c.Data))
		}
	}

	return sb.String()
}

func parseTr(lang string, tr *html.Node) string {
	var sb strings.Builder

	sb.Write([]byte("|"))

	for col := range tr.ChildNodes() {
		if col.Type == html.ElementNode {
			if col.DataAtom == atom.Th {
				sb.Write([]byte(fmt.Sprintf("%s|", GetAllText(col))))
			} else if col.DataAtom == atom.Td {
				sb.Write([]byte(fmt.Sprintf("%s|", parseTd(lang, col))))
			}
		}
	}

	sb.Write([]byte("\n"))

	return sb.String()
}

func parseTHeadTBody(lang string, thb *html.Node) string {
	var sb strings.Builder

	for c := range thb.ChildNodes() {
		if c.Type == html.ElementNode && c.DataAtom == atom.Tr {
			sb.Write([]byte(parseTr(lang, c)))
		}
	}

	return sb.String()
}

func ParseTable(lang string, table *html.Node) string {
	var sb strings.Builder
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

	sb.Write([]byte(parseTHeadTBody(lang, thead)))
	sb.Write([]byte(fmt.Sprintf("%s|\n", strings.Repeat("|---", columns))))
	sb.Write([]byte(parseTHeadTBody(lang, tbody)))

	return sb.String()
}

func parseLi(lang string, li *html.Node) string {
	var sb strings.Builder

	for c := range li.ChildNodes() {
		switch {
		case c.Type == html.TextNode:
			sb.Write([]byte(c.Data))

		case c.Type == html.ElementNode:
			switch c.DataAtom {
			case atom.Strong:
				sb.Write([]byte(fmt.Sprintf("%s\n", parseStrong(c))))

			case atom.Pre:
				sb.Write([]byte(ParsePre(lang, c)))

			case atom.P:
				sb.Write([]byte(ParseP(c)))
			}
		}
	}

	return fmt.Sprintf("- %s\n", sb.String())
}

func ParseUl(lang string, ul *html.Node) string {
	var sb strings.Builder

	for c := range ul.ChildNodes() {
		switch {
		case c.Type == html.ElementNode && c.DataAtom == atom.Li:
			sb.Write([]byte(parseLi(lang, c)))
		}
	}

	return sb.String()
}
