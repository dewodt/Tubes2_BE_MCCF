package controllers

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

func getAllInternalLinks(url string) []string {
	// Initialize result array
	links := make(map[string]bool)
	result := make([]string, 0)

	// ext := []string{".jpg", ".jpeg", ".png", ".gif", ".svg"}

	// Initialize colly collector
	cl := colly.NewCollector()

	// On HTML element a
	cl.OnHTML("div#bodyContent a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// skip := false
		// for _, ext := range ext {
		// 	if strings.HasSuffix(link, ext) {
		// 		skip = true
		// 		break
		// 	}
		// }
		// Check if the link is an internal link and correspond to a specific article
		if len(link) > 6 && link[:6] == "/wiki/" && !strings.Contains(link, ":") {

			fullLink := "https://en.wikipedia.org" + link

			// If the link is not in the map, add it
			if _, ok := links[fullLink]; !ok {
				links[fullLink] = true
				result = append(result, fullLink)
			}
		}
	})

	// Visit the URL
	cl.Visit(url)

	return result
}

// still need to be corrected
func getThumbnail(url string) string {
	cl := colly.NewCollector()
	var thumbnail string
	var link string
	cl.OnHTML("class#mw-file-description a[href]", func(e *colly.HTMLElement) {
		thumbnail = e.Attr("href")
		link = "https://en.wikipedia.org" + thumbnail
	})
	cl.Visit(url)
	return link
}
