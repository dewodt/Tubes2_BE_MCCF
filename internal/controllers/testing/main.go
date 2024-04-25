package main

import (
	"fmt"
	"time"
	"tubes2-be-mccf/internal/controllers"
	"tubes2-be-mccf/internal/utils"
)

func IDStest() {
	// startURL := "https://en.wikipedia.org/wiki/Umbraculidae"
	// endURL := "https://en.wikipedia.org/wiki/Chicken"
	startURL := "https://en.wikipedia.org/wiki/Inauguration_of_Joko_Widodo"
	endURL := "https://en.wikipedia.org/wiki/Indonesia"

	// startURL := "https://en.wikipedia.org/wiki/Chicken"
	// endURL := "https://en.wikipedia.org/wiki/Joko_Widodo"

	// startURL := "https://en.wikipedia.org/wiki/New_car_smell"
	// endURL := "https://en.wikipedia.org/wiki/Venus_flytrap"
	// startURL := "https://en.wikipedia.org/wiki/Joko_Widodo"
	// endURL := "https://en.wikipedia.org/wiki/Rengasdengklok_Incident"

	// res := getAllInternalLinks(startURL)
	// fmt.Println(res)
	// results := getThumbnail(startURL)
	startTime := time.Now()
	// // fmt.Println(res)
	results, traversed := controllers.IDS(startURL, endURL)
	elapsedTime := time.Since(startTime)

	// fmt.Println(results)
	// print(results)
	utils.PrintResultPath(results)
	fmt.Println("Elapsed Time: ", elapsedTime)
	fmt.Println("Total Traversed: ", traversed)
}

func BFStest() {
	// startURL := "https://en.wikipedia.org/wiki/Prince_Edward_Viaduct"
	// endURL := "https://en.wikipedia.org/wiki/Humber_Cinemas"
	// startURL := "https://en.wikipedia.org/wiki/Humber_Cinemas"
	// endURL := "https://en.wikipedia.org/wiki/Prince_Edward_Viaduct"

	startURL := "https://en.wikipedia.org/wiki/Inauguration_of_Joko_Widodo"
	endURL := "https://en.wikipedia.org/wiki/Indonesia"
	// startURL := "https://en.wikipedia.org/wiki/Prabowo_Subianto"
	// endURL := "https://en.wikipedia.org/wiki/Joko_Widodo"

	// startURL := "https://en.wikipedia.org/wiki/Prabowo_Subianto"
	// endURL := "https://en.wikipedia.org/wiki/Joko_Widodo"

	// res := getAllInternalLinks(startURL)
	// fmt.Println(res)
	// results := getThumbnail(startURL)
	startTime := time.Now()
	// // fmt.Println(res)
	results, path := controllers.BFS(startURL, endURL)
	elapsedTime := time.Since(startTime)

	// fmt.Println(results)
	// print(results)
	utils.PrintResultPath(results)
	fmt.Println("Path: ", path)
	fmt.Println("Elapsed Time: ", elapsedTime)
	// fmt.Println("Total Traversed: ", traversed)
}

func main() {
	BFStest()
}
