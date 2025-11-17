package main

import (
	"log"
	"net/http"

	"chatapp/database"
	"chatapp/handlers"
	"chatapp/models"
	"chatapp/store"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	if err := database.InitDB("chatapp.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate models
	if err := database.AutoMigrate(&models.Message{}, &models.Room{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize stores
	userStore := store.NewUserStore()
	roomStore := store.NewRoomStore()
	messageStore := store.NewMessageStore()

	// Create default room if it doesn't exist
	defaultRoom, err := roomStore.GetRoom(1)
	if err != nil {
		defaultRoom, err = roomStore.CreateRoom("General")
		if err != nil {
			log.Printf("Warning: Could not create default room: %v", err)
		} else {
			log.Printf("Created default room: %s (ID: %d)", defaultRoom.Name, defaultRoom.ID)
		}
	}

	// Initialize the WebSocket hub
	hub := handlers.NewHub(roomStore, messageStore)
	go hub.Run()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userStore)
	roomHandler := handlers.NewRoomHandler(roomStore)

	// Create router
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Auth routes
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// Room routes
	router.HandleFunc("/api/rooms", roomHandler.ListRooms).Methods("GET")
	router.HandleFunc("/api/rooms", roomHandler.CreateRoom).Methods("POST")
	router.HandleFunc("/api/rooms/{id}", roomHandler.GetRoom).Methods("GET")

	// WebSocket route
	router.HandleFunc("/ws", handlers.WSHandler(hub, roomStore)).Methods("GET")

	// Home route
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Authentication enabled - users must register/login to chat")
	log.Println("Database: chatapp.db")
	log.Fatal(http.ListenAndServe(":8080", router))
}
