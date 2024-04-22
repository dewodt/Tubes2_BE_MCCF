package controllers

import (
	"errors"
	"fmt"
	"math"
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

func dfs(paths *[][]string, path *[]string, parent map[string][]string, end string) {
	if parent[end] == nil {
		// path = append(path, end)
		*path = append(*path, end)
		// paths = append(paths, path)
		*paths = append(*paths, *path)
		// path = path[:len(path)-1]
		*path = (*path)[:len(*path)-1]
	} else {
		for i := 0; i < len(parent[end]); i++ {
			// path = append(path, end)
			*path = append(*path, end)
			dfs(paths, path, parent, parent[end][i])
			// path = path[:len(path)-1]
			*path = (*path)[:len(*path)-1]
		}
	}
}

func BFS(startURL string, targetURL string) ([][]string,[]string) {
	fmt.Println("Solving with BFS")
	fmt.Println("Start URL:", startURL)
	fmt.Println("Target URL:", targetURL)
	// var adj [][]int
	maxInt := math.MaxInt32
	adj := make(map[string][]string)
	parent := make(map[string][]string)
	parent[startURL] = nil
	q := Queue{Size: 0}
	dist := make(map[string]int)
	dist[startURL] = 0
	dist[targetURL] = maxInt
	q.Enqueue(startURL)
	//making bfs tree
	for !q.IsEmpty() {
		u, err := q.Peek()
		if err != nil {
			fmt.Println("Queue is empty")
		}
		q.Dequeue()
		if dist[u] >= dist[targetURL] {
			continue
		}
		links := getAllInternalLinks(startURL)
		for i := 0; i < len(links); i++ {
			adj[u] = append(adj[u], links[i])
		}
		for _, v := range links {
			if v != startURL && dist[v] == 0 {
				dist[v] = maxInt
			}
		}
		for i := 0; i < len(links); i++ {
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
		}
	}
	//change bfs tree to array of array of solution
	var paths [][]string
	var path []string

	dfs(&paths, &path, parent, targetURL)
	return paths,path
	//fill solution type with solution

	// Placeholder

}
