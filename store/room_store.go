package store

import (
	"errors"
	"sync"

	"chatapp/database"
	"chatapp/models"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrRoomExists   = errors.New("room already exists")
)

// RoomStore manages chat rooms and their subscriptions
type RoomStore struct {
	mu sync.RWMutex
	// roomID -> map of clients
	rooms map[uint]map[string]bool
}

// NewRoomStore creates a new room store
func NewRoomStore() *RoomStore {
	return &RoomStore{
		rooms: make(map[uint]map[string]bool),
	}
}

// CreateRoom creates a new room
func (s *RoomStore) CreateRoom(name string) (*models.Room, error) {
	room := &models.Room{
		Name: name,
	}

	result := database.DB.Create(room)
	if result.Error != nil {
		return nil, ErrRoomExists
	}

	s.mu.Lock()
	s.rooms[room.ID] = make(map[string]bool)
	s.mu.Unlock()

	return room, nil
}

// GetRoom retrieves a room by ID
func (s *RoomStore) GetRoom(roomID uint) (*models.Room, error) {
	var room models.Room
	result := database.DB.First(&room, roomID)
	if result.Error != nil {
		return nil, ErrRoomNotFound
	}

	s.mu.Lock()
	if _, exists := s.rooms[roomID]; !exists {
		s.rooms[roomID] = make(map[string]bool)
	}
	s.mu.Unlock()

	return &room, nil
}

// GetAllRooms retrieves all rooms
func (s *RoomStore) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	result := database.DB.Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}

	s.mu.Lock()
	for _, room := range rooms {
		if _, exists := s.rooms[room.ID]; !exists {
			s.rooms[room.ID] = make(map[string]bool)
		}
	}
	s.mu.Unlock()

	return rooms, nil
}

// AddClientToRoom adds a client to a room
func (s *RoomStore) AddClientToRoom(roomID uint, clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.rooms[roomID]; !exists {
		s.rooms[roomID] = make(map[string]bool)
	}
	s.rooms[roomID][clientID] = true
}

// RemoveClientFromRoom removes a client from a room
func (s *RoomStore) RemoveClientFromRoom(roomID uint, clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if clients, exists := s.rooms[roomID]; exists {
		delete(clients, clientID)
	}
}

// GetRoomClients returns all client IDs in a room
func (s *RoomStore) GetRoomClients(roomID uint) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make([]string, 0)
	if roomClients, exists := s.rooms[roomID]; exists {
		for clientID := range roomClients {
			clients = append(clients, clientID)
		}
	}
	return clients
}
