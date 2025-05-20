package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WebSocketServer struct {
	upgrader   websocket.Upgrader
	mu         sync.Mutex
	activeConn *websocket.Conn
	router     *Router
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		router: NewRouter(),
	}
}

func (s *WebSocketServer) Start(addr string) error {
	http.HandleFunc("/ws", s.handleWebSocket)

	log.Printf("WebSocket server listening on %s/ws\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.handleSingleConnection(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client connected")

	defer s.closeConnection()
	s.readMessages(conn)
}

func (s *WebSocketServer) handleSingleConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeConn != nil {
		http.Error(w, "A client is already connected", http.StatusServiceUnavailable)
		return nil, fmt.Errorf("Rejected a new connection: already connected")
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Upgrade error:", err)
		return nil, err
	}

	s.activeConn = conn
	return s.activeConn, nil
}

func (s *WebSocketServer) closeConnection() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeConn != nil {
		s.activeConn.Close()
		s.activeConn = nil
		log.Println("Client disconnected")
	}
}

func (s *WebSocketServer) readMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		go s.decodeAndRoute(message)
	}
}

func (s *WebSocketServer) decodeAndRoute(message []byte) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Println("JSON unmarshal error:", err)
		return
	}

	s.router.RouteMessage(msg.Type, msg.Payload, s)
}

func (s *WebSocketServer) SendMessage(payload any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.activeConn.WriteJSON(payload)
}
