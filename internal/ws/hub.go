package ws

import (
	"encoding/json"
	"fmt"
)

// room specific messages
type BroadcastMessage struct {
	Room string
	Data []byte
}

type Hub struct {
	// room -> clients
	Rooms       map[string]map[*Client]bool
	Broadcast   chan BroadcastMessage
	Register    chan *Client
	Unregister  chan *Client
	OnlineUsers map[string]bool
}

func NewHub() *Hub {
	return &Hub{
		Rooms:       make(map[string]map[*Client]bool),
		Broadcast:   make(chan BroadcastMessage),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		OnlineUsers: make(map[string]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		// client joins a room
		case client := <-h.Register:
			// create room if not present
			if _, ok := h.Rooms[client.Room]; !ok {
				h.Rooms[client.Room] = make(map[*Client]bool)
			}
			h.Rooms[client.Room][client] = true
			fmt.Printf("Client %s joined room: %s\n", client.Username, client.Room)

			h.OnlineUsers[client.Username] = true
			// broadcast online sts
			statusMsg, _ := json.Marshal(map[string]string{
				"type":   "status",
				"user":   client.Username,
				"status": "online",
				"room":   client.Room,
			})

			h.Broadcast <- BroadcastMessage{
				Room: client.Room,
				Data: statusMsg,
			}

		// client leaves a room
		case client := <-h.Unregister:
			if clients, ok := h.Rooms[client.Room]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					fmt.Printf("Client %s left room: %s\n", client.Username, client.Room)
					// clean up empty room
					if len(clients) == 0 {
						delete(h.Rooms, client.Room)
						fmt.Printf("Room %s deleted\n", client.Room)
					}
				}
			}

			delete(h.OnlineUsers, client.Username)
			// broadcast offline sts
			statusMsg, _ := json.Marshal(map[string]string{
				"type":   "status",
				"user":   client.Username,
				"status": "offline",
				"room":   client.Room,
			})

			h.Broadcast <- BroadcastMessage{
				Room: client.Room,
				Data: statusMsg,
			}

		case msg := <-h.Broadcast:
			if clients, ok := h.Rooms[msg.Room]; ok {
				for client := range clients {
					select {
					case client.Send <- msg.Data:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
		}
	}
}
