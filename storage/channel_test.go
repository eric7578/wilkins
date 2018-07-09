package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChannels_Create(t *testing.T) {
	channels := NewChannels()
	sessionID := "sessionID"
	ch := channels.Create(sessionID)

	assert.Equal(t, sessionID, ch.Sessions[0])
	assert.Equal(t, Orphan, ch.Status)

	t.Run("InChannel", func(t *testing.T) {
		inChannel := channels.InChannel(sessionID, ch.ID)

		assert.True(t, inChannel)
	})

	t.Run("Get", func(t *testing.T) {
		found, err := channels.Get(ch.ID)

		assert.Nil(t, err)
		assert.Equal(t, ch.ID, found.ID)
		assert.Equal(t, sessionID, found.Sessions[0])
	})
}

func TestChannels_PostTo(t *testing.T) {
	channels := NewChannels()
	sessionID := "sessionID"
	ch := channels.Create(sessionID)

	errors := []error{
		channels.PostTo(sessionID, ch.ID, "msg1"),
		channels.PostTo(sessionID, ch.ID, "msg2"),
		channels.PostTo(sessionID, ch.ID, "msg3"),
	}

	for _, err := range errors {
		assert.Nil(t, err)
	}

	t.Run("ReadMessages should return message by offset from left", func(t *testing.T) {
		msgs, err := channels.ReadMessages(sessionID, ch.ID, 1)

		assert.Nil(t, err)
		assert.True(t, assert.ObjectsAreEqual([]string{"msg1", "msg2"}, msgs))
	})
}
