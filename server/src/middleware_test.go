package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestAuthenticationRequired(t *testing.T) {
	server := setupTestServer(t)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/jobs"},
		{"POST", "/api/jobs"},
		{"GET", "/api/jobs/1"},
		{"PUT", "/api/jobs/1"},
		{"DELETE", "/api/jobs/1"},
		{"GET", "/api/jobs/search"},
		{"GET", "/api/auth/me"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			req := httptest.NewRequest(ep.method, ep.path, nil)
			rec := httptest.NewRecorder()

			server.Router.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
			}
		})
	}
}

func TestUserIsolation(t *testing.T) {
	server := setupTestServer(t)

	// Create two users
	user1ID, token1 := createTestUser(t, server, "user1@example.com", "password123", "User 1")
	_, token2 := createTestUser(t, server, "user2@example.com", "password123", "User 2")

	// Create job for user1
	query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, _ := server.DB.Exec(query, user1ID, "Company A", "Position A", "applied",
		now, "Notes", "Contact", "50000", "Location", now, now)

	user1JobID, _ := result.LastInsertId()

	// Try to access user1's job with user2's token
	req := httptest.NewRequest("GET", "/api/jobs/"+strconv.Itoa(int(user1JobID)), nil)
	req.Header.Set("Authorization", "Bearer "+token2)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("User2 should not be able to access User1's job. Expected status %d, got %d",
			http.StatusNotFound, rec.Code)
	}

	// Verify user1 can access their own job
	req = httptest.NewRequest("GET", "/api/jobs/"+strconv.Itoa(int(user1JobID)), nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	rec = httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("User1 should be able to access their own job. Expected status %d, got %d",
			http.StatusOK, rec.Code)
	}
}
