package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGetAllJobs_Empty(t *testing.T) {
	server := setupTestServer(t)

	_, token := createTestUser(t, server, "empty@example.com", "password123", "Empty User")

	req := httptest.NewRequest("GET", "/api/jobs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response PaginatedResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	jobs := response.Data.([]interface{})
	if len(jobs) != 0 {
		t.Errorf("Expected empty job list, got %d jobs", len(jobs))
	}
}

func TestCreateJob_Success(t *testing.T) {
	server := setupTestServer(t)

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

func TestCreateJob_MissingFields(t *testing.T) {
	server := setupTestServer(t)

	_, token := createTestUser(t, server, "validation@example.com", "password123", "Validation User")

	tests := []struct {
		name    string
		request CreateJobApplicationRequest
	}{
		{
			name: "Missing company",
			request: CreateJobApplicationRequest{
				Position: "Engineer",
				Status:   "applied",
			},
		},
		{
			name: "Missing position",
			request: CreateJobApplicationRequest{
				Company: "Test Corp",
				Status:  "applied",
			},
		},
		{
			name: "Missing status",
			request: CreateJobApplicationRequest{
				Company:  "Test Corp",
				Position: "Engineer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/jobs", bytes.NewReader(body))
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			server.Router.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestGetJob_Success(t *testing.T) {
	server := setupTestServer(t)

	userID, token := createTestUser(t, server, "getjob@example.com", "password123", "Get Job User")

	// Create a job first
	now := time.Now()
	query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, _ := server.DB.Exec(query, userID, "Netflix", "Data Engineer", "interview", now,
		"Technical round", "careers@netflix.com", "$160k", "Los Gatos", now, now)

	id, _ := result.LastInsertId()

	// Get the job
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/jobs/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var job JobApplication
	if err := json.Unmarshal(rec.Body.Bytes(), &job); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if job.Company != "Netflix" {
		t.Errorf("Expected company 'Netflix', got '%s'", job.Company)
	}

	if job.Position != "Data Engineer" {
		t.Errorf("Expected position 'Data Engineer', got '%s'", job.Position)
	}
}

func TestGetJob_NotFound(t *testing.T) {
	server := setupTestServer(t)

	_, token := createTestUser(t, server, "notfound@example.com", "password123", "Not Found User")

	req := httptest.NewRequest("GET", "/api/jobs/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateJob_Success(t *testing.T) {
	server := setupTestServer(t)

	userID, token := createTestUser(t, server, "update@example.com", "password123", "Update User")

	// Create a job first
	now := time.Now()
	query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, _ := server.DB.Exec(query, userID, "Spotify", "Backend Developer", "applied", now,
		"Music platform", "jobs@spotify.com", "$125k", "Stockholm", now, now)

	id, _ := result.LastInsertId()

	// Update the job
	newStatus := "offer"
	newSalary := "$135k"
	updateData := UpdateJobApplicationRequest{
		Status: &newStatus,
		Salary: &newSalary,
	}

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/jobs/%d", id), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var job JobApplication
	if err := json.Unmarshal(rec.Body.Bytes(), &job); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if job.Status != "offer" {
		t.Errorf("Expected status 'offer', got '%s'", job.Status)
	}

	if job.Salary != "$135k" {
		t.Errorf("Expected salary '$135k', got '%s'", job.Salary)
	}
}

func TestDeleteJob_Success(t *testing.T) {
	server := setupTestServer(t)

	userID, token := createTestUser(t, server, "delete@example.com", "password123", "Delete User")

	// Create a job first
	now := time.Now()
	query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, _ := server.DB.Exec(query, userID, "Uber", "Mobile Developer", "rejected", now,
		"Not selected", "hiring@uber.com", "$135k", "San Francisco", now, now)

	id, _ := result.LastInsertId()

	// Delete the job
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/jobs/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify it's deleted
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/api/jobs/%d", id), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRec := httptest.NewRecorder()

	server.Router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Errorf("Expected NotFound after deletion, got %d", getRec.Code)
	}
}

func TestGetJobsByStatus(t *testing.T) {
	server := setupTestServer(t)

	userID, token := createTestUser(t, server, "status@example.com", "password123", "Status User")

	// Create multiple jobs with different statuses
	now := time.Now()
	query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	server.DB.Exec(query, userID, "Company X", "Role X", "applied", now, "", "", "", "", now, now)
	server.DB.Exec(query, userID, "Company Y", "Role Y", "interview", now, "", "", "", "", now, now)
	server.DB.Exec(query, userID, "Company Z", "Role Z", "applied", now, "", "", "", "", now, now)
	server.DB.Exec(query, userID, "Company W", "Role W", "offer", now, "", "", "", "", now, now)

	// Get jobs with "applied" status
	req := httptest.NewRequest("GET", "/api/jobs/status/applied", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var jobs []JobApplication
	if err := json.Unmarshal(rec.Body.Bytes(), &jobs); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs with status 'applied', got %d", len(jobs))
	}

	// Verify all returned jobs have the correct status
	for _, job := range jobs {
		if job.Status != "applied" {
			t.Errorf("Expected all jobs to have status 'applied', got '%s'", job.Status)
		}
	}
}

func TestPagination(t *testing.T) {
	server := setupTestServer(t)

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
	server := setupTestServer(t)

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

func TestCompleteWorkflow(t *testing.T) {
	server := setupTestServer(t)

	_, token := createTestUser(t, server, "workflow@example.com", "password123", "Workflow User")

	// 1. Create a job
	jobReq := CreateJobApplicationRequest{
		Company:     "Meta",
		Position:    "Full Stack Engineer",
		Status:      "applied",
		AppliedDate: time.Now().Format(time.RFC3339),
		Notes:       "Meta opportunity",
		ContactInfo: "recruiting@meta.com",
		Salary:      "$170k",
		Location:    "Menlo Park",
	}

	body, _ := json.Marshal(jobReq)
	req := httptest.NewRequest("POST", "/api/jobs", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Failed to create job: %d", rec.Code)
	}

	var createdJob JobApplication
	json.Unmarshal(rec.Body.Bytes(), &createdJob)

	// 2. Retrieve all jobs
	req = httptest.NewRequest("GET", "/api/jobs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Failed to get jobs: %d", rec.Code)
	}

	var paginatedResponse PaginatedResponse
	json.Unmarshal(rec.Body.Bytes(), &paginatedResponse)

	jobs := paginatedResponse.Data.([]interface{})
	if len(jobs) != 1 {
		t.Errorf("Expected 1 job, got %d", len(jobs))
	}

	// 3. Update the job status
	newStatus := "interview"
	updateData := UpdateJobApplicationRequest{
		Status: &newStatus,
	}

	body, _ = json.Marshal(updateData)
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/jobs/%d", createdJob.ID), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Failed to update job: %d", rec.Code)
	}

	var updatedJob JobApplication
	json.Unmarshal(rec.Body.Bytes(), &updatedJob)

	if updatedJob.Status != "interview" {
		t.Errorf("Expected status 'interview', got '%s'", updatedJob.Status)
	}

	// 4. Search for the job
	req = httptest.NewRequest("GET", "/api/jobs/search?company=Meta", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Failed to search jobs: %d", rec.Code)
	}

	var searchResponse PaginatedResponse
	json.Unmarshal(rec.Body.Bytes(), &searchResponse)

	searchJobs := searchResponse.Data.([]interface{})
	if len(searchJobs) != 1 {
		t.Errorf("Expected 1 job in search results, got %d", len(searchJobs))
	}

	// 5. Delete the job
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/jobs/%d", createdJob.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Failed to delete job: %d", rec.Code)
	}

	// 6. Verify deletion
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/jobs/%d", createdJob.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()

	server.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected NotFound after deletion, got %d", rec.Code)
	}
}
