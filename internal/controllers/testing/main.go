package main

import (
	"fmt"
	"time"
	"tubes2-be-mccf/internal/controllers"
)

func main() {
	// startURL := "https://en.wikipedia.org/wiki/Umbraculidae"
	// endURL := "https://en.wikipedia.org/wiki/Chicken"
	startURL := "https://en.wikipedia.org/wiki/Humber_Cinemas"
	endURL := "https://en.wikipedia.org/wiki/Prince_Edward_Viaduct"

	// startURL := "https://en.wikipedia.org/wiki/Chicken"
	// endURL := "https://en.wikipedia.org/wiki/Joko_Widodo"

	// res := getAllInternalLinks(startURL)
	// fmt.Println(res)
	// results := getThumbnail(startURL)
	startTime := time.Now()
	// // fmt.Println(res)
	results, traversed := controllers.IDS(startURL, endURL)
	elapsedTime := time.Since(startTime)

	// fmt.Println(results)
	// print(results)
	controllers.PrintResultPath(results)
	fmt.Println("Elapsed Time: ", elapsedTime)
	fmt.Println("Total Traversed: ", traversed)
}
