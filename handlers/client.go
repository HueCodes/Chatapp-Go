package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chatapp/models"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

// Client represents a websocket client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	userID   string
	roomID   uint
}

// clientID returns a unique identifier for the client
func (c *Client) clientID() string {
	return c.userID + "-" + c.username
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		var rawMessage map[string]interface{}
		if err := json.Unmarshal(messageBytes, &rawMessage); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Check message type
		msgType, ok := rawMessage["type"].(string)
		if !ok {
			msgType = "text"
		}

		// Handle typing indicators
		if msgType == "typing" {
			isTyping := false
			if isTypingVal, ok := rawMessage["is_typing"].(bool); ok {
				isTyping = isTypingVal
			}

			indicator := &models.TypingIndicator{
				Type:     models.TypingMessage,
				UserID:   c.userID,
				Username: c.username,
				RoomID:   c.roomID,
				IsTyping: isTyping,
			}

			select {
			case c.hub.typingIndicator <- indicator:
			default:
				// Channel full, skip
			}
			continue
		}

		// Handle regular text messages
		var message models.Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Set message properties from the client
		message.UserID = c.userID
		message.Username = c.username
		message.RoomID = c.roomID
		message.Type = models.TextMessage
		message.Timestamp = time.Now()

		broadcastMsg := &BroadcastMessage{
			Message: message,
			RoomID:  c.roomID,
		}

		select {
		case c.hub.broadcast <- broadcastMsg:
		default:
			close(c.send)
			return
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
