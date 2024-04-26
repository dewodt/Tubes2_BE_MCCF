package db

import (
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

// Global session variable
var Session *gocql.Session

// Initialize cassadnra database
func InitDB() {
	// Create cluster config
	var err error
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 1 * time.Hour

	// Connect to cassandra
	sessionInitKeySpace, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Error connecting to cassandra:", err.Error())
	}

	// Create keyspace if not exists
	err = sessionInitKeySpace.Query(`
		CREATE KEYSPACE IF NOT EXISTS wikirace
		WITH REPLICATION = {
			'class' : 'SimpleStrategy',
			'replication_factor' : 1
		}`).Exec()
	if err != nil {
		log.Fatal("Error creating keyspace:", err.Error())
	}

	// Close session
	sessionInitKeySpace.Close()

	// Initialize global valid session
	cluster.Keyspace = "wikirace"
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Error connecting to cassandra:", err.Error())
	}

	// Create table if not exists
	/*
		Data Structure:
		{
			url1: [containedUrl1, containedUrl2, containedUrl3, ..., containedUrlN],
			url2: [containedUrl1, containedUrl2, containedUrl3, ..., containedUrlM],
			...
		}
	*/
	err = Session.Query(`
    	CREATE TABLE IF NOT EXISTS wikipedia_cache (
        	url TEXT PRIMARY KEY,
        	internal_urls LIST<TEXT>
    	)`).Exec()
	if err != nil {
		log.Fatal("Error creating table:", err.Error())
	}
}
