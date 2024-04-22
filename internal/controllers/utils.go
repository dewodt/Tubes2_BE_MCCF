package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"tubes2-be-mccf/internal/models"
)

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

func PrintResultPath(resultPath [][]string) {
	fmt.Println("Found", len(resultPath), "paths : ")
	for i, path := range resultPath {
		fmt.Print("Path ", i+1, " : ")
		fmt.Println(path)
	}
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
