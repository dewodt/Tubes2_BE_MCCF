package controllers

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"time"
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

func BFS(startURL string, targetURL string) ([][]string, []string) {
	fmt.Println("Solving with BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)
	runtime.GOMAXPROCS(runtime.NumCPU())
	// var adj [][]int
	if startURL == targetURL {
		return [][]string{{startURL}}, []string{startURL}
	}

	gm := NewGoRoutineManager(50)
	maxInt := math.MaxInt32
	adj := make(map[string][]string)
	parent := make(map[string][]string)
	parent[startURL] = nil
	q := Queue{Size: maxInt}
	dist := make(map[string]int)
	dist[startURL] = 0
	dist[targetURL] = maxInt
	q.Enqueue(startURL)
	//making bfs tree
	for !q.IsEmpty() {
		// u, err := q.Peek()
		// if err != nil {
		// 	fmt.Println("Queue is empty")
		// }
		// q.Dequeue()

		// fmt.Println(u,dist[u],"Target url: ",dist[targetURL])
		// if dist[u] >= dist[targetURL] {
		// 	// fmt.Println(u,dist[u],"skipped")
		// 	continue
		// )
		timenow4 := time.Now()
		timenow := time.Now()
		fmt.Println("len q", q.GetLength())

		for i := 0; i < q.GetLength(); i++ {
			i := i
			gm.Run(func() {
				mu.Lock()
				check := dist[q.Elements[i]] < dist[targetURL]
				mu.Unlock()
				if check {
					// fmt.Println(q.Elements[i], dist[q.Elements[i]])
					links := getAllInternalLinks(q.Elements[i])
					fmt.Println(q.Elements[i], "done")
					mu.Lock()
					adj[q.Elements[i]] = links
					mu.Unlock()
				}
			})
		}
		wg.Wait()

		fmt.Println("time elapsed1: ", time.Since(timenow))
		// perform multithreading on this
		length := q.GetLength()
		timenow3 := time.Now()
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

			var isFirst bool
			isFirst = false

			// timenow2 := time.Now()
			for i := 0; i < len(adj[u]); i++ {

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
				if dist[targetURL] == 1 {
					isFirst = true
					break
				}
				// fmt.Println(adj[u][i], dist[adj[u][i]], "Target url: ", dist[targetURL])
			}
			// fmt.Println( "time elapsed2: ",time.Since(timenow2))
			if isFirst {
				// fmt.Println()
				break
			}

		}
		fmt.Println("time elapsed3: ", time.Since(timenow3))
		fmt.Println("time elapsed4 : ", time.Since(timenow4))
		q.Elements = q.Elements[length:]
	}
	// fmt.Println(parent[targetURL],"debug")
	//change bfs tree to array of array of solution
	paths := make([][]string, 0)
	path := make([]string, 0)
	fmt.Println(len(parent))
	paths = dfs(paths, path, parent, targetURL)
	//perform multihtreading on this
	for i := 0; i < len(paths); i++ {
		paths[i] = reverse(paths[i])
	}
	return paths, path
	//fill solution type with solution
	// Placeholder
}
