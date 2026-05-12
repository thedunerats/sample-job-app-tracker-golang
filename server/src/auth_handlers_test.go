package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestServer(t *testing.T) *Server {
	dbPath := "./test_server.db"
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return NewServer(db)
}

func createTestUser(t *testing.T, server *Server, email, password, name string) (int, string) {
	regReq := RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	}

	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Failed to create test user. Status: %d, Body: %s", rec.Code, rec.Body.String())
	}

	var response LoginResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	return response.User.ID, response.Token
}

func TestRegister_Success(t *testing.T) {
	server := setupTestServer(t)

	regReq := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Token == "" {
		t.Error("Token should not be empty")
	}

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", response.User.Email)
	}
}

func TestRegister_ValidationErrors(t *testing.T) {
	server := setupTestServer(t)

	tests := []struct {
		name           string
		request        RegisterRequest
		expectedStatus int
	}{
		{
			name:           "Missing email",
			request:        RegisterRequest{Password: "password123", Name: "Test User"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing password",
			request:        RegisterRequest{Email: "test@example.com", Name: "Test User"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Short password",
			request:        RegisterRequest{Email: "test@example.com", Password: "123", Name: "Test User"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing name",
			request:        RegisterRequest{Email: "test@example.com", Password: "password123"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			server.Router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestLogin_Success(t *testing.T) {
	server := setupTestServer(t)

	// Create user first
	email := "login@example.com"
	password := "password123"
	createTestUser(t, server, email, password, "Login User")

	// Now login
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Token == "" {
		t.Error("Token should not be empty")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	server := setupTestServer(t)

	// Create user
	email := "test@example.com"
	password := "password123"
	createTestUser(t, server, email, password, "Test User")

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"Wrong password", email, "wrongpassword"},
		{"Non-existent user", "nonexistent@example.com", "password123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginReq := LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			server.Router.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
			}
		})
	}
}
