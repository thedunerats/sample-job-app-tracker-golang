package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

func setupAuthTestServer(t *testing.T) *Server {
	dbPath := "./test_auth.db"
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
	server := setupAuthTestServer(t)

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
	server := setupAuthTestServer(t)

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
	server := setupAuthTestServer(t)

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
	server := setupAuthTestServer(t)

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

func TestPagination(t *testing.T) {
	server := setupAuthTestServer(t)

	// Create user and jobs
	userID, token := createTestUser(t, server, "page@example.com", "password123", "Page User")

	// Create 25 test jobs
	now := time.Now()
	for i := 1; i <= 25; i++ {
		query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
				  notes, contact_info, salary, location, created_at, updated_at) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		company := fmt.Sprintf("Company%d", i)
		position := fmt.Sprintf("Position%d", i)
		server.DB.Exec(query, userID, company, position, "applied",
			now, "Notes", "Contact", "50000", "Location", now, now)
	}

	tests := []struct {
		name          string
		page          int
		limit         int
		expectedCount int
		expectedPages int
		expectedTotal int
	}{
		{"First page", 1, 10, 10, 3, 25},
		{"Second page", 2, 10, 10, 3, 25},
		{"Last page", 3, 10, 5, 3, 25},
		{"Larger page size", 1, 20, 20, 2, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/jobs?page=%d&limit=%d", tt.page, tt.limit)
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			server.Router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
				return
			}

			var response PaginatedResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			jobs := response.Data.([]interface{})
			if len(jobs) != tt.expectedCount {
				t.Errorf("Expected %d jobs, got %d", tt.expectedCount, len(jobs))
			}

			if response.TotalPages != tt.expectedPages {
				t.Errorf("Expected %d total pages, got %d", tt.expectedPages, response.TotalPages)
			}

			if response.TotalCount != tt.expectedTotal {
				t.Errorf("Expected total count %d, got %d", tt.expectedTotal, response.TotalCount)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	server := setupAuthTestServer(t)

	// Create user
	userID, token := createTestUser(t, server, "search@example.com", "password123", "Search User")

	// Create test jobs with specific attributes
	testJobs := []struct {
		company  string
		position string
		status   string
		location string
	}{
		{"Google", "Software Engineer", "applied", "Mountain View"},
		{"Microsoft", "Software Engineer", "interview", "Seattle"},
		{"Amazon", "Data Scientist", "offer", "Seattle"},
		{"Google", "Product Manager", "rejected", "New York"},
		{"Apple", "Software Engineer", "applied", "Cupertino"},
	}

	now := time.Now()
	for _, job := range testJobs {
		query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
				  notes, contact_info, salary, location, created_at, updated_at) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		server.DB.Exec(query, userID, job.company, job.position, job.status,
			now, "Notes", "Contact", "100000", job.location, now, now)
	}

	tests := []struct {
		name          string
		params        map[string]string
		expectedCount int
	}{
		{"Search by company", map[string]string{"company": "Google"}, 2},
		{"Search by position", map[string]string{"position": "Software Engineer"}, 3},
		{"Search by status", map[string]string{"status": "applied"}, 2},
		{"Search by location", map[string]string{"location": "Seattle"}, 2},
		{"Combined search", map[string]string{"company": "Google", "position": "Software Engineer"}, 1},
		{"No matches", map[string]string{"company": "NonExistent"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/api/jobs/search?"
			for key, value := range tt.params {
				urlPath += key + "=" + url.QueryEscape(value) + "&"
			}

			req := httptest.NewRequest("GET", urlPath, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			server.Router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
				return
			}

			var response PaginatedResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			jobs := response.Data.([]interface{})
			if len(jobs) != tt.expectedCount {
				t.Errorf("Expected %d jobs, got %d", tt.expectedCount, len(jobs))
			}
		})
	}
}

func TestAuthenticationRequired(t *testing.T) {
	server := setupAuthTestServer(t)

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
	server := setupAuthTestServer(t)

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

func TestCreateJobWithAuth(t *testing.T) {
	server := setupAuthTestServer(t)

	_, token := createTestUser(t, server, "create@example.com", "password123", "Create User")

	jobReq := CreateJobApplicationRequest{
		Company:     "Test Company",
		Position:    "Test Position",
		Status:      "applied",
		AppliedDate: time.Now().Format(time.RFC3339),
		Notes:       "Test notes",
		ContactInfo: "test@test.com",
		Salary:      "100000",
		Location:    "Test Location",
	}

	body, _ := json.Marshal(jobReq)
	req := httptest.NewRequest("POST", "/api/jobs", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var job JobApplication
	if err := json.Unmarshal(rec.Body.Bytes(), &job); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if job.Company != jobReq.Company {
		t.Errorf("Expected company %s, got %s", jobReq.Company, job.Company)
	}
}
