package handlers

import (
	"log"
	"net/http"
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
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
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
			username: username,
		}

		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in new goroutines
		go client.writePump()
		go client.readPump()
	}
}
