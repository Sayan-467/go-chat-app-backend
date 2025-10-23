package ws

import (
	"chat-app-backend/internal/config"
	"chat-app-backend/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// manages each websocket's client connection and communicates with hub
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Id       string
	Username string
	Room     string
	Receiver string
	Conn     *websocket.Conn
	Hub      *Hub
	Send     chan []byte
}

type Message struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

// ServeWS upgrades the HTTP request to a WebSocket and registers the client
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error: ", err)
		return
	}

	// get username from query param
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}
	
	// get room from query param
	room := r.URL.Query().Get("room")
	if room == "" {
		room = "general"
	}
	receiver := r.URL.Query().Get("receiver")

	// determine room type (for dm)
	if receiver != "" {
		// dm:alice:bob
		if username < receiver {
			room = fmt.Sprintf("dm:%s:%s", username, receiver)
		} else {
			room = fmt.Sprintf("dm:%s:%s", receiver, username)
		}
	}

	log.Printf("[CONNECT] RemoteAddr=%s, username=%s, receiver=%s, room=%s\n",
		r.RemoteAddr, username, receiver, room)

	client := &Client{
		Id:       r.RemoteAddr,
		Username: username,
		Room:     room,
		Receiver: receiver,
		Conn:     conn,
		Hub:      hub,
		Send:     make(chan []byte),
	}

	log.Printf("[REGISTER] Registering client %s in room %s", username, room)

	// register client in hub
	hub.Register <- client

	// Immediately send "online" status to the room
	statusMsg, _ := json.Marshal(map[string]string{
		"type":   "status",
		"user":   client.Username,
		"status": "online",
		"room":   client.Room,
	})
	hub.Broadcast <- BroadcastMessage{
		Room: client.Room,
		Data: statusMsg,
	}

	// send last messages as chat history
	var messages []models.Message
	if err := config.DB.
		Where("room = ?", client.Room).
		Order("createdAt desc").Limit(20).Find(&messages).Error; err == nil {
		// send oldest msg first
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}

		for _, m := range messages {
			msgjson, _ := json.Marshal(map[string]string{
				"type":      "message",
				"sender":    m.Sender,
				"receiver":  m.Receiver,
				"room":      m.Room,
				"message":   m.Content,
				"timestamp": m.CreatedAt.Format("15:04"),
			})
			client.Send <- msgjson
		}
	} else {
		log.Println("error fetching chat history", err)
	}

	// start read and write pump
	go client.readPump()
	go client.writePump()
}

// readPump reads messages from this client and broadcasts them
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c

		// Send offline status when client disconnects
		statusMsg, _ := json.Marshal(map[string]string{
			"type":   "status",
			"user":   c.Username,
			"status": "offline",
			"room":   c.Room,
		})
		c.Hub.Broadcast <- BroadcastMessage{
			Room: c.Room,
			Data: statusMsg,
		}

		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			log.Printf("[DISCONNECT] %s encountered read error: %v", c.Username, err)
			break
		}

		// typing indicator
		var incoming map[string]interface{}
		if err := json.Unmarshal(msg, &incoming); err == nil {
			if incoming["type"] == "typing" {
				c.Hub.Broadcast <- BroadcastMessage{
					Room: c.Room,
					Data: msg,
				}
				continue
			}
		}

		message := models.Message{
			Sender:   c.Username,
			Receiver: c.Receiver,
			Room:     c.Room,
			Content:  string(msg),
		}

		// broadcast to all connected clients
		jsonMsg, _ := json.Marshal(map[string]string{
			"type":      "message",
			"sender":    message.Sender,
			"receiver":  message.Receiver,
			"room":      message.Room,
			"message":   message.Content,
			"timestamp": time.Now().Format("15:04"),
		})

		log.Printf("[BROADCAST] %s -> room %s: %s", c.Username, c.Room, string(jsonMsg))
		// broadcase with room infos
		c.Hub.Broadcast <- BroadcastMessage{
			Room: c.Room,
			Data: jsonMsg,
		}

		// save to database
		msgRecord := models.Message{
			Sender:   c.Username,
			Receiver: c.Receiver,
			Room:     c.Room,
			Content:  string(msg),
		}

		if err := config.DB.Create(&msgRecord).Error; err != nil {
			println("Error in saving ", err)
		}
	}
}

// writePump sends messages from the hub to this client
func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Printf("[WRITE] %s channel closed", c.Username)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("[SEND] Sending message to user %s: %s", c.Username, string(message));
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
