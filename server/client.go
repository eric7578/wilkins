package server

import (
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

// Client represet a connection from client
type Client interface {
	ID() string
	Receive(message interface{}) error
	Send(message interface{}) error
}

// NewWebSocketClient get Client wrapped with websocket connection
func NewWebSocketClient(conn *websocket.Conn) Client {
	return &websocketClient{
		Conn: conn,
		id:   uuid.Must(uuid.NewV4()).String(),
	}
}

type websocketClient struct {
	*websocket.Conn
	id string
}

func (c *websocketClient) ID() string {
	return c.id
}

func (c *websocketClient) Send(message interface{}) error {
	return c.Conn.WriteJSON(message)
}

func (c *websocketClient) Receive(message interface{}) error {
	return c.Conn.ReadJSON(message)
}
