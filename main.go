package main

import (
	// "fmt"

	"log"

	"github.com/adamfrenzen/gohst/internal/websocket"
)

func main() {
	server := websocket.NewServer()
	err := server.Start("localhost:64057")

	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
