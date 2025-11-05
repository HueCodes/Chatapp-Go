package models

import "time"

// MessageType represents the type of message
type MessageType string

const (
	TextMessage     MessageType = "text"
	UserJoinMessage MessageType = "user_join"
	UserLeftMessage MessageType = "user_left"
	SystemMessage   MessageType = "system"
)

// Message represents a chat message
type Message struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	Username  string      `json:"username"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
	RoomID    string      `json:"room_id,omitempty"`
}

// User represents a chat user
// NOTE: User and Room types were removed to keep the initial codebase minimal.
// Add domain types as needed when implementing rooms, persistence, or user
// management features.
