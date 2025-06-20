package parserTools

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
)

func LogValue(value string, exists ...bool) string {
	if len(exists) > 0 && !exists[0] {
		return "empty"
	}
	if value == "" {
		return "empty"
	}
	return value
}

func PrintElementHTML(s *goquery.Selection) {
	html, err := s.Html()
	if err != nil {
		log.Printf("Error getting HTML: %v", err)
		return
	}
	log.Println("Element HTML:")
	log.Println(html)
}

func SaveHTMLToFile(doc *goquery.Document, filename string) error {
	html, err := doc.Html()
	if err != nil {
		return err
	}

	// Format the HTML for better readability
	formatted := strings.Replace(html, "><", ">\n<", -1)

	return os.WriteFile(filename, []byte(formatted), 0644)
}
