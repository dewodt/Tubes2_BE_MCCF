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

	cl.Visit(url)

	result := make([]string, 0, len(links))
	for link := range links {
		result = append(result, link)
	}

	return result
}

func printResultpath(resultPath [][]string) {
	for _, path := range resultPath {
		fmt.Println(path)
	}
}

func IDS(startURL string, targetURL string) [][]string {
	resultPath := make([][]string, 0)
	path := make([]string, 0)
	cache:= make(map[string][]string)
	// visited := make(map[string]bool)
	depth := 1
	for {
		if DLS(startURL, targetURL, path, &resultPath, depth,&cache) {
			fmt.Println("found")
			// fmt.Println(len(path))
			// fmt.Println(resultPath)
			printResultpath(resultPath)
			fmt.Println(len(resultPath))
			// insert startURL at the start of path
			for i := range resultPath {
				resultPath[i] = append([]string{startURL}, resultPath[i]...)
			} 
			
			return resultPath
		}
		depth++
	}
	return nil
	// }
	// return path

}

func DLS(startURL string, targetURL string, path []string, resultpath *[][]string, depth int, cache *map[string][]string) bool {
	

	if startURL == targetURL {
		*resultpath = append(*resultpath, path)
		return true
	}
	if depth == 0 {
		return false
	}

	
	var links []string
	
	if(*cache)[startURL] == nil{
		links = getAllInternalLinks(startURL)
		(*cache)[startURL] = links
	}else{
		links = (*cache)[startURL]
	}
	
	fmt.Println("current processed : ", startURL)

	fmt.Println("depth : ", depth)

	
	result := false
	flag := false
	for _, link := range links {
		currpath := append(path, link)
		result = DLS(link, targetURL, currpath, resultpath, depth-1, cache)
		if result {
			flag = true
		}
	}
	return flag
	// return false
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
