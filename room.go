package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

func (rm *room) run() {
	for {
		select {
		case client := <-rm.join:
			// joining
			rm.clients[client] = true
		case client := <-rm.leave:
			// leaving
			delete(rm.clients, client)
			close(client.send)
		case msg := <-rm.forward:
			// forward message to all clients
			for client := range rm.clients {
				select {
				case client.send <- msg:
					// send the message
				default:
					// failed to send
					delete(rm.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (rm *room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   rm,
	}
	rm.join <- client

	defer func() {
		rm.leave <- client
	}()

	go client.write()
	client.read()
}
