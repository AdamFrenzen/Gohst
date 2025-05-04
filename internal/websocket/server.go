package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	conn      *websocket.Conn
	connMutex sync.Mutex
	sendChan  chan []byte
}

func NewServer() *Server {
	return &Server{
		sendChan: make(chan []byte),
	}
}
