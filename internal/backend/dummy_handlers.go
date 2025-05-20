package backend

import (
	"fmt"
)

func handleChat(prompt string, ws *WebSocketServer) {
	fmt.Println("response to", prompt)
	ws.SendMessage(ChatResponsePayload{Response: "HEY"})
}

func handleEcho() {
	fmt.Println("handled echo")
}
