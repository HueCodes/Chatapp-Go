package store

import (
	"chatapp/database"
	"chatapp/models"
)

// MessageStore manages message persistence
type MessageStore struct{}

// NewMessageStore creates a new message store
func NewMessageStore() *MessageStore {
	return &MessageStore{}
}

// Save persists a message to the database
func (s *MessageStore) Save(message *models.Message) error {
	// Only persist text messages, not typing indicators or ephemeral messages
	if message.Type != models.TextMessage {
		return nil
	}

	result := database.DB.Create(message)
	return result.Error
}

// GetByRoom retrieves messages for a specific room with a limit
func (s *MessageStore) GetByRoom(roomID uint, limit int) ([]models.Message, error) {
	var messages []models.Message
	result := database.DB.Where("room_id = ? AND type = ?", roomID, models.TextMessage).
		Order("timestamp DESC").
		Limit(limit).
		Find(&messages)
	
	if result.Error != nil {
		return nil, result.Error
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
