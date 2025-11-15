package handlers

import (
	"encoding/json"
	"net/http"

	"chatapp/auth"
	"chatapp/models"
	"chatapp/store"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userStore *store.UserStore
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(userStore *store.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" || req.Email == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Create user
	user, err := h.userStore.CreateUser(req.Username, req.Email, passwordHash)
	if err != nil {
		if err == store.ErrUserExists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	response := models.AuthResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := h.userStore.GetUser(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	response := models.AuthResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
