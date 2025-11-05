package handlers

import (
	"encoding/json"
	"log"
	"time"

	"chatapp/models"

	"github.com/google/uuid"
)

const (
	// maximum number of messages kept in memory
	maxMessageHistory = 100
	// how many recent messages to send to a newly connected client
	recentMessagesToSend = 50
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Store recent messages (in-memory for simplicity)
	messages []models.Message
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		messages:   make([]models.Message, 0),
	}
}

// Run starts the hub and handles client registration/unregistration and message broadcasting
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected: %s", client.username)

			// Send recent messages to newly connected client
			h.sendRecentMessages(client)

			// Broadcast user join message
			joinMessage := models.Message{
				ID:        uuid.New().String(),
				Type:      models.UserJoinMessage,
				Username:  client.username,
				Content:   client.username + " joined the chat",
				Timestamp: time.Now(),
			}
			h.addMessage(joinMessage)
			h.broadcastMessage(joinMessage)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected: %s", client.username)

				// Broadcast user left message
				leftMessage := models.Message{
					ID:        uuid.New().String(),
					Type:      models.UserLeftMessage,
					Username:  client.username,
					Content:   client.username + " left the chat",
					Timestamp: time.Now(),
				}
				h.addMessage(leftMessage)
				h.broadcastMessage(leftMessage)
			}

		case message := <-h.broadcast:
			var msg models.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Add ID and timestamp to message
			msg.ID = uuid.New().String()
			msg.Timestamp = time.Now()

			h.addMessage(msg)
			h.broadcastMessage(msg)
		}
	}
}

// addMessage adds a message to the hub's message history
func (h *Hub) addMessage(message models.Message) {
	h.messages = append(h.messages, message)

	// Keep only the last maxMessageHistory messages
	if len(h.messages) > maxMessageHistory {
		// keep the last maxMessageHistory items
		h.messages = h.messages[len(h.messages)-maxMessageHistory:]
	}
}

// broadcastMessage sends a message to all connected clients
func (h *Hub) broadcastMessage(message models.Message) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range h.clients {
		select {
		case client.send <- messageBytes:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// sendRecentMessages sends the last 50 messages to a newly connected client
func (h *Hub) sendRecentMessages(client *Client) {
	start := 0
	if len(h.messages) > recentMessagesToSend {
		start = len(h.messages) - recentMessagesToSend
	}

	for i := start; i < len(h.messages); i++ {
		messageBytes, err := json.Marshal(h.messages[i])
		if err != nil {
			log.Printf("Error marshaling recent message: %v", err)
			continue
		}

		select {
		case client.send <- messageBytes:
		default:
			return
		}
	}
}
