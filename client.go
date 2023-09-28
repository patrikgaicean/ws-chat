package main

import (
	"github.com/gorilla/websocket"
)

// a single chatting user
type client struct {
	// websocket for current client
	socket *websocket.Conn
	// channel on which messages are sent
	send chan []byte
	// room this client is chatting in
	room *room
}

type room struct {
	// chan that holds incoming messages that should be forwarded to other
	// clients
	forward chan []byte
}
