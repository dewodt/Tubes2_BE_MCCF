package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
	"tubes2-be-mccf/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

type Queue struct {
	Elements []string
	Size     int
}

func (q *Queue) Enqueue(elem string) {
	if q.GetLength() == q.Size {
		fmt.Println("Overflow")
		return
	}
	q.Elements = append(q.Elements, elem)
}

func (q *Queue) Dequeue() string {
	if q.IsEmpty() {
		fmt.Println("UnderFlow")
		return ""
	}
	element := q.Elements[0]
	if q.GetLength() == 1 {
		q.Elements = nil
		return element
	}
	q.Elements = q.Elements[1:]
	return element // Slice off the element once it is dequeued.
}

func (q *Queue) GetLength() int {
	return len(q.Elements)
}

func (q *Queue) IsEmpty() bool {
	return len(q.Elements) == 0
}

func (q *Queue) Peek() (string, error) {
	if q.IsEmpty() {
		return "", errors.New("empty queue")
	}
	return q.Elements[0], nil
}

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

const maxConcurrent = 150

// Error Response Data Structure
type PlayErrorResponse struct {
	Error       string       `json:"error"`       // Error message
	Message     string       `json:"message"`     // Error message
	ErrorFields []FieldError `json:"errorFields"` // List of fields that caused the error
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

type goRoutineManager struct {
	goRoutineCnt chan bool
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
	cl.OnHTML("a[href]", func(e *colly.HTMLElement) {
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

func IDS(startURL string, targetURL string) ([][]string, int32) {
	resultPath := make([][]string, 0)
	path := make([]string, 0)

	depth := 1
	gm := NewGoRoutineManager(maxConcurrent)
	var totalTraversed int32 = 0
	for {

		DLS(startURL, targetURL, path, &resultPath, depth, gm, &totalTraversed)

		if len(resultPath) > 0 {
			fmt.Println("found")
			fmt.Println(len(resultPath))

			for i := range resultPath {
				resultPath[i] = append([]string{startURL}, resultPath[i]...)
			}

			return resultPath, totalTraversed
		}

		path = path[:0]
		if depth > 10 {
			break
		}
		depth++
	}
	return nil, 0

}

func DLS(startURL string, targetURL string, path []string, resultpath *[][]string, depth int, gm *goRoutineManager, totalTraversed *int32) {
	atomic.AddInt32(totalTraversed, 1)
	if startURL == targetURL {
		*resultpath = append(*resultpath, path)
		return
	}
	if depth == 0 {

		return
	}

	links := getAllInternalLinks(startURL)

	fmt.Println("current processed : ", startURL)

	fmt.Println("depth : ", depth)
	for _, link := range links {
		currpath := append(path, link)

		// capture the link so each goroutine is unique
		link := link
		gm.Run(func() {
			DLS(link, targetURL, currpath, resultpath, depth-1, gm, totalTraversed)
		})

	}

}

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

func getTitleFromURL(url string) string {

	title := url[30:]

	title = strings.ReplaceAll(title, "_", " ")
	return title
}

func getArticlesFromResultPath(path [][]string) []models.Article {
	//
	articleSet := make(map[string]bool)
	articles := make([]models.Article, 0)

	id := 0
	for _, p := range path {

		for _, link := range p {
			if _, ok := articleSet[link]; !ok {
				articleSet[link] = true
				articles = append(articles, models.Article{
					ID:          id,
					Title:       getTitleFromURL(link),
					Description: "",
					Thumbnail:   getThumbnail(link),
					URL:         link,
				})
				id++
			}

		}
	}
	return articles

}

func getPathsFromResultPath(path [][]string, articles []models.Article) []models.Path {
	paths := make([]models.Path, 0)
	for _, p := range path {
		pathRes := make([]int, 0)
		for _, link := range p {
			for _, article := range articles {
				if article.URL == link {
					pathRes = append(pathRes, article.ID)
					break
				}
			}
		}
		paths = append(paths, pathRes)
	}
	return paths
}

func SolveIDS(startURL string, targetURL string) (PlaySuccessResponse, error) {
	fmt.Println("Solving with IDS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)
	startTime := time.Now()
	resultPath, totalTraversed := IDS(startURL, targetURL)

	elapseTime := time.Since(startTime).Seconds() * 1000

	// Placeholder
	if len(resultPath) == 0 {
		return PlaySuccessResponse{}, nil
	} else {
		articles := getArticlesFromResultPath(resultPath)
		paths := getPathsFromResultPath(resultPath, articles)
		return PlaySuccessResponse{
			TotalTraversed:     int(totalTraversed),
			ShortestPathLength: len(resultPath[0]),
			Duration:           float32(elapseTime),
			Articles:           articles,
			Paths:              paths,
		}, nil
	}
}
func dfs(paths [][]string, path []string, parent map[string][]string, end string) {
	if parent[end] == nil {
		path = append(path, end)
		paths = append(paths, path)
		path = path[:len(path)-1]
	} else {
		for i := 0; i < len(parent[end]); i++ {
			path = append(path, end)
			dfs(paths, path, parent, parent[end][i])
			path = path[:len(path)-1]
		}
	}
}
func solveBFS(startURL string, targetURL string) (PlaySuccessResponse, error) {
	fmt.Println("Solving with BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)
	// var adj [][]int
	maxInt := math.MaxInt32
	adj := make(map[string][]string)
	parent := make(map[string][]string)
	parent[startURL] = nil
	q := Queue{Size: 0}
	dist := make(map[string]int)
	dist[startURL] = 0
	dist[targetURL] = maxInt
	//making bfs tree
	for !q.IsEmpty() {
		u, err := q.Peek()
		if err != nil {
			fmt.Println("Queue is empty")
		}
		q.Dequeue()
		if dist[u] >= dist[targetURL] {
			continue
		}
		links := getAllInternalLinks(startURL)
		for i := 0; i < len(links); i++ {
			adj[u] = append(adj[u], links[i])
		}
		for _, v := range links {
			if v != startURL && dist[v] == 0 {
				dist[v] = maxInt
			}
		}
		for i := 0; i < len(links); i++ {
			if dist[adj[u][i]] > dist[u]+1 {
				dist[adj[u][i]] = dist[u] + 1
				q.Enqueue(adj[u][i])
				//parent[adj[u][i]].clear(),push_back
				parent[adj[u][i]] = nil
				parent[adj[u][i]] = append(parent[adj[u][i]], u)
			} else if dist[adj[u][i]] == dist[u]+1 {
				//parent[adj[u][i]].pushback
				parent[adj[u][i]] = append(parent[adj[u][i]], u)
			}
		}
	}
	//change bfs tree to array of array of solution
	var paths [][]string
	var path []string

	dfs(paths, path, parent, targetURL)

	//fill solution type with solution

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
	// print the json
	// fmt.Println(result)
	fmt.Println(result.Articles)
	fmt.Println(result.Paths)
}
