package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	mu         sync.Mutex
	activeConn *websocket.Conn
)

// StartServer starts a WebSocket server on the given address.
func StartServer(addr string) {
	http.HandleFunc("/ws", handleWebSocket)

	log.Printf("WebSocket server listening on %s/ws\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	if activeConn != nil {
		mu.Unlock()
		http.Error(w, "A client is already connected", http.StatusServiceUnavailable)
		log.Println("Rejected a new connection: already connected")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		mu.Unlock()
		log.Println("Upgrade error:", err)
		return
	}
	activeConn = conn
	mu.Unlock()

	log.Println("Client connected")

	defer func() {
		conn.Close()
		log.Println("Client disconnected")
		mu.Lock()
		activeConn = nil
		mu.Unlock()
	}()

	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s\n", message)

		// Echo message back to client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
