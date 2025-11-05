package main

import (
	"log"
	"net/http"

	"chatapp/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the WebSocket hub
	hub := handlers.NewHub()
	go hub.Run()

	// Create router
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Routes
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/ws", handlers.WSHandler(hub)).Methods("GET")

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
