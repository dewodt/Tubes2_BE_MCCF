package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func InitServer() *http.Server {
	// Get port from environment
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(),
		IdleTimeout:  1 * time.Hour,
		ReadTimeout:  0 * time.Second, // 0 or negative = no timeout
		WriteTimeout: 0 * time.Second, // 0 or negative = no timeout
	}

	return server
}
