package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"chatapp/auth"
	"chatapp/store"
)

// HomeHandler serves the main chat page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the static HTML file. This keeps handlers lightweight and lets the
	// static directory contain the full UI (HTML/CSS/JS).
	http.ServeFile(w, r, "static/index.html")
}

// WSHandler handles websocket requests from the peer
func WSHandler(hub *Hub, roomStore *store.RoomStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from query parameter or Authorization header
		token := r.URL.Query().Get("token")
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			http.Error(w, "Authentication token is required", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Get room_id from query parameter (default to 1)
		roomIDStr := r.URL.Query().Get("room_id")
		roomID := uint(1)
		if roomIDStr != "" {
			parsed, err := strconv.ParseUint(roomIDStr, 10, 32)
			if err != nil {
				http.Error(w, "Invalid room_id", http.StatusBadRequest)
				return
			}
			roomID = uint(parsed)
		}

		// Validate room exists
		_, err = roomStore.GetRoom(roomID)
		if err != nil {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := &Client{
			hub:      hub,
			conn:     conn,
			send:     make(chan []byte, 256),
			username: claims.Username,
			userID:   claims.UserID,
			roomID:   roomID,
		}

		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in new goroutines
		go client.writePump()
		go client.readPump()
	}
}
