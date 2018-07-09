package packet

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Channel represent channel instance
type Channel struct {
	ID       string   `json:"id"`
	Status   int      `json:"status"`
	Sessions []string `json:"sessions"`
}

// ChannelInfo contains the complete info of a channel
type ChannelInfo struct {
	ID       string     `json:"id"`
	Status   int        `json:"status"`
	Sessions []*Session `json:"sessions"`
}

// NewChannelInfo return new ChannelInfo instance
func NewChannelInfo(channel *Channel) *ChannelInfo {
	return &ChannelInfo{
		ID:       channel.ID,
		Status:   channel.Status,
		Sessions: make([]*Session, len(channel.Sessions)),
	}
}

// Session represet a session instance
type Session struct {
	ID      string `json:"id"`
	Profile string `json:"profile,omitempty"`
	Created int64  `json:"created"`
}

// MessagesBody is the POST body of messages
type MessagesBody struct {
	Messages string `json:"messages"`
}

// Message represent a single message sent from user
type Message struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

// IsValidMessageString check if message string is valid
func IsValidMessageString(messageString string) error {
	messages := make([]Message, 0)
	err := json.Unmarshal([]byte(messageString), &messages)
	if err != nil {
		return errors.New("string should be array of messages")
	}

	for idx, message := range messages {
		var err error
		switch message.Type {
		case "text":
			var textContent string
			err = json.Unmarshal(message.Content, &textContent)
		default:
			return fmt.Errorf("invalid message type: %s", message.Type)
		}

		if err != nil {
			return fmt.Errorf("invalid %s content at %v", message.Type, idx)
		}
	}

	return nil
}
