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
	color := rand.Uint32() & 0x6F6F6F;
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

func (c *client) send(message string, color string) error {
	var msg struct{
		Message string `json:"message"`
		Color string `json:"color"`
	}
	msg.Message = message
	msg.Color = color
	bytes, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, bytes)
	return err
}

func (c *client) receive() (string, error) {
	var msgobj struct{
		Message string `json:"message"`
	}
	// TODO handle WebSocket message types other than TextMessage
	_, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bytes, &msgobj)
	if err != nil {
		return "", err
	}
	return msgobj.Message, nil
}

