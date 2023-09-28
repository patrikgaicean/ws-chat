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

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err != nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

type room struct {
	// chan that holds incoming messages that should be forwarded to other
	// clients
	forward chan []byte
	// chan for clients wishing to join the room
	join chan *client
	// chan for clients wishing to leave the room
	leave chan *client
	// holds all current clients in the room
	clients map[*client]bool
}
