package controllers

import (
	"errors"
	"fmt"
	"math"
	"tubes2-be-mccf/internal/cache"
	"tubes2-be-mccf/internal/utils"
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
func reverse(arr []string) []string {
	var ans []string

	for i := len(arr) - 1; i >= 0; i-- {
		ans = append(ans, arr[i])
	}
	return ans
}
func dfs(paths [][]string, path []string, parent map[string][]string, end string) [][]string {
	if parent[end] == nil {
		path = append(path, end)
		// *path = append(*path, end)

		paths = append(paths, path)
		// *paths = append(*paths, *path)
		path = path[:len(path)-1]
		// *path = (*path)[:len(*path)-1]
	} else {
		for i := 0; i < len(parent[end]); i++ {
			path = append(path, end)
			// *path = append(*path, end)
			paths = dfs(paths, path, parent, parent[end][i])
			path = path[:len(path)-1]
			// *path = (*path)[:len(*path)-1]
		}
	}
	return paths
}

const maxConcurrentBFS = 450

func BFS(startURL string, targetURL string, isSingle bool) ([][]string, int) {
	if isSingle {
		return BFSSingle(startURL, targetURL)
	} else {
		return BFSMulti(startURL, targetURL)
	}
}

func BFSMulti(startURL string, targetURL string) ([][]string, int) {
	fmt.Println("Solving with BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)

	traversed := 1

	if startURL == targetURL {
		return [][]string{{startURL}}, 1
	}

	gm := NewGoRoutineManager(maxConcurrentBFS)
	maxInt := math.MaxInt32
	adj := make(map[string][]string)
	parent := make(map[string][]string)
	parent[startURL] = nil
	q := Queue{Size: maxInt}
	dist := make(map[string]int)
	dist[startURL] = 0
	dist[targetURL] = maxInt
	q.Enqueue(startURL)

	for !q.IsEmpty() {
		var isFirst bool
		isFirst = false
		for i := 0; i < q.GetLength(); i++ {
			i := i
			gm.Run(func() {
				check := dist[q.Elements[i]] < dist[targetURL]

				if check {
					// rwmu.RLock()
					cacheCheck := cache.Cache[q.Elements[i]]
					// rwmu.RUnlock()
					if cacheCheck == nil {
						links := utils.GetAllInternalLinks(q.Elements[i])
						rwmu.Lock()
						adj[q.Elements[i]] = links
						rwmu.Unlock()
					} else {
						rwmu.Lock()
						adj[q.Elements[i]] = cacheCheck
						rwmu.Unlock()
					}
					fmt.Println(q.Elements[i], "dist: ", dist[q.Elements[i]], "target URL: ", dist[targetURL])

				}
			})
		}
		wg.Wait()
		length := q.GetLength()
		for i := 0; i < length; i++ {
			u := q.Elements[i]
			traversed++
			if dist[u] >= dist[targetURL] {
				continue
			}
			for _, v := range adj[u] {
				if v != startURL && dist[v] == 0 {
					dist[v] = maxInt
				}
			}

			for i := 0; i < len(adj[u]); i++ {

				if dist[adj[u][i]] > dist[u]+1 {
					dist[adj[u][i]] = dist[u] + 1
					q.Enqueue(adj[u][i])
					parent[adj[u][i]] = nil
					parent[adj[u][i]] = append(parent[adj[u][i]], u)
				} else if dist[adj[u][i]] == dist[u]+1 {
					parent[adj[u][i]] = append(parent[adj[u][i]], u)
				}
				if dist[targetURL] == 1 {
					isFirst = true
					break
				}
			}
			if isFirst {
				break
			}

		}
		if isFirst {
			break
		}

		q.Elements = q.Elements[length:]
	}

	paths := make([][]string, 0)
	path := make([]string, 0)

	paths = dfs(paths, path, parent, targetURL)

	for i := 0; i < len(paths); i++ {
		paths[i] = reverse(paths[i])
	}
	// cache.UpdateMapInFile("./cache/cache.json")
	// updateMapInFile(cache, "../cache/cache.json")
	return paths, traversed

}

func BFSSingle(startURL string, targetURL string) ([][]string, int) {
	fmt.Println("Solving with BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// var adj [][]int
	traversed := 1

	if startURL == targetURL {
		return [][]string{{startURL}}, 1
	}

	// cache := make(map[string][]string)
	// // cache, err := readMapFromFile("../cache/cache.json")
	// cache, err := readMapFromFile("./cache/cache.json")

	// if err != nil {
	// 	fmt.Println("error reading cache file")
	// 	cache = make(map[string][]string)
	// }

	gm := NewGoRoutineManager(maxConcurrentBFS)
	maxInt := math.MaxInt32
	adj := make(map[string][]string)
	parent := make(map[string][]string)

	parent[startURL] = nil
	q := Queue{Size: maxInt}
	dist := make(map[string]int)
	dist[startURL] = 0
	dist[targetURL] = maxInt
	q.Enqueue(startURL)
	isFound := false
	for !q.IsEmpty() {
		length := min(q.GetLength(), 200)
		for i := 0; i < length; i++ {
			i := i
			gm.Run(func() {
				check := dist[q.Elements[i]] < dist[targetURL]
				// mu.Unlock()
				if check {
					// rwmu.RLock()
					cacheCheck := cache.Cache[q.Elements[i]]
					// rwmu.RUnlock()
					if cacheCheck == nil {
						links := utils.GetAllInternalLinks(q.Elements[i])
						rwmu.Lock()
						// cache.Cache[q.Elements[i]] = links
						adj[q.Elements[i]] = links
						rwmu.Unlock()
					} else {
						rwmu.Lock()
						adj[q.Elements[i]] = cacheCheck
						rwmu.Unlock()
					}
					fmt.Println(q.Elements[i], "dist: ", dist[q.Elements[i]], "target URL: ", dist[targetURL])

				}
			})
		}
		wg.Wait()
		// length := q.GetLength()
		for i := 0; i < length; i++ {
			u := q.Elements[i]
			if dist[u] >= dist[targetURL] {
				continue
			}
			for _, v := range adj[u] {
				if v != startURL && dist[v] == 0 {
					dist[v] = maxInt
				}
			}
			for i := 0; i < len(adj[u]); i++ {
				traversed++

				if dist[adj[u][i]] > dist[u]+1 {
					dist[adj[u][i]] = dist[u] + 1
					q.Enqueue(adj[u][i])
					parent[adj[u][i]] = nil
					parent[adj[u][i]] = append(parent[adj[u][i]], u)
				}

				if dist[targetURL] != maxInt {
					isFound = true
					break
				}
			}
			if isFound {
				break
			}

		}
		if isFound {
			break
		}
		q.Elements = q.Elements[length:]
	}

	paths := make([][]string, 0)
	path := make([]string, 0)

	paths = dfs(paths, path, parent, targetURL)

	for i := 0; i < len(paths); i++ {
		paths[i] = reverse(paths[i])
	}
	// cache.UpdateMapInFile("./cache/cache.json")
	// updateMapInFile(cache, "../cache/cache.json")
	paths = paths[:1]
	return paths, traversed

}

// func BFS(startURL string, targetURL string, isSingle bool) ([][]string, int) {
// 	fmt.Println("Solving with BFS")
// 	fmt.Println("Start URL:", startURL)
// 	fmt.Println("Target URL:", targetURL)
// 	// runtime.GOMAXPROCS(runtime.NumCPU())
// 	// var adj [][]int
// 	traversed := 1

// 	if startURL == targetURL {
// 		return [][]string{{startURL}}, 1
// 	}
// 	cache := make(map[string][]string)

// 	// cache, err := readMapFromFile("./internal/controllers/cache/cache.json")
// 	cache, err := readMapFromFile("../cache/cache.json")

// 	if err != nil {
// 		fmt.Println("error reading cache file")
// 		cache = make(map[string][]string)
// 	}

// 	gm := NewGoRoutineManager(200)
// 	maxInt := math.MaxInt32
// 	adj := make(map[string][]string)
// 	parent := make(map[string][]string)
// 	parent[startURL] = nil
// 	q := Queue{Size: maxInt}
// 	dist := make(map[string]int)
// 	dist[startURL] = 0
// 	dist[targetURL] = maxInt
// 	q.Enqueue(startURL)
// 	var isFound bool
// 	var isFirst bool
// 	isFirst = false
// 	isFound = false
// 	for !q.IsEmpty() {
// 		length := min(q.GetLength(),q.GetLength())
// 		for i := 0; i < length; i++ {
// 			i := i
// 			gm.Run(func() {
// 				// mu.Lock()
// 				check := dist[q.Elements[i]] < dist[targetURL]
// 				// mu.Unlock()
// 				if check {
// 					rwmu.RLock()
// 					cacheCheck := cache[q.Elements[i]]
// 					rwmu.RUnlock()
// 					if cacheCheck == nil {
// 						links := utils.GetAllInternalLinks(q.Elements[i])
// 						rwmu.Lock()
// 						adj[q.Elements[i]] = links
// 						cache[q.Elements[i]] = links
// 						rwmu.Unlock()
// 					} else {
// 						// rwmu.Lock()
// 						// rwmu.Lock()
// 						adj[q.Elements[i]] = cacheCheck
// 						// rwmu.Unlock()

// 						// mu.Unlock()
// 					}
// 					fmt.Println(q.Elements[i],"dist: ",dist[q.Elements[i]],"target URL: ", dist[targetURL])

// 				}
// 			})
// 		}
// 		wg.Wait()
// 		// length := q.GetLength()
// 		for i := 0; i < length; i++ {
// 			u := q.Elements[i]

// 			if dist[u] >= dist[targetURL] {
// 				continue
// 			}
// 			for _, v := range adj[u] {
// 				if v != startURL && dist[v] == 0 {
// 					dist[v] = maxInt
// 				}
// 			}

// 			for i := 0; i < len(adj[u]); i++ {
// 				// fmt.Println(adj[u][i], "TargetURL: ", targetURL, "Apakah sama: ", targetURL == adj[u][i])
// 				traversed++
// 				if dist[adj[u][i]] > dist[u]+1 {

// 					dist[adj[u][i]] = dist[u] + 1
// 					q.Enqueue(adj[u][i])
// 					parent[adj[u][i]] = nil
// 					parent[adj[u][i]] = append(parent[adj[u][i]], u)

// 				} else if dist[adj[u][i]] == dist[u]+1 {
// 					if !isSingle {
// 						parent[adj[u][i]] = append(parent[adj[u][i]], u)
// 					}
// 				}
// 				if dist[targetURL] == 1 {
// 					isFirst = true
// 					break
// 				}
// 				if isSingle && dist[targetURL] != maxInt {
// 					isFound = true
// 					break
// 				}
// 			}
// 			if isFirst {
// 				break
// 			}
// 			if isFound && isSingle {
// 				break
// 			}

// 		}
// 		if isFirst {
// 			break
// 		}
// 		if isFound && isSingle {
// 			break
// 		}

// 		q.Elements = q.Elements[length:]
// 		// wg.Wait()
// 	}

// 	paths := make([][]string, 0)
// 	path := make([]string, 0)

// 	paths = dfs(paths, path, parent, targetURL)

// 	for i := 0; i < len(paths); i++ {
// 		paths[i] = reverse(paths[i])
// 	}

// 	// updateMapInFile(cache, "./internal/controllers/cache/cache.json")
// 	updateMapInFile(cache, "../cache/cache.json")

// 	return paths, traversed

// }
