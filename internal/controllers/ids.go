package controllers

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const maxConcurrent = 100

var mu sync.Mutex
var wg sync.WaitGroup

// var cache = make(map[string][]string)

// var depthh = 1

func IDS(startURL string, targetURL string) ([][]string, int32) {

	resultPath := make([][]string, 0)
	cache := make(map[string][]string)
	path := make([]string, 0)
	var totalTraversed int32 = 0
	gm := NewGoRoutineManager(maxConcurrent)
	depth := 1
	for {
		fmt.Println("===============================================")
		fmt.Println("depth : ", depth)
		fmt.Println("==============================================")
		// wg.Add(1)
		DLS(startURL, targetURL, path, &resultPath, depth, gm, &totalTraversed, &cache)
		// wg.Done()/
		wg.Wait()
		if len(resultPath) > 0 {
			fmt.Println("found")
			// fmt.Println(len(resultPath))

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
		totalTraversed = 0
	}
	return nil, 0

}

func DLS(startURL string, targetURL string, path []string, resultpath *[][]string, depth int, gm *goRoutineManager, totalTraversed *int32, cache *map[string][]string) {

	atomic.AddInt32(totalTraversed, 1)
	if startURL == targetURL {
		// mu.Lock()
		mu.Lock()
		*resultpath = append(*resultpath, path)

		mu.Unlock()

		return
	}
	if depth == 0 {

		return
	}

	var links []string
	if depth > 1 {
		// links = (&cache[startURL])
		links = (*cache)[startURL]
	} else {
		// check if startURL is in cache
		// mu.Lock()
		// _, ok := cache[startURL]
		// mu.Unlock()
		// if ok {
		// 	return
		// }

		links = getAllInternalLinks(startURL)
		mu.Lock()
		(*cache)[startURL] = links
		// visited[startURL] = true

		mu.Unlock()
	}

	fmt.Println("current processed : ", startURL)

	// fmt.Println("depth : ", depth)

	for _, link := range links {
		currpath := append(path, link)

		// capture the link so each goroutine is unique
		link := link

		gm.Run(func() {

			DLS(link, targetURL, currpath, resultpath, depth-1, gm, totalTraversed, cache)

		})

	}
	// wg.Wait()
	return

}

// func IDSSingle(startURL string, targetURL string) ([][]string, int32) {
// 	// resultPath := make([][]string, 0)
// 	cache := make(map[string][]string)
// 	path := make([]string, 0)
// 	var totalTraversed int32 = 0
// 	gm := NewGoRoutineManager(maxConcurrent)
// 	depth := 1
// 	for {
// 		fmt.Println("===============================================")
// 		fmt.Println("depth : ", depth)
// 		fmt.Println("==============================================")
// 		// wg.Add(1)
// 		res := DLSSingle(startURL, targetURL, path, depth, gm, &totalTraversed, &cache)
// 		// DLS(startURL, targetURL, path,depth, gm, &totalTraversed, &cache)
// 		// wg.Done()/
// 		// wg.Wait()
// 		if res != nil {
// 			fmt.Println("found")
// 			// fmt.Println(len(resultPath))

// 			fmt.Println("res : ", res)

// 			return [][]string{res} , totalTraversed
// 		}

// 		path = path[:0]
// 		if depth > 10 {
// 			break
// 		}
// 		depth++
// 		totalTraversed = 0
// 	}
// 	return nil, 0

// }

func IDSSingle(startURL string, targetURL string) ([][]string, int32) {
	resultPath := make([][]string, 0)
	cache := make(map[string][]string)
	path := make([]string, 0)
	var totalTraversed int32 = 0
	gm := NewGoRoutineManager(maxConcurrent)
	depth := 1
	for {
		fmt.Println("===============================================")
		fmt.Println("depth : ", depth)
		fmt.Println("==============================================")
		// wg.Add(1)
		DLSSingle(startURL, targetURL, path, &resultPath, depth, gm, &totalTraversed, &cache)
		// wg.Done()/
		wg.Wait()
		if len(resultPath) > 0 {
			fmt.Println("found")
			// fmt.Println(len(resultPath))

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
		totalTraversed = 0
	}
	return nil, 0

}

var targetFound int32 = 0

func DLSSingle(startURL string, targetURL string, path []string, resultpath *[][]string, depth int, gm *goRoutineManager, totalTraversed *int32, cache *map[string][]string) {
	if targetFound != 0 {
		return
	}
	atomic.AddInt32(totalTraversed, 1)
	if startURL == targetURL {
		// mu.Lock()
		targetFound++
		mu.Lock()
		*resultpath = append(*resultpath, path)
		fmt.Println(targetFound)
		mu.Unlock()

		return
	}
	if depth == 0 {

		return
	}

	var links []string
	if depth > 1 {
		// links = (&cache[startURL])
		links = (*cache)[startURL]
	} else {
		// check if startURL is in cache
		// mu.Lock()
		// _, ok := cache[startURL]
		// mu.Unlock()
		// if ok {
		// 	return
		// }

		links = getAllInternalLinks(startURL)
		mu.Lock()
		(*cache)[startURL] = links
		// visited[startURL] = true

		mu.Unlock()
	}

	fmt.Println("current processed : ", startURL)

	// fmt.Println("depth : ", depth)

	for _, link := range links {
		currpath := append(path, link)

		// capture the link so each goroutine is unique
		link := link

		gm.Run(func() {

			DLSSingle(link, targetURL, currpath, resultpath, depth-1, gm, totalTraversed, cache)

		})

	}
	// wg.Wait()
	return

}
