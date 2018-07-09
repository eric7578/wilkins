package server

import (
	"testing"
	"time"

	"github.com/eric7578/wilkins/packet"
	"github.com/stretchr/testify/assert"
)

func Test_generateToken(t *testing.T) {
	sess := &packet.Session{
		ID:      "session-id",
		Created: time.Now().Unix(),
	}
	token, err := generateToken(sess)

	assert.NotEqual(t, "", token)
	assert.Nil(t, err)

	t.Run("validateToken", func(t *testing.T) {
		claims, err := validateToken(token)

		assert.Nil(t, err)
		assert.Nil(t, claims.Valid())
		assert.NotEqual(t, "", claims.Session.ID)
	})
}
