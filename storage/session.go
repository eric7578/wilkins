package storage

import (
	"strconv"
	"time"

	"github.com/eric7578/wilkins/packet"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
)

// Sessions contains the session data
type Sessions struct {
	cli *redis.Client
}

// NewSessions return storage instance of Sessions
func NewSessions() *Sessions {
	return &Sessions{
		cli: mustGetRedisClient(),
	}
}

// Get returns the session instance by sessionID
func (s *Sessions) Get(sessionID string) (*packet.Session, error) {
	key := getSessionKey(sessionID)
	sessResult := s.cli.HGetAll(key)
	if sessResult == nil {
		return nil, newErrNotFound(key)
	}

	m := sessResult.Val()
	created, _ := strconv.Atoi(m["Created"])
	sess := &packet.Session{
		ID:      sessionID,
		Profile: m["Profile"],
		Created: int64(created),
	}

	return sess, nil
}

// Create returns a new session instance
func (s *Sessions) Create() *packet.Session {
	sess := &packet.Session{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Created: time.Now().Unix(),
	}

	m := map[string]interface{}{
		"ID":      sess.ID,
		"Created": sess.Created,
	}
	err := s.cli.HMSet(getSessionKey(sess.ID), m).Err()
	if err != nil {
		panic(err)
	}

	return sess
}
