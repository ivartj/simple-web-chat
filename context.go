package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"log"
)

type context struct{
	http.ServeMux
	clients map[string]*client
}

func newContext(assets string) *context {
	ctx := &context{
		ServeMux: *http.NewServeMux(),
		clients: map[string]*client{},
	}

	ctx.HandleFunc("/websocket", ctx.handleWebSocket)
	ctx.Handle("/", http.FileServer(http.Dir(assets + "/static/")))

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
	ctx.broadcast(&message{
		Type: "join",
		User: client.name,
	})
	for {
		msg, err := client.receive()
		if err != nil {
			// TODO better error handling
			panic(err)
		}
		switch msg.Type {
		case "message":
			err = ctx.broadcast(&message{
				Type: "message",
				User: client.name,
				Color: client.color,
				Text: msg.Text,
			})
			if err != nil {
				panic(err)
			}
		case "leave":
			delete(ctx.clients, client.name)
			err = ctx.broadcast(&message{
				Type: "leave",
				User: client.name,
			})
			return
		default:
			panic(fmt.Errorf("Unexpected message type '%s'", msg.Type))
		}
	}
}

func (ctx *context) broadcast(msg *message) error {

	for _, c := range ctx.clients {
		err := c.send(msg)
		if err != nil {
			log.Printf("WARNING: Error on sending message to '%s': %s", c.name, err.Error())
		}
	}

	return nil
}

