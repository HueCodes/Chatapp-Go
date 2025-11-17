package models

import "time"

// Room represents a chat room
type Room struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// RoomResponse represents a room in API responses
type RoomResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateRoomRequest represents a request to create a room
type CreateRoomRequest struct {
	Name string `json:"name"`
}
