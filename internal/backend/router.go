package backend

import (
	"encoding/json"
	"log"
)

type Router struct{}

func New() *Router {
	return &Router{}
}

func (r *Router) RouteMessage(msgType string, payload json.RawMessage, server *WebSocketServer) {
	switch msgType {
	case "chat":
		var p ChatPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			log.Println("Chat payload error:", err)
			return
		}
		go handleChat(p.Prompt)

	case "echo":
		var p EchoPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			log.Println("Echo payload error:", err)
			return
		}
		go handleEcho()

	default:
		log.Println("Unknown message type:", msgType)
		return
	}
}
