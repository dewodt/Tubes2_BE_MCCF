package controllers

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const maxConcurrent = 200

var mu sync.Mutex
var wg sync.WaitGroup

var cache = make(map[string][]string)

// var depthh = 1

func IDS(startURL string, targetURL string) ([][]string, int32) {
	resultPath := make([][]string, 0)
	// cache := make(map[string][]string)
	path := make([]string, 0)
	var totalTraversed int32 = 0
	gm := NewGoRoutineManager(maxConcurrent)
	depth := 1
	for {
		// wg.Add(1)
		DLS(startURL, targetURL, path, &resultPath, depth, gm, &totalTraversed)
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
	}
	return nil, 0

}

// var visited = make(map[string]bool)

func DLS(startURL string, targetURL string, path []string, resultpath *[][]string, depth int, gm *goRoutineManager, totalTraversed *int32) {

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

	// links := getAllInternalLinks(startURL)

	// mu.Lock()
	// mu.lock()
	// checks := (*cache)[startURL] == nil
	// mu.Unlock()
	// add links to cache
	var links []string
	if depth > 1 {
		links = cache[startURL]
	} else {
		// check if startURL is in cache
		mu.Lock()
		_, ok := cache[startURL]
		mu.Unlock()
		if ok {
			return
		}

		links = getAllInternalLinks(startURL)
		mu.Lock()
		cache[startURL] = links
		// visited[startURL] = true

		mu.Unlock()
	}

	// if checks {
	// 	mu.Lock()
	// 	(*cache)[startURL] = links
	// 	mu.Unlock()
	// } else {
	// 	// mu.Lock()
	// 	links = (*cache)[startURL]
	// 	// mu.Unlock()
	// }

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

	return

}
