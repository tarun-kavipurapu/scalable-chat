package internal

import (
	"context"
	"fmt"
	"log"
	db "tarun-kavipurapu/test-go-chat/db/sqlc"
	"tarun-kavipurapu/test-go-chat/internal/handlers"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients     map[int64]*Client
	register    chan *Client
	unregister  chan *Client
	broadcast   chan *Message
	chatHandler *handlers.ChatHandler
	store       db.Store
}

// NewHub initializes and returns a new Hub instance.
func NewHub(chatHandler *handlers.ChatHandler, store db.Store) *Hub {
	return &Hub{
		broadcast:   make(chan *Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[int64]*Client),
		chatHandler: chatHandler,

		store: store,
	}
}

// Run handles the registration, unregistration, and message broadcasting to clients.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.HandleUserRegisterEvent(client, ctx)

		case client := <-h.unregister:
			h.HandleUserDisconnectEvent(client)

		case message := <-h.broadcast:
			h.HandleMessageBroadcast(message, ctx)
		}
	}
}

// HandleUserRegisterEvent handles the user registration event.
func (h *Hub) HandleUserRegisterEvent(client *Client, ctx context.Context) {

	h.clients[client.userId] = client
	// Log client registration
	fmt.Printf("Registered client: %d\n", client.userId)
}

// HandleUserDisconnectEvent handles the user disconnection event.
func (h *Hub) HandleUserDisconnectEvent(client *Client) {
	if _, ok := h.clients[client.userId]; ok {
		delete(h.clients, client.userId)
		close(client.sendTo)
		// Log client unregistration
		log.Printf("Unregistered client: %d\n", client.userId)
	}
}

// HandleMessageBroadcast handles the message broadcasting.
func (h *Hub) HandleMessageBroadcast(message *Message, ctx context.Context) {
	//before uploading to databse if the User with from and to UserId Exists
	// h.store.GetUserById(message.ID)
	//first upload the message to the database

	//in future try to send the error to client wither through a websocket message or http JSOn response
	_, err := h.store.GetUserById(ctx, int64(message.From))
	if err != nil {
		log.Println("From userId is not available in Database Please Register")
		return
	}
	_, err = h.store.GetUserById(ctx, int64(message.To))
	if err != nil {
		log.Println("To userId is not available in Database Please Register")
		return
	}
	//Insert into the Databsae
	err = h.chatHandler.InsertMessage(ctx, message.From, message.To, message.Content)
	if err != nil {
		log.Println("Unable to Insert Messaages In the Database")
		return

	}
	client, ok := h.clients[message.To]
	if ok {
		select {
		case client.sendTo <- message:
			log.Printf("Sent message to client: %d\n", message.To)
		default:
			// If the client's sendTo channel is blocked
			close(client.sendTo)
			delete(h.clients, client.userId)
			log.Printf("Closed channel and removed client due to blocked sendTo: %d\n", message.To)
		}
	} else {
		//In this Section IT should Publish to the redis CLient of
		log.Printf("Attempted to send message to non-existent client: %d\n", message.To)
	}
}
