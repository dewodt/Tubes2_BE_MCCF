package utils

import "fmt"

// Print an array of strings.
func PrintResultPath(resultPath [][]string) {
	fmt.Println("Found", len(resultPath), "paths : ")
	for i, path := range resultPath {
		fmt.Print("Path ", i+1, " : ")
		fmt.Println(path)
	}
}
