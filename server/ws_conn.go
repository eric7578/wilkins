package server

import (
	"encoding/json"
	"errors"
	"github.com/eric7578/wilkins/packet"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type wsconn struct {
	conn  *websocket.Conn
	send  chan *packet.WSMessage
	err   chan error
	recv  chan *packet.WSMessage
	close chan struct{}
}

func newWSConn(conn *websocket.Conn) *wsconn {
	return &wsconn{
		conn: conn,
		send: make(chan *packet.WSMessage),
		err:  make(chan error),
	}
}

func (c *wsconn) handleOutgoing() {
	defer func() {
		c.conn.Close()
	}()

	for {
		var err error
		select {
		case message := <-c.send:
			err = c.conn.WriteJSON(message)
		case err := <-c.err:
			message := packet.NewNewErrorMessage(err)
			err = c.conn.WriteJSON(message)
		}

		if err != nil {
			break
		}
	}
}

func (c *wsconn) handleIncoming() {
	defer func() {
		c.conn.Close()
	}()

	authenticated := false

	for {
		message := new(packet.WSMessage)
		err := c.conn.ReadJSON(message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(err)
			}
			break
		}

		log.Infof("receive message of %s", message.Type)

		if !authenticated {
			err := validateAuthenticatePacket(message)
			if err != nil {
				log.Error(err)
				break
			}

			log.Info("validate success")
			authenticated = true
		} else {
			switch message.Type {
			default:
				log.Warnf("type %s is not support", message.Type)
			}
		}
	}
}

func validateAuthenticatePacket(message *packet.WSMessage) error {
	if message.Type != packet.Authenticate {
		return errors.New("invalid authenticate message")
	}

	payload := new(packet.AuthenticatePayload)
	if err := json.Unmarshal(message.Payload, payload); err != nil {
		return errors.New("invalid authenticate message")
	}

	_, err := validateToken(payload.Token)
	return err
}
