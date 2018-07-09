package storage

import (
	"fmt"
	"strconv"

	"github.com/eric7578/wilkins/packet"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	// Picked status means the channel has already been picked by other sessions
	Picked int = iota
	// Orphan means the channel is not picked
	Orphan
)

// Channels contains the methods that manipulates channel data
type Channels struct {
	cli *redis.Client
}

// NewChannels returns new Channels instance
func NewChannels() *Channels {
	return &Channels{
		cli: mustGetRedisClient(),
	}
}

// List list all channels the session had
func (c *Channels) List(sessionID string) []string {
	key := getSessionChannelsKey(sessionID)
	return c.cli.ZRevRange(key, 0, -1).Val()
}

// Create create channel only if there is not orphan channel belong to the session
func (c *Channels) Create(sessionID string) *packet.Channel {
	channelID := uuid.Must(uuid.NewV4()).String()
	sessChKey := getSessionChannelsKey(sessionID)
	chSessKey := getChannelSessionsKey(channelID)
	statusKey := getChannelStatusKey(channelID)

	p := c.cli.TxPipeline()
	p.SAdd(sessChKey, channelID)
	p.SAdd(chSessKey, sessionID)
	p.Set(statusKey, Orphan, 0)
	_, err := p.Exec()
	if err != nil {
		panic(err)
	}

	ch := &packet.Channel{
		ID:       channelID,
		Status:   Orphan,
		Sessions: []string{sessionID},
	}
	return ch
}

// HasOrphan will check if session contains any orphan channel
func (c *Channels) HasOrphan(sessionID string) bool {
	sessChKey := getSessionChannelsKey(sessionID)

	channelIDs := c.cli.SMembers(sessChKey).Val()
	if len(channelIDs) == 0 {
		return false
	}

	// check if orphan channel is included
	statusKeys := getChannelStatusKeys(channelIDs...)
	cmd := c.cli.MGet(statusKeys...)
	if err := cmd.Err(); err != nil {
		panic(err)
	}

	for _, val := range cmd.Val() {
		intVal, _ := strconv.Atoi(val.(string))
		if intVal == Orphan {
			return true
		}
	}
	return false
}

// InChannel check if session has a channel
func (c *Channels) InChannel(sessionID string, channelID string) bool {
	key := getSessionChannelsKey(sessionID)
	cmd := c.cli.SIsMember(key, channelID)
	return cmd.Val()
}

// Get returns channel instance by channelID
func (c *Channels) Get(channelID string) (*packet.Channel, error) {
	// get Status
	statusKey := getChannelStatusKey(channelID)
	statusCmd := c.cli.Get(statusKey)
	if err := statusCmd.Err(); err != nil {
		return nil, newErrNotFound(channelID)
	}

	status, err := strconv.Atoi(statusCmd.Val())
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("invalid status: %v", status)
	}

	// get sessions
	chSessKey := getChannelSessionsKey(channelID)
	chSessCmd := c.cli.SMembers(chSessKey)
	if err := chSessCmd.Err(); err != nil {
		log.Error(err)
		return nil, fmt.Errorf("invalid sessions")
	}

	ch := &packet.Channel{
		ID:       channelID,
		Status:   status,
		Sessions: chSessCmd.Val(),
	}
	return ch, nil
}

// PostTo push a new message to channel
func (c *Channels) PostTo(sessionID string, channelID string, message string) error {
	msgKey := getMessagesKey(channelID)
	msgCmd := c.cli.RPush(msgKey, message)
	return msgCmd.Err()
}

// FindBySession returns channels where the session belongs
func (c *Channels) FindBySession(sessionID string) ([]*packet.Channel, error) {
	sessChKey := getSessionChannelsKey(sessionID)
	sessChCmd := c.cli.SMembers(sessChKey)
	channelIDs := sessChCmd.Val()
	channels := make([]*packet.Channel, len(channelIDs))
	for idx, channelID := range channelIDs {
		ch, err := c.Get(channelID)
		if err != nil {
			return nil, err
		}
		channels[idx] = ch
	}
	return channels, nil
}

// Join add a session info channel
func (c *Channels) Join(sessionID string, channelID string) error {
	statusKey := getChannelStatusKey(channelID)
	sessChKey := getSessionChannelsKey(sessionID)
	chSessKey := getChannelSessionsKey(channelID)

	statusCmd := c.cli.Get(statusKey)
	if err := statusCmd.Err(); err != nil {
		return newErrNotFound(channelID)
	}

	p := c.cli.TxPipeline()
	p.SAdd(sessChKey, channelID)
	p.SAdd(chSessKey, sessionID)
	_, err := p.Exec()
	if err != nil {
		return err
	}

	return nil
}

// ReadMessages returns all messages and sessions in channel
func (c *Channels) ReadMessages(sessionID, channelID string, offset int64) ([]string, error) {
	messageKey := getMessagesKey(channelID)
	messageCmd := c.cli.LRange(messageKey, 0, offset)
	if err := messageCmd.Err(); err != nil {
		return nil, err
	}

	messageLenCmd := c.cli.LLen(messageKey)
	if err := messageLenCmd.Err(); err != nil {
		return nil, err
	}

	readKey := getMessagesReadKey(channelID, sessionID)
	c.cli.HSet(readKey, sessionID, messageLenCmd.Val())
	return messageCmd.Val(), nil
}
