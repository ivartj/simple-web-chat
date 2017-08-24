package main

import (
	"github.com/gorilla/websocket"
	"math/rand"
	"fmt"
	"encoding/json"
)

type client struct{
	conn *websocket.Conn
	name string
	color string
}

func generateClientName() string {
	r := rand.Uint32()
	m := r & 0xFF
	p1 := (r >> 8) % 3
	p2 := (p1 + ((r >> 16) % 2) + 1) % 3
	color := ((0xFF - m) << (p1 * 8)) | (m << (p2 * 8))
	name := fmt.Sprintf("#%06X", color)
	return name
}

func newClient(ctx *context, conn *websocket.Conn) *client {
	var name string
	for {
		// TODO ensure that this can not become an infinite loop (by too many users)
		name = generateClientName()
		_, exists := ctx.clients[name]
		if exists == false {
			break
		}
	}
	c := &client{
		conn: conn,
		name: name,
		color: name,
	}
	ctx.clients[name] = c

	return c
}

func (c *client) send(msg *message) error {
	err := c.conn.WriteMessage(websocket.TextMessage, msg.encode())
	return err
}

func (c *client) receive() (*message, error) {
	var msg message
	msgType, bytes, err := c.conn.ReadMessage()
	if websocket.IsUnexpectedCloseError(err) {
		return &message{
			Type: "leave",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	switch msgType {
	case websocket.TextMessage:
		err = json.Unmarshal(bytes, &msg)
		return &msg, err
	case websocket.CloseMessage:
		return &message{
			Type: "leave",
		}, nil
	default:
		return nil, fmt.Errorf("Unexpected WebSocket message type %d", msgType)
	}
}

