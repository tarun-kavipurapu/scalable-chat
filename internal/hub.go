package internal

import (
	"fmt"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client

	broadcast chan *Message
	//used to send messages to the client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userId] = client
			// Optional: Log client registration
			fmt.Printf("Registered client: %s\n", client.userId)

		case client := <-h.unregister:
			if _, ok := h.clients[client.userId]; ok {
				delete(h.clients, client.userId)
				close(client.sendTo)
				// Optional: Log client unregistration
				log.Printf("Unregistered client: %s\n", client.userId)
			}

		case message := <-h.broadcast:
			// Find the recipient client
			client, ok := h.clients[message.To]
			if ok {
				select {
				case client.sendTo <- message:
					log.Printf("Sent message to client: %s\n", message.To)
				default:
					// If the client's sendTo channel is blocked
					close(client.sendTo)
					delete(h.clients, client.userId)
					log.Printf("Closed channel and removed client due to blocked sendTo: %s\n", message.To)
				}
			} else {
				log.Printf("Attempted to send message to non-existent client: %s\n", message.To)
			}
		}
	}
}
