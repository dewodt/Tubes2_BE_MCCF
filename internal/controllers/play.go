package controllers

import (
	"fmt"
	"net/http"
	"tubes2-be-mccf/internal/models"
	"tubes2-be-mccf/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

// Result Request Data Structure
type PlayRequest struct {
	Algorithm string `form:"algorithm" binding:"required,oneof='IDS' 'BFS'"` // Algorithm type (IDS or BFS)
	Start     string `form:"start" binding:"required"`                       // Start wikipedia article title
	Target    string `form:"target" binding:"required"`                      // Target wikipedia article title
}

// Result Response Data Structure
type PlaySuccessResponse struct {
	TotalTraversed     int              `json:"totalTraversed"`     // Total article traversed
	ShortestPathLength int              `json:"shortestPathLength"` // Shortest path length
	Duration           float32          `json:"duration"`           // Duration of the search
	Articles           []models.Article `json:"articles"`           // Articles data used in paths: Path[]
	Paths              []models.Path    `json:"paths"`              // List of the shortest paths from start to target found
}

type FieldError struct {
	Field   string `json:"field"`   // Form field that caused the error
	Message string `json:"message"` // Error message
}

// Error Response Data Structure
type PlayErrorResponse struct {
	Error       string       `json:"error"`       // Error message
	Message     string       `json:"message"`     // Error message
	ErrorFields []FieldError `json:"errorFields"` // List of fields that caused the error
}

// Get Wikipedia URL from title.
//
// Returns the wikipedia URL of the article with the given title.
//
// Returns error if the title is invalid or the article is not found.
func getWikipediaURLFromTitle(title string) (string, error) {
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

// Get all internal links from a wikipedia URL.
//
// Returns a list of internal links found in the given wikipedia URL (english only en.wikipedia.org).
func getAllInternalLinks(url string) []string {
	// Initialize result array
	var links []string

	// Initialize colly collector
	cl := colly.NewCollector()

	// On HTML element a
	cl.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Check if the link is an internal link and correspond to a specific article
		if len(link) > 6 && link[:6] == "/wiki/" {
			links = append(links, "https://en.wikipedia.org"+link)
		}
	})

	// Visit the URL
	cl.Visit(url)

	return links
}

func SolveIDS(startURL string, targetURL string) (PlaySuccessResponse, error) {
	fmt.Println("Solving with IDS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)

	links := getAllInternalLinks(startURL)
	utils.PrintArrayString(links)

	// Placeholder
	return PlaySuccessResponse{}, nil
}

func solveBFS(startURL string, targetURL string) (PlaySuccessResponse, error) {
	fmt.Println("Solving with BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)

	links := getAllInternalLinks(startURL)
	utils.PrintArrayString(links)

	// Placeholder
	return PlaySuccessResponse{}, nil
}

func Solve(algorithm string, startURL string, targetURL string) (PlaySuccessResponse, error) {
	if algorithm == "IDS" {
		return SolveIDS(startURL, targetURL)
	} else {
		return solveBFS(startURL, targetURL)
	}
}

func PlayHandler(c *gin.Context) {
	// Validate request data
	var reqJSON PlayRequest
	err := c.ShouldBind(&reqJSON)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	// Get data
	algorithm := reqJSON.Algorithm
	startTitle := reqJSON.Start
	targetTitle := reqJSON.Target

	// Get start wikipedia URL (and validate title)
	startURL, err := getWikipediaURLFromTitle(startTitle)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request", "message": "Wikipedia start article not found", "errorFields": []FieldError{{"start", "Wikipedia start article not found"}}})
		return
	}
	// Get target wikipedia URL (and validate title)
	targetURL, err := getWikipediaURLFromTitle(targetTitle)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request", "message": "Wikipedia target article not found", "errorFields": []FieldError{{"target", "Wikipedia target article not found"}}})
		return
	}

	// Solve
	result, err := Solve(algorithm, startURL, targetURL)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the result
	c.JSON(200, result)
}
