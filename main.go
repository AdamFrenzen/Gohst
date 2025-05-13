package main

import (
	"log"

	"github.com/adamfrenzen/gohst/internal/backend"
)

func main() {
	server := backend.NewServer()
	err := server.Start("localhost:64057")

	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
