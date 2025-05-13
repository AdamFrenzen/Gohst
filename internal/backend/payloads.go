package backend

type ChatPayload struct {
	Prompt string `json:"prompt"`
}

type ChatResponsePayload struct {
	Response string `json:"response"`
}

type EchoPayload struct {
	Message string `json:"message"`
}
