package controllers

import (
	"fmt"
	"time"
	"tubes2-be-mccf/internal/models"
	"github.com/gin-gonic/gin"
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
func solveBFS(startURL string, targetURL string) (PlaySuccessResponse, error) {
	fmt.Println("Solving with BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)
	// var adj [][]int

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
