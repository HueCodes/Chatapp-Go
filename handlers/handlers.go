package handlers

import (
	"log"
	"net/http"
	"strings"

	"chatapp/auth"
)

// HomeHandler serves the main chat page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the static HTML file. This keeps handlers lightweight and lets the
	// static directory contain the full UI (HTML/CSS/JS).
	http.ServeFile(w, r, "static/index.html")
}

// WSHandler handles websocket requests from the peer
func WSHandler(hub *Hub) http.HandlerFunc {
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
		}

		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in new goroutines
		go client.writePump()
		go client.readPump()
	}
}
