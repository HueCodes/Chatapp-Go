package main

import (
	"log"
	"net/http"

	"chatapp/handlers"
	"chatapp/store"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize user store
	userStore := store.NewUserStore()

	// Initialize the WebSocket hub
	hub := handlers.NewHub()
	go hub.Run()

	// Initialize auth handler
	authHandler := handlers.NewAuthHandler(userStore)

	// Create router
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Auth routes
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// Routes
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/ws", handlers.WSHandler(hub)).Methods("GET")

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Authentication enabled - users must register/login to chat")
	log.Fatal(http.ListenAndServe(":8080", router))
}
