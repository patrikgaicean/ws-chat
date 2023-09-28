package main

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

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send the message
				default:
					// failed to send
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
