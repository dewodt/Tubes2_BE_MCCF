package utils

import (
	"fmt"
	"tubes2-be-mccf/internal/models"
)

func GetArticlesAndPaths(path [][]string) ([]models.Article, []models.Path) {
	// Initialize
	var articles []models.Article
	var paths []models.Path

	// Initialize article[url]->index map
	articleSet := make(map[string]int)
	index := 0

	// Iterate over each path
	for _, p := range path {
		fmt.Println("Original Path: ")
		fmt.Println(p) // Array of links
		// Initialize pathRes
		pathRes := make([]int, 0)

		// Iterate over each link in path
		for _, link := range p {
			// Check link is not in articleSet
			if mapIdx, ok := articleSet[link]; !ok {
				// Get article data from link
				article := GetArticleData(link)

				// Add article to articles
				articles = append(articles, article)

				// Update map
				articleSet[link] = index

				// Add index to pathRes
				pathRes = append(pathRes, index)

				// Add index
				index++
			} else {
				// Add index to pathRes
				pathRes = append(pathRes, mapIdx)
			}
		}

		// Add pathRes to paths
		paths = append(paths, pathRes)
	}

	return articles, paths
}
