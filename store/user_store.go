package store

import (
	"errors"
	"sync"
	"time"

	"chatapp/models"

	"github.com/google/uuid"
)

var (
	ErrUserExists      = errors.New("username already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// UserStore manages user data (in-memory for now)
type UserStore struct {
	mu    sync.RWMutex
	users map[string]*models.User // username -> user
}

// NewUserStore creates a new user store
func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*models.User),
	}
}

// CreateUser creates a new user
func (s *UserStore) CreateUser(username, email, passwordHash string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return nil, ErrUserExists
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	s.users[username] = user
	return user, nil
}

// GetUser retrieves a user by username
func (s *UserStore) GetUser(username string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserStore) GetUserByID(userID string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.ID == userID {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}
