package scraper

import "fmt"

func ParseCode(lang, text string) string {
	return fmt.Sprintf("```%s\n%s\n```", lang, text)
}
