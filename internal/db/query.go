package db

import (
	"github.com/gocql/gocql"
)

// Find an URL in the database
func FindURL(url string) ([]string, error) {
	// Query the URL
	var internalURLs []string
	err := Session.Query(`SELECT internal_urls FROM wikipedia_cache WHERE url = ?`, url).Scan(&internalURLs)
	if err != nil {
		return nil, err
	}

	return internalURLs, nil
}

// Insert many URLS to the database to update the cache
func InsertURLs(urls map[string][]string) error {
	// Insert all URLs
	batch := Session.NewBatch(gocql.LoggedBatch)

	statement := `INSERT INTO wikipedia_cache (url, internal_urls) VALUES (?, ?) IF NOT EXISTS` // If not exists to prevent overwriting

	i := 0
	for url, internalURLs := range urls {
		batch.Query(statement, url, internalURLs)

		// Execute batch every 50 queries
		if i%50 == 0 {
			if err := Session.ExecuteBatch(batch); err != nil {
				return err
			}
			// Reset batch for the next set of queries
			batch = Session.NewBatch(gocql.LoggedBatch)
		}

		i++
	}

	// Execute any remaining queries in the batch
	if err := Session.ExecuteBatch(batch); err != nil {
		return err
	}

	// no error
	return nil
}

// Select all URLs from the database
// ONLY FOR TESTING PURPOSES
func SelectAllURLs() (map[string][]string, error) {
	// Query all URLs
	iter := Session.Query(`SELECT url, internal_urls FROM wikipedia_cache`).Iter()
	urls := make(map[string][]string)

	// Scan all URLs
	var url string
	var internalURLs []string
	for iter.Scan(&url, &internalURLs) {
		urls[url] = internalURLs
	}

	// Return error if any
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return urls, nil
}
