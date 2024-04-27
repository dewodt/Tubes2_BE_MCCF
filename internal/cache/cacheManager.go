package cache

import (
	"encoding/json"
	"fmt"
	"os"
)

// Define global variable for storing cache
var Cache map[string][]string

func InitCache() {
	// Load the cache from the file
	fmt.Println("Reading cache, loading...")
	err := ReadMapFromFile("./cache/cache.json")

	if err != nil {
		fmt.Println("Error reading cache file")
		// Cache = make(map[string][]string)
		return
	}
	fmt.Println("Cache loaded")
	
}

func WriteMapToFile(filename string) error {
	// Convert the map to JSON
	jsonData, err := json.Marshal(Cache)
	if err != nil {
		return err
	}

	// Write the JSON to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadMapFromFile(filename string) error {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map
	err = json.Unmarshal(data, &Cache)
	if err != nil {
		return err
	}

	return nil
}

func UpdateMapInFile(filename string) error {
	// Read the existing data from the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map
	err = json.Unmarshal(data, &Cache)
	if err != nil {
		return err
	}

	// Merge the existing map with the new map
	for key, value := range Cache {
		if _, exists := Cache[key]; !exists {
			Cache[key] = value
		}
	}

	// Convert the merged map to JSON
	jsonData, err := json.Marshal(Cache)
	if err != nil {
		return err
	}

	// Write the JSON data back to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
