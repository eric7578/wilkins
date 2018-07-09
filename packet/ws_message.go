package packet

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

// WSMessage is the basic struct of ws packet
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	// Authenticate is the type of AuthenticatePayload
	Authenticate string = "authenticate"
	// Error is the type of ErrorPayload
	Error string = "error"
)

// AuthenticatePayload should be sent immediately once a websocket connection is connected
type AuthenticatePayload struct {
	Token string `json:"token"`
}

// ErrorPayload represent general errors
type ErrorPayload struct {
	ID    string `json:"ID"`
	Error string `json:"error"`
}

// NewNewErrorMessage ...
func NewNewErrorMessage(err error) *WSMessage {
	payload := ErrorPayload{
		ID:    uuid.Must(uuid.NewV4()).String(),
		Error: err.Error(),
	}
	bytes, err := json.Marshal(payload)

	if err != nil {
		panic(err)
	}

	return &WSMessage{
		Type:    Error,
		Payload: bytes,
	}
}
