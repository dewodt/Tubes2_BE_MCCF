package controllers

import (
	"fmt"
	"sync"
	"sync/atomic"
	"tubes2-be-mccf/internal/utils"
)

const maxConcurrent = 200

// mutex so race condition doesnt happen
var mu sync.Mutex
// waitgroup so we can wait for all goroutine to finish before continuing to the next IDS iteration
var wg sync.WaitGroup
// targetfound used for single path IDS
var targetFound int32 = 0

func IDS(startURL string, targetURL string, isSingle bool) ([][]string, int32) {

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
		if(isSingle){
			DLSSingle(startURL, targetURL, path, &resultPath, depth, gm, &totalTraversed, &cache)

		}else{
			DLS(startURL, targetURL, path, &resultPath, depth, gm, &totalTraversed, &cache)
		}

		wg.Wait()
		if len(resultPath) > 0 {
			fmt.Println("found")

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
		links = (*cache)[startURL]
	} else {
		links = utils.GetAllInternalLinks(startURL)
		mu.Lock()
		(*cache)[startURL] = links
		mu.Unlock()
	}

	fmt.Println("current processed : ", startURL)

	for _, link := range links {
		currpath := append(path, link)

		// capture the link so each goroutine is unique
		link := link
		gm.Run(func() {

			DLS(link, targetURL, currpath, resultpath, depth-1, gm, totalTraversed, cache)

		})
	}
	
}



func DLSSingle(startURL string, targetURL string, path []string, resultpath *[][]string, depth int, gm *goRoutineManager, totalTraversed *int32, cache *map[string][]string) {
	if targetFound != 0 {
		return
	}
	atomic.AddInt32(totalTraversed, 1)
	if startURL == targetURL {
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
	
		links = (*cache)[startURL]
	} else {
		links = utils.GetAllInternalLinks(startURL)
		mu.Lock()
		(*cache)[startURL] = links


		mu.Unlock()
	}

	fmt.Println("current processed : ", startURL)

	for _, link := range links {
		currpath := append(path, link)

		// capture the link so each goroutine is unique
		link := link

		gm.Run(func() {

			DLSSingle(link, targetURL, currpath, resultpath, depth-1, gm, totalTraversed, cache)

		})

	}

}
