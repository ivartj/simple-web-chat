package main

import (
	"github.com/gorilla/websocket"
	"math/rand"
	"fmt"
	"encoding/json"
	"sync"
	"time"
)

type client struct{
	ctx *context
	conn *websocket.Conn
	name string
	color string
	readLock sync.Mutex
	writeLock sync.Mutex
	lastPongReceived time.Time
}

// Generates color code with a value of 255 distributed randomly across two primary colours.
func generateColorCode() string {
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
		name = generateColorCode()
		_, exists := ctx.clients[name]
		if exists == false {
			break
		}
	}
	c := &client{
		ctx: ctx,
		conn: conn,
		name: name,
		color: name,
	}
	ctx.clients[name] = c
	c.conn.SetPongHandler(c.pongHandle)

	return c
}

func (c *client) send(msg *message) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	err := c.conn.WriteMessage(websocket.TextMessage, msg.encode())
	return err
}

func (c *client) receive() (*message, error) {
	var msg message

	c.readLock.Lock()
	defer c.readLock.Unlock()

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

func (c *client) ping() error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second * 5))
}

func (c *client) pongHandle(appData string) error {
	c.lastPongReceived = time.Now()
	return nil
}

const pingInterval = 20

func (c *client) loop() error {
	go func() {
		for {
			pingSent := time.Now()
			err := c.ping()
			if err != nil {
				c.conn.Close()
				break
			}
			time.Sleep(time.Second * pingInterval)
			if c.lastPongReceived.IsZero() || c.lastPongReceived.Before(pingSent) {
				c.send(&message{
					Type: "server-system",
					Text: "Closing connection due to ping timeout",
				})
				c.conn.Close()
			}
		}
	}()

	defer func() {
		delete(c.ctx.clients, c.name)
		c.ctx.broadcast(&message{
			Type: "leave",
			User: c.name,
		})
	}()

	for {
		msg, err := c.receive()
		if err != nil {
			return err
		}
		switch msg.Type {
		case "message":
			c.ctx.broadcast(&message{
				Type: "message",
				User: c.name,
				Color: c.color,
				Text: msg.Text,
			})
		case "leave":
			// Handled in deferred function above
			return nil
		default:
			return fmt.Errorf("Unexpected message type '%s'", msg.Type)
		}
	}

	return nil
}

