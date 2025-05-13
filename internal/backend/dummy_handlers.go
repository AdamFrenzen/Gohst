package backend

import (
	"fmt"
)

func handleChat(prompt string) {
	fmt.Println("response to", prompt)
}

func handleEcho() {
	fmt.Println("handled echo")
}
