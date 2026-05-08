package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Server holds the database and router
type Server struct {
	DB     *Database
	Router *mux.Router
}

// NewServer creates a new server instance
func NewServer(db *Database) *Server {
	s := &Server{
		DB:     db,
		Router: mux.NewRouter(),
	}
	s.routes()
	return s
}

// routes defines all the API routes
func (s *Server) routes() {
	s.Router.HandleFunc("/api/jobs", s.handleGetAllJobs).Methods("GET")
	s.Router.HandleFunc("/api/jobs", s.handleCreateJob).Methods("POST")
	s.Router.HandleFunc("/api/jobs/{id}", s.handleGetJob).Methods("GET")
	s.Router.HandleFunc("/api/jobs/{id}", s.handleUpdateJob).Methods("PUT")
	s.Router.HandleFunc("/api/jobs/{id}", s.handleDeleteJob).Methods("DELETE")
	s.Router.HandleFunc("/api/jobs/status/{status}", s.handleGetJobsByStatus).Methods("GET")
}

// handleGetAllJobs retrieves all job applications
func (s *Server) handleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, company, position, status, applied_date, notes, 
			  contact_info, salary, location, created_at, updated_at 
			  FROM job_applications ORDER BY applied_date DESC`

	rows, err := s.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching jobs")
		return
	}
	defer rows.Close()

	jobs := []JobApplication{}
	for rows.Next() {
		var job JobApplication
		err := rows.Scan(&job.ID, &job.Company, &job.Position, &job.Status,
			&job.AppliedDate, &job.Notes, &job.ContactInfo, &job.Salary,
			&job.Location, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning job")
			return
		}
		jobs = append(jobs, job)
	}

	respondWithJSON(w, http.StatusOK, jobs)
}

// handleCreateJob creates a new job application
func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req CreateJobApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Company == "" || req.Position == "" || req.Status == "" {
		respondWithError(w, http.StatusBadRequest, "Company, position, and status are required")
		return
	}

	appliedDate, err := time.Parse(time.RFC3339, req.AppliedDate)
	if err != nil {
		appliedDate = time.Now()
	}

	query := `INSERT INTO job_applications (company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := s.DB.Exec(query, req.Company, req.Position, req.Status, appliedDate,
		req.Notes, req.ContactInfo, req.Salary, req.Location, now, now)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating job application")
		return
	}

	id, _ := result.LastInsertId()
	job := JobApplication{
		ID:          int(id),
		Company:     req.Company,
		Position:    req.Position,
		Status:      req.Status,
		AppliedDate: appliedDate,
		Notes:       req.Notes,
		ContactInfo: req.ContactInfo,
		Salary:      req.Salary,
		Location:    req.Location,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	respondWithJSON(w, http.StatusCreated, job)
}

// handleGetJob retrieves a single job application by ID
func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	query := `SELECT id, company, position, status, applied_date, notes, 
			  contact_info, salary, location, created_at, updated_at 
			  FROM job_applications WHERE id = ?`

	var job JobApplication
	err = s.DB.QueryRow(query, id).Scan(&job.ID, &job.Company, &job.Position,
		&job.Status, &job.AppliedDate, &job.Notes, &job.ContactInfo,
		&job.Salary, &job.Location, &job.CreatedAt, &job.UpdatedAt)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Job application not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching job application")
		return
	}

	respondWithJSON(w, http.StatusOK, job)
}

// handleUpdateJob updates an existing job application
func (s *Server) handleUpdateJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	var req UpdateJobApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Build dynamic update query
	query := "UPDATE job_applications SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Company != nil {
		query += ", company = ?"
		args = append(args, *req.Company)
	}
	if req.Position != nil {
		query += ", position = ?"
		args = append(args, *req.Position)
	}
	if req.Status != nil {
		query += ", status = ?"
		args = append(args, *req.Status)
	}
	if req.AppliedDate != nil {
		appliedDate, err := time.Parse(time.RFC3339, *req.AppliedDate)
		if err == nil {
			query += ", applied_date = ?"
			args = append(args, appliedDate)
		}
	}
	if req.Notes != nil {
		query += ", notes = ?"
		args = append(args, *req.Notes)
	}
	if req.ContactInfo != nil {
		query += ", contact_info = ?"
		args = append(args, *req.ContactInfo)
	}
	if req.Salary != nil {
		query += ", salary = ?"
		args = append(args, *req.Salary)
	}
	if req.Location != nil {
		query += ", location = ?"
		args = append(args, *req.Location)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := s.DB.Exec(query, args...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating job application")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Job application not found")
		return
	}

	// Fetch and return the updated job
	s.handleGetJob(w, r)
}

// handleDeleteJob deletes a job application
func (s *Server) handleDeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	query := "DELETE FROM job_applications WHERE id = ?"
	result, err := s.DB.Exec(query, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting job application")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Job application not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Job application deleted successfully"})
}

// handleGetJobsByStatus retrieves all job applications with a specific status
func (s *Server) handleGetJobsByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	query := `SELECT id, company, position, status, applied_date, notes, 
			  contact_info, salary, location, created_at, updated_at 
			  FROM job_applications WHERE status = ? ORDER BY applied_date DESC`

	rows, err := s.DB.Query(query, status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching jobs")
		return
	}
	defer rows.Close()

	jobs := []JobApplication{}
	for rows.Next() {
		var job JobApplication
		err := rows.Scan(&job.ID, &job.Company, &job.Position, &job.Status,
			&job.AppliedDate, &job.Notes, &job.ContactInfo, &job.Salary,
			&job.Location, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning job")
			return
		}
		jobs = append(jobs, job)
	}

	respondWithJSON(w, http.StatusOK, jobs)
}

// respondWithError sends an error response in JSON format
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
