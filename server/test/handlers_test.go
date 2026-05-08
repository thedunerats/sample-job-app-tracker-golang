package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func setupTestServer(t *testing.T) *Server {
	dbPath := "./test_handlers.db"
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

func TestHandleGetAllJobs_Empty(t *testing.T) {
	server := setupTestServer(t)
	
	req, _ := http.NewRequest("GET", "/api/jobs", nil)
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var jobs []JobApplication
	json.Unmarshal(rr.Body.Bytes(), &jobs)
	
	if len(jobs) != 0 {
		t.Errorf("Expected empty array, got %d jobs", len(jobs))
	}
}

func TestHandleCreateJob_Success(t *testing.T) {
	server := setupTestServer(t)
	
	jobData := CreateJobApplicationRequest{
		Company:     "Tesla",
		Position:    "Electrical Engineer",
		Status:      "applied",
		AppliedDate: time.Now().Format(time.RFC3339),
		Notes:       "Exciting opportunity",
		ContactInfo: "hr@tesla.com",
		Salary:      "$120k",
		Location:    "Palo Alto",
	}
	
	body, _ := json.Marshal(jobData)
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	
	var job JobApplication
	json.Unmarshal(rr.Body.Bytes(), &job)
	
	if job.Company != "Tesla" {
		t.Errorf("Expected company 'Tesla', got '%s'", job.Company)
	}
	if job.Position != "Electrical Engineer" {
		t.Errorf("Expected position 'Electrical Engineer', got '%s'", job.Position)
	}
	if job.ID == 0 {
		t.Error("Expected non-zero ID")
	}
}

func TestHandleCreateJob_MissingFields(t *testing.T) {
	server := setupTestServer(t)
	
	jobData := CreateJobApplicationRequest{
		Company: "Incomplete Corp",
		// Missing Position and Status
	}
	
	body, _ := json.Marshal(jobData)
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandleGetJob_Success(t *testing.T) {
	server := setupTestServer(t)
	
	// Create a job first
	now := time.Now()
	query := `INSERT INTO job_applications (company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, _ := server.DB.Exec(query, "Netflix", "Data Engineer", "interview", now,
		"Technical round", "careers@netflix.com", "$160k", "Los Gatos", now, now)
	
	id, _ := result.LastInsertId()
	
	// Get the job
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/jobs/%d", id), nil)
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var job JobApplication
	json.Unmarshal(rr.Body.Bytes(), &job)
	
	if job.Company != "Netflix" {
		t.Errorf("Expected company 'Netflix', got '%s'", job.Company)
	}
}

func TestHandleGetJob_NotFound(t *testing.T) {
	server := setupTestServer(t)
	
	req, _ := http.NewRequest("GET", "/api/jobs/9999", nil)
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestHandleUpdateJob_Success(t *testing.T) {
	server := setupTestServer(t)
	
	// Create a job first
	now := time.Now()
	query := `INSERT INTO job_applications (company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, _ := server.DB.Exec(query, "Spotify", "Backend Developer", "applied", now,
		"Music platform", "jobs@spotify.com", "$125k", "Stockholm", now, now)
	
	id, _ := result.LastInsertId()
	
	// Update the job
	newStatus := "offer"
	updateData := UpdateJobApplicationRequest{
		Status: &newStatus,
	}
	
	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/jobs/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var job JobApplication
	json.Unmarshal(rr.Body.Bytes(), &job)
	
	if job.Status != "offer" {
		t.Errorf("Expected status 'offer', got '%s'", job.Status)
	}
}

func TestHandleDeleteJob_Success(t *testing.T) {
	server := setupTestServer(t)
	
	// Create a job first
	now := time.Now()
	query := `INSERT INTO job_applications (company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, _ := server.DB.Exec(query, "Uber", "Mobile Developer", "rejected", now,
		"Not selected", "hiring@uber.com", "$135k", "San Francisco", now, now)
	
	id, _ := result.LastInsertId()
	
	// Delete the job
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/jobs/%d", id), nil)
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Verify it's deleted
	getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/jobs/%d", id), nil)
	getRR := httptest.NewRecorder()
	
	server.Router.ServeHTTP(getRR, getReq)
	
	if status := getRR.Code; status != http.StatusNotFound {
		t.Errorf("Expected NotFound after deletion, got %v", status)
	}
}

func TestHandleGetJobsByStatus_Success(t *testing.T) {
	server := setupTestServer(t)
	
	// Create multiple jobs with different statuses
	now := time.Now()
	query := `INSERT INTO job_applications (company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	server.DB.Exec(query, "Company X", "Role X", "applied", now, "", "", "", "", now, now)
	server.DB.Exec(query, "Company Y", "Role Y", "interview", now, "", "", "", "", now, now)
	server.DB.Exec(query, "Company Z", "Role Z", "applied", now, "", "", "", "", now, now)
	
	// Get jobs with "applied" status
	req, _ := http.NewRequest("GET", "/api/jobs/status/applied", nil)
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var jobs []JobApplication
	json.Unmarshal(rr.Body.Bytes(), &jobs)
	
	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs with status 'applied', got %d", len(jobs))
	}
}

func TestCompleteWorkflow(t *testing.T) {
	server := setupTestServer(t)
	
	// 1. Create a job
	jobData := CreateJobApplicationRequest{
		Company:     "Facebook",
		Position:    "Full Stack Engineer",
		Status:      "applied",
		AppliedDate: time.Now().Format(time.RFC3339),
		Notes:       "Meta opportunity",
		ContactInfo: "recruiting@fb.com",
		Salary:      "$170k",
		Location:    "Menlo Park",
	}
	
	body, _ := json.Marshal(jobData)
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	var createdJob JobApplication
	json.Unmarshal(rr.Body.Bytes(), &createdJob)
	
	// 2. Retrieve all jobs
	req, _ = http.NewRequest("GET", "/api/jobs", nil)
	rr = httptest.NewRecorder()
	server.Router.ServeHTTP(rr, req)
	
	var allJobs []JobApplication
	json.Unmarshal(rr.Body.Bytes(), &allJobs)
	
	if len(allJobs) != 1 {
		t.Errorf("Expected 1 job, got %d", len(allJobs))
	}
	
	// 3. Update the job status
	newStatus := "interview"
	updateData := UpdateJobApplicationRequest{
		Status: &newStatus,
	}
	
	body, _ = json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/jobs/%d", createdJob.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	
	server.Router.ServeHTTP(rr, req)
	
	var updatedJob JobApplication
	json.Unmarshal(rr.Body.Bytes(), &updatedJob)
	
	if updatedJob.Status != "interview" {
		t.Errorf("Expected status 'interview', got '%s'", updatedJob.Status)
	}
	
	// 4. Delete the job
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/jobs/%d", createdJob.ID), nil)
	rr = httptest.NewRecorder()
	server.Router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Delete failed with status %v", status)
	}
}
