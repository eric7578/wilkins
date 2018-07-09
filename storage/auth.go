package storage

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Auth contains operations to check if session is allowed to access resources
type Auth struct {
	cli *redis.Client
}

func NewAuth() *Auth {
	return &Auth{
		cli: mustGetRedisClient(),
	}
}

func (a *Auth) CanSessionAccessChannel(sessionID, channelID string) error {
	sessChKey := getSessionChannelsKey(sessionID)
	sessChCmd := a.cli.SIsMember(sessChKey, channelID)

	if err := sessChCmd.Err(); err != nil {
		return newErrNotFound(fmt.Sprintf("session %s", sessionID))
	}

	if !sessChCmd.Val() {
		return newErrOperationNotAllowed()
	}

	return nil
}
