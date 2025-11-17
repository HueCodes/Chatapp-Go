package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chatapp/models"
	"chatapp/store"

	"github.com/gorilla/mux"
)

// RoomHandler handles room-related requests
type RoomHandler struct {
	roomStore *store.RoomStore
}

// NewRoomHandler creates a new room handler
func NewRoomHandler(roomStore *store.RoomStore) *RoomHandler {
	return &RoomHandler{
		roomStore: roomStore,
	}
}

// ListRooms handles GET /api/rooms - lists all rooms
func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomStore.GetAllRooms()
	if err != nil {
		http.Error(w, "Failed to retrieve rooms", http.StatusInternalServerError)
		return
	}

	response := make([]models.RoomResponse, len(rooms))
	for i, room := range rooms {
		response[i] = models.RoomResponse{
			ID:        room.ID,
			Name:      room.Name,
			CreatedAt: room.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateRoom handles POST /api/rooms - creates a new room
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	room, err := h.roomStore.CreateRoom(req.Name)
	if err != nil {
		if err == store.ErrRoomExists {
			http.Error(w, "Room already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	response := models.RoomResponse{
		ID:        room.ID,
		Name:      room.Name,
		CreatedAt: room.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetRoom handles GET /api/rooms/{id} - gets a specific room
func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomIDStr := vars["id"]
	
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	room, err := h.roomStore.GetRoom(uint(roomID))
	if err != nil {
		if err == store.ErrRoomNotFound {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve room", http.StatusInternalServerError)
		return
	}

	response := models.RoomResponse{
		ID:        room.ID,
		Name:      room.Name,
		CreatedAt: room.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
