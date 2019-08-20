package ws

import (
	"github.com/sirupsen/logrus"
)

// Hub represents a controller for handling all WebSocket communications.
type Hub struct {
	Broadcast  chan []byte
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
}

// CreateHub returns a new instance of Hub.
func CreateHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Execute runs the main loop for handling WebSocket Hub events.
func (h *Hub) Execute() {
	logrus.Debug("Starting WS Hub executor.")
	for {
		select {
		case client := <-h.Register:
			logrus.Debug("Received Register message on WS Hub.")
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			logrus.Debug("Received Broadcast message on WS Hub.")
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
