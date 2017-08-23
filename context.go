package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

type context struct{
	http.ServeMux
	clients map[string]*client
}

func newContext() *context {
	ctx := &context{
		ServeMux: *http.NewServeMux(),
		clients: map[string]*client{},
	}

	ctx.HandleFunc("/websocket", ctx.handleWebSocket)
	ctx.Handle("/", http.FileServer(http.Dir("./static/")))

	return ctx
}

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


func (ctx *context) handleWebSocket(w http.ResponseWriter, req *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, req, nil)
	if err != nil {
		// TODO better error handling
		panic(err)
	}

	client := newClient(ctx, conn)
	ctx.broadcast("* " + client.name + " joined", "#000000")
	for {
		msg, err := client.receive()
		if err != nil {
			// TODO better error handling
			panic(err)
		}
		err = ctx.broadcast(fmt.Sprintf("%s: %s", client.name, msg), client.color)
		if err != nil {
			panic(err)
		}
	}
}

func (ctx *context) broadcast(message string, color string) error {
	for _, c := range ctx.clients {
		err := c.send(message, color)
		if err != nil {
			// TODO continue and gather up errors
			return err
		}
	}

	return nil
}

