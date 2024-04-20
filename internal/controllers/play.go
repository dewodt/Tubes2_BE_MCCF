package controllers

import (
	"fmt"
	"net/http"
	"tubes2-be-mccf/internal/models"
	"tubes2-be-mccf/internal/utils"

	"encoding/json"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

// Result Request Data Structure
type PlayRequest struct {
	Algorithm string `json:"algorithm" binding:"required,oneof='IDS' 'BFS'"` // Algorithm type (IDS or BFS)
	Start     string `json:"start" binding:"required"`                       // Start wikipedia article title
	Target    string `json:"target" binding:"required"`                      // Target wikipedia article title
}

// Result Response Data Structure
type PlaySuccessResponse struct {
	TotalTraversed     int              `json:"totalTraversed"`     // Total article traversed
	ShortestPathLength int              `json:"shortestPathLength"` // Shortest path length
	Duration           float32          `json:"duration"`           // Duration of the search
	Articles           []models.Article `json:"articles"`           // Articles data used in paths: Path[]
	Paths              []models.Path    `json:"paths"`              // List of the shortest paths from start to target found
}

// Error Response Data Structure
type PlayErrorResponse struct {
	Field   string `json:"field"`   // Form field that caused the error
	Message string `json:"message"` // Error message
}

type BacklinkResponse struct {
	BatchComplete string `json:"batchcomplete"`
	Continue      struct {
		BlContinue string `json:"blcontinue"`
		Cont       string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Backlinks []struct {
			PageId int    `json:"pageid"`
			Ns     int    `json:"ns"`
			Title  string `json:"title"`
		} `json:"backlinks"`
	} `json:"query"`
}

// Get all backlinks from a wikipedia URL.

func getBackLinkTitles(title string) []string {
	baseUrl := "https://en.wikipedia.org/w/api.php"
	params := url.Values{
		"action":  {"query"},
		"format":  {"json"},
		"list":    {"backlinks"},
		"bltitle": {title},
		"bllimit": {"max"},
	}
	// create array of string
	var backlinks []string

	count := 0
	for {
		resp, err := http.Get(baseUrl + "?" + params.Encode())
		if err != nil {
			panic(err)
		}

		var response BacklinkResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		for _, backlink := range response.Query.Backlinks {
			// fmt.Println(backlink.Title)
			// backlinks = append(backlinks, backlink.Title)
			title := strings.ReplaceAll(backlink.Title, " ", "_")
			title = "https://en.wikipedia.org/wiki/" + title
			backlinks = append(backlinks, title)
			count++
		}

		if response.Continue.BlContinue != "" {
			params.Set("blcontinue", response.Continue.BlContinue)
		} else {
			break
		}
	}

	// fmt.Println(count)
	return backlinks
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

func getPath(path map[string]string, startURL string, endURL string) []string {
	result := []string{endURL}
	current := endURL
	for current != startURL {
		current = path[current]
		result = append(result, current)
	}
	// rever result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func IDS(startURL string, targetURL string) []string {
	path :=make(map[string]string)
	visited := make(map[string]bool)
	// depth := 1
	// for {
	if DLS(startURL, targetURL, visited, path, 2) {
		// fmt.Println("found")
		return getPath(path, startURL, targetURL)
	}
	return nil
	// }
	// return path

}

func DLS(startURL string, targetURL string, visited map[string]bool, path map[string]string, depth int) bool {
	// Mark the current URL as visited
	visited[startURL] = true


	if startURL == targetURL {
		return true
	}
	if depth == 0 {
		return false
	}
	// If the current URL is the target URL, return the path

	// Get all internal links from the current URL
	links := getAllInternalLinks(startURL)
	fmt.Println("current processed : ", startURL)
	for _, link := range links {
		fmt.Println(link)
	}
	fmt.Println("depth : ", depth)

	// For each link, if it has not been visited, call DFS recursively
	for _, link := range links {
		if !visited[link] {
			path[link] = startURL
			result := DLS(link, targetURL, visited, path, depth-1)
			if(result){
				return true
			}
		}
	}

	return false
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get data
	algorithm := reqJSON.Algorithm
	startTitle := reqJSON.Start
	targetTitle := reqJSON.Target

	// Get start wikipedia URL (and validate title)
	startURL, err := getWikipediaURLFromTitle(startTitle)
	if err != nil {
		c.JSON(400, gin.H{"field": "start", "message": err.Error()})
	}
	// Get target wikipedia URL (and validate title)
	targetURL, err := getWikipediaURLFromTitle(targetTitle)
	if err != nil {
		c.JSON(400, gin.H{"field": "target", "message": err.Error()})
	}

	// Solve
	result, err := Solve(algorithm, startURL, targetURL)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	// Return the result
	c.JSON(200, result)
}


func main(){

	// StartURL = "https://en.wikipedia.org/wiki/Computer_Science"
	// EndURL = "https://en.wikipedia.org/wiki/Joko_Widodo"




}