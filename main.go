package main

import (
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	server := NewServer()
	server.Run(":" + port)
	fmt.Printf("Running server on port %v\n", port)
}
