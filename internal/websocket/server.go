package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader   websocket.Upgrader
	mu         sync.Mutex
	activeConn *websocket.Conn
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/ws", s.handleWebSocket)

	log.Printf("WebSocket server listening on %s/ws\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.handleSingleConnection(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	s.mu.Lock()
	s.activeConn = conn
	s.mu.Unlock()

	defer s.closeConnection()
	s.readMessages()
}

func (s *Server) handleSingleConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeConn != nil {
		http.Error(w, "A client is already connected", http.StatusServiceUnavailable)
		log.Println("Rejected a new connection: already connected")
		return nil, fmt.Errorf("a client is already connected")
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Upgrade error:", err)
		return nil, err
	}

	return conn, nil
}

func (s *Server) closeConnection() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeConn != nil {
		s.activeConn.Close()
		s.activeConn = nil
		log.Println("Client disconnected")
	}
}

func (s *Server) readMessages() {
	for {
		s.mu.Lock()
		conn := s.activeConn
		s.mu.Unlock()

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
