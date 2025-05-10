package main

import (
	// "fmt"

	"github.com/adamfrenzen/gohst/internal/websocket"
)

func main() {
	websocket.StartServer("localhost:64057")
}
