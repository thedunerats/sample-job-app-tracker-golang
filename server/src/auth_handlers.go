package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// handleRegister creates a new user account
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Email, password, and name are required")
		return
	}

	if len(req.Password) < 6 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing password")
		return
	}

	// Insert user into database
	query := `INSERT INTO users (email, password, name, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := s.DB.Exec(query, req.Email, hashedPassword, req.Name, now, now)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			respondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	userID, _ := result.LastInsertId()

	// Generate token
	token, err := GenerateToken(int(userID), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	user := User{
		ID:        int(userID),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// handleLogin authenticates a user and returns a JWT token
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Fetch user from database
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = ?`

	var user User
	err := s.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching user")
		return
	}

	// Check password
	if !CheckPasswordHash(req.Password, user.Password) {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// handleGetCurrentUser returns the current authenticated user's information
func (s *Server) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = ?`

	var user User
	err := s.DB.QueryRow(query, claims.UserID).Scan(
		&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
