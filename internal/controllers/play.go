package controllers

import (
	"fmt"
	"net/http"
	"os"
	"sync"
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

const maxConcurrent = 50

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

type goRoutineManager struct {
	goRoutineCnt chan bool
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
	// var links []string
	links := make(map[string]bool)

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
		if len(link) > 6 && link[:6] == "/wiki/" {

			fullLink := "https://en.wikipedia.org" + link

			// If the link is not in the map, add it
			if _, ok := links[fullLink]; !ok {
				links[fullLink] = true
			}
		}
	})

	// Visit the URL

	err := cl.Visit(url)
	if err != nil {
		fmt.Println("error visit")
		os.Exit(1)

	}
	// check if error

	result := make([]string, 0, len(links))
	for link := range links {
		result = append(result, link)
	}

	return result
}

func printResultpath(resultPath [][]string) {
	fmt.Println("Found", len(resultPath), "paths : ")
	for i, path := range resultPath {
		fmt.Print("Path ", i+1, " : ")
		fmt.Println(path)
	}
}

func (g *goRoutineManager) Run(f func()) {
	select {
	case g.goRoutineCnt <- true:
		go func() {
			f()
			<-g.goRoutineCnt
		}()
	default:
		f()
	}
}

func NewGoRoutineManager(goRoutineLimit int) *goRoutineManager {
	return &goRoutineManager{
		goRoutineCnt: make(chan bool, goRoutineLimit),
	}
}

func IDS(startURL string, targetURL string) [][]string {
	resultPath := make([][]string, 0)
	path := make([]string, 0)
	cache := make(map[string][]string)
	mu := sync.Mutex{}
	depth := 1
	gm := NewGoRoutineManager(maxConcurrent)
	for {

		DLS(startURL, targetURL, path, &resultPath, depth, &cache, gm, &mu)

		if len(resultPath) > 0 {
			fmt.Println("found")
			fmt.Println(len(resultPath))

			for i := range resultPath {
				resultPath[i] = append([]string{startURL}, resultPath[i]...)
			}

			return resultPath
		}

		path = path[:0]
		if depth > 10 {
			break
		}
		depth++
	}
	return nil

}

func DLS(startURL string, targetURL string, path []string, resultpath *[][]string, depth int, cache *map[string][]string, gm *goRoutineManager, mu *sync.Mutex) {
	if startURL == targetURL {
		mu.Lock()
		*resultpath = append(*resultpath, path)
		mu.Unlock()

		return
	}
	if depth == 0 {

		return
	}

	var links []string

	if (*cache)[startURL] == nil {
		links = getAllInternalLinks(startURL)
		mu.Lock()
		(*cache)[startURL] = links
		mu.Unlock()
	} else {
		links = (*cache)[startURL]
	}

	fmt.Println("current processed : ", startURL)

	fmt.Println("depth : ", depth)
	for _, link := range links {
		currpath := append(path, link)
		gm.Run(func() {
			DLS(link, targetURL, currpath, resultpath, depth-1, cache, gm, mu)
		})

	}

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

func main() {

	// StartURL = "https://en.wikipedia.org/wiki/Computer_Science"
	// EndURL = "https://en.wikipedia.org/wiki/Joko_Widodo"

}
