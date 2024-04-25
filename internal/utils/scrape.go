package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tubes2-be-mccf/internal/models"

	"github.com/gocolly/colly/v2"
)

// Get Wikipedia URL from title.
//
// Returns the wikipedia URL of the article with the given title.
//
// Returns error if the title is invalid or the article is not found.
func GetWikipediaURLFromTitle(title string) (string, error) {
	// Validate title
	res, err := http.Get("https://en.wikipedia.org/wiki/" + title)

	// Http protocol error / too many redirects
	if err != nil {
		return "", err
	}

	// Article not found
	if res.StatusCode == 404 {
		return "", fmt.Errorf("start article not found")
	}

	return res.Request.URL.String(), nil
}

// Get all internal links from the URL
func GetAllInternalLinks(url string) []string {
	// Initialize result array
	links := make(map[string]bool)
	result := make([]string, 0)

	// Initialize colly collector
	cl := colly.NewCollector()

	// On HTML element a
	cl.OnHTML("div#bodyContent a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Check if the link is an internal link and correspond to a specific article
		if len(link) > 6 && link[:6] == "/wiki/" && !strings.HasPrefix(link, "/wiki/File:") {

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

// ArticleScript is a struct to parse the JSON data from the <script type=""application/ld+json""> tag
type ArticleScript struct {
	Type       string `json:"@type"`
	MainEntity string `json:"mainEntity"` // ID "https://www.wikidata.org/entity/Q_ID"
	Name       string `json:"name"`       // Title
	URL        string `json:"url"`        // URL
	Headline   string `json:"headline"`   // Description
	Image      string `json:"image"`      // Thumbnaik
}

// Get article data from the URL
func GetArticleData(url string) models.Article {
	// Initialize colly collector
	cl := colly.NewCollector()

	// Initialize article
	var article models.Article

	// Find the <script> tag with type "application/ld+json"
	cl.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		// Get the text inside the <script> tag
		data := e.Text

		// Parse to JSON
		var articleScript ArticleScript
		json.Unmarshal([]byte(data), &articleScript)

		// Get article ID from the URL
		articleID := strings.Split(articleScript.MainEntity, "http://www.wikidata.org/entity/")[1]

		// Set the article data
		article.ID = articleID
		article.Title = articleScript.Name
		article.Description = articleScript.Headline
		article.URL = articleScript.URL
		article.Thumbnail = articleScript.Image
	})

	// Visit the URL
	cl.Visit(url)

	return article
}
