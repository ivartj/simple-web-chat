package main

import (
	"net/http"
	"github.com/gorilla/websocket"
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
	err = client.loop()
	if err != nil {
		// TODO better error handling
		panic(err)
	}
}

func (ctx *context) broadcast(msg *message) {

	for _, c := range ctx.clients {
		err := c.send(msg)
		if err != nil {
			log.Printf("WARNING: Error on sending message to '%s': %s\n", c.name, err.Error())
		}
	}
}

