package controllers

import (
	"encoding/json"
	"io/ioutil"
)

func writeMapToFile(m map[string][]string, filename string) error {
	// Convert the map to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Write the JSON to the file
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readMapFromFile(filename string) (map[string][]string, error) {
	// Read the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a map
	var m map[string][]string
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func updateMapInFile(m map[string][]string, filename string) error {
	// Read the existing data from the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map
	var existingMap map[string][]string
	err = json.Unmarshal(data, &existingMap)
	if err != nil {
		return err
	}

	// Merge the existing map with the new map
	for key, value := range m {
		if _, exists := existingMap[key]; !exists {
			existingMap[key] = value
		}
	}

	// Convert the merged map to JSON
	jsonData, err := json.Marshal(existingMap)
	if err != nil {
		return err
	}

	// Write the JSON data back to the file
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
