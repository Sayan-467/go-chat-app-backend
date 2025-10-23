package websocket

import "github.com/gorilla/websocket"

// manage all connected clients and message broadcasting
type Client struct {
	Hub 	*Hub
	Conn	*websocket.Conn
	Send 	chan []byte 
}

type Hub struct {
	Clients 	map[*Client]bool
	Broadcast 	chan []byte
	Register 	chan *Client
	Unregister 	chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h* Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}