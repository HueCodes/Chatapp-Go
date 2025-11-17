package models

import (
	"time"

	"gorm.io/gorm"
)

// MessageType represents the type of message
type MessageType string

const (
	TextMessage      MessageType = "text"
	UserJoinMessage  MessageType = "user_join"
	UserLeftMessage  MessageType = "user_left"
	SystemMessage    MessageType = "system"
	TypingMessage    MessageType = "typing"
)

// Message represents a chat message (both in-memory and persisted)
type Message struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Type      MessageType     `gorm:"size:20;not null" json:"type"`
	UserID    string          `gorm:"size:100;index" json:"user_id"`
	Username  string          `gorm:"size:100;not null" json:"username"`
	RoomID    uint            `gorm:"index;not null" json:"room_id"`
	Content   string          `gorm:"type:text;not null" json:"content"`
	Timestamp time.Time       `gorm:"autoCreateTime" json:"timestamp"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TypingIndicator represents a typing indicator message
type TypingIndicator struct {
	Type      MessageType `json:"type"`
	UserID    string      `json:"user_id"`
	Username  string      `json:"username"`
	RoomID    uint        `json:"room_id"`
	IsTyping  bool        `json:"is_typing"`
}
