package main

import (
	"fmt"
	"tubes2-be-mccf/internal/server"
)

func main() {

	server := server.InitServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
