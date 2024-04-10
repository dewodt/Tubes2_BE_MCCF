package utils

import "fmt"

// Print an array of strings.
func PrintArrayString(arr []string) {
	for _, s := range arr {
		fmt.Println(s)
	}
}
