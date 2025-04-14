package scraper

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func GetAllText(h *html.Node) string {
	var sb strings.Builder

	for c := range h.Descendants() {
		if c.Type == html.TextNode {
			sb.Write([]byte(c.Data))
		}
	}

	text := sb.String()
	text = strings.Replace(text, "|", "\\|", -1)
	text = strings.Replace(text, "*", "\\*", -1)

	return text
}

func GetFirstLevelText(h *html.Node) string {
	var sb strings.Builder

	for c := range h.ChildNodes() {
		if c.Type == html.TextNode {
			sb.Write([]byte(c.Data))
		}
	}

	return sb.String()
}

func GetTHeadColumns(thead *html.Node) int {
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
			columns++
		}
	}

	return columns
}

func AfterParse(text string) string {
	reg := regexp.MustCompile(`(?s)(\n\s*){3,}`)

	return reg.ReplaceAllString(text, "\n\n")
}
