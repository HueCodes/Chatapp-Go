package handlers

import (
	"encoding/json"
	"log"
	"time"

	"chatapp/models"
	"chatapp/store"
)

const (
	// how many recent messages to send to a newly connected client
	recentMessagesToSend = 50
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients (keyed by client pointer)
	clients map[*Client]bool

	// Rooms store for managing room subscriptions
	roomStore *store.RoomStore

	// Message store for persistence
	messageStore *store.MessageStore

	// Inbound messages from the clients
	broadcast chan *BroadcastMessage

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Typing indicator channel
	typingIndicator chan *models.TypingIndicator
}

// BroadcastMessage wraps a message with room information
type BroadcastMessage struct {
	Message models.Message
	RoomID  uint
}

// NewHub creates a new Hub instance
func NewHub(roomStore *store.RoomStore, messageStore *store.MessageStore) *Hub {
	return &Hub{
		broadcast:       make(chan *BroadcastMessage),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		typingIndicator: make(chan *models.TypingIndicator),
		clients:         make(map[*Client]bool),
		roomStore:       roomStore,
		messageStore:    messageStore,
	}
}

// Run starts the hub and handles client registration/unregistration and message broadcasting
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.roomStore.AddClientToRoom(client.roomID, client.clientID())
			log.Printf("Client connected: %s to room %d", client.username, client.roomID)

			// Send recent messages to newly connected client
			h.sendRecentMessages(client)

			// Broadcast user join message
			joinMessage := models.Message{
				Type:      models.UserJoinMessage,
				UserID:    client.userID,
				Username:  client.username,
				RoomID:    client.roomID,
				Content:   client.username + " joined the chat",
				Timestamp: time.Now(),
			}
			h.broadcastToRoom(&BroadcastMessage{Message: joinMessage, RoomID: client.roomID})

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.roomStore.RemoveClientFromRoom(client.roomID, client.clientID())
				close(client.send)
				log.Printf("Client disconnected: %s from room %d", client.username, client.roomID)

				// Broadcast user left message
				leftMessage := models.Message{
					Type:      models.UserLeftMessage,
					UserID:    client.userID,
					Username:  client.username,
					RoomID:    client.roomID,
					Content:   client.username + " left the chat",
					Timestamp: time.Now(),
				}
				h.broadcastToRoom(&BroadcastMessage{Message: leftMessage, RoomID: client.roomID})
			}

		case broadcastMsg := <-h.broadcast:
			// Save text messages to database asynchronously
			if broadcastMsg.Message.Type == models.TextMessage {
				go func(msg models.Message) {
					if err := h.messageStore.Save(&msg); err != nil {
						log.Printf("Error saving message to database: %v", err)
					}
				}(broadcastMsg.Message)
			}

			h.broadcastToRoom(broadcastMsg)

		case typingIndicator := <-h.typingIndicator:
			h.broadcastTypingIndicator(typingIndicator)
		}
	}
}

// broadcastToRoom sends a message to all clients in a specific room
func (h *Hub) broadcastToRoom(broadcastMsg *BroadcastMessage) {
	messageBytes, err := json.Marshal(broadcastMsg.Message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range h.clients {
		if client.roomID == broadcastMsg.RoomID {
			select {
			case client.send <- messageBytes:
			default:
				close(client.send)
				delete(h.clients, client)
				h.roomStore.RemoveClientFromRoom(client.roomID, client.clientID())
			}
		}
	}
}

// broadcastTypingIndicator sends a typing indicator to all clients in a room
func (h *Hub) broadcastTypingIndicator(indicator *models.TypingIndicator) {
	indicatorBytes, err := json.Marshal(indicator)
	if err != nil {
		log.Printf("Error marshaling typing indicator: %v", err)
		return
	}

	for client := range h.clients {
		// Don't send typing indicator to the user who is typing
		if client.roomID == indicator.RoomID && client.userID != indicator.UserID {
			select {
			case client.send <- indicatorBytes:
			default:
				// Client buffer full, skip
			}
		}
	}
}

// sendRecentMessages sends the last 50 messages to a newly connected client
func (h *Hub) sendRecentMessages(client *Client) {
	messages, err := h.messageStore.GetByRoom(client.roomID, recentMessagesToSend)
	if err != nil {
		log.Printf("Error retrieving recent messages: %v", err)
		return
	}

	for _, message := range messages {
		messageBytes, err := json.Marshal(message)
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
