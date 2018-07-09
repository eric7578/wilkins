package storage

import (
	"bytes"
)

// sess:{sessionID}
func getSessionKey(sessionID string) string {
	var buf bytes.Buffer
	buf.WriteString("sess:")
	buf.WriteString(sessionID)
	return buf.String()
}

// ch:{channelID}:msg
// List
func getMessagesKey(channelID string) string {
	var buf bytes.Buffer
	buf.WriteString("ch:")
	buf.WriteString(channelID)
	buf.WriteString(":msg")
	return buf.String()
}

// ch:{channelID}:read
// Map [sessionID] {readIndex}
func getMessagesReadKey(channelID, sessionID string) string {
	var buf bytes.Buffer
	buf.WriteString("ch:")
	buf.WriteString(channelID)
	buf.WriteString(":read")
	return buf.String()
}

// ch:{channelID}:st
// KV
func getChannelStatusKeys(channelIDs ...string) []string {
	keys := make([]string, len(channelIDs))
	var buf bytes.Buffer
	for i, channelID := range channelIDs {
		buf.WriteString("ch:")
		buf.WriteString(channelID)
		buf.WriteString(":st")
		keys[i] = buf.String()
		buf.Reset()
	}
	return keys
}

func getChannelStatusKey(channelID string) string {
	var buf bytes.Buffer
	buf.WriteString("ch:")
	buf.WriteString(channelID)
	buf.WriteString(":st")
	return buf.String()
}

// ch:{channelID}:sess
// Set, sessionID
func getChannelSessionsKey(channelID string) string {
	var buf bytes.Buffer
	buf.WriteString("ch:")
	buf.WriteString(channelID)
	buf.WriteString(":sess")
	return buf.String()
}

// sess:{sessionID}:ch
// Set, channelID
func getSessionChannelsKey(sessionID string) string {
	var buf bytes.Buffer
	buf.WriteString("sess:")
	buf.WriteString(sessionID)
	buf.WriteString(":ch")
	return buf.String()
}
