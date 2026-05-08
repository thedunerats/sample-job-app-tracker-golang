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
	// Public routes
	s.Router.HandleFunc("/api/auth/register", s.handleRegister).Methods("POST")
	s.Router.HandleFunc("/api/auth/login", s.handleLogin).Methods("POST")

	// Protected routes (require authentication)
	api := s.Router.PathPrefix("/api").Subrouter()
	api.Use(s.AuthMiddleware)

	api.HandleFunc("/auth/me", s.handleGetCurrentUser).Methods("GET")
	api.HandleFunc("/jobs", s.handleGetAllJobs).Methods("GET")
	api.HandleFunc("/jobs", s.handleCreateJob).Methods("POST")
	api.HandleFunc("/jobs/search", s.handleSearchJobs).Methods("GET")
	api.HandleFunc("/jobs/{id}", s.handleGetJob).Methods("GET")
	api.HandleFunc("/jobs/{id}", s.handleUpdateJob).Methods("PUT")
	api.HandleFunc("/jobs/{id}", s.handleDeleteJob).Methods("DELETE")
	api.HandleFunc("/jobs/status/{status}", s.handleGetJobsByStatus).Methods("GET")
}

// handleGetAllJobs retrieves all job applications with pagination
func (s *Server) handleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse pagination parameters
	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 10)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM job_applications WHERE user_id = ?`
	err := s.DB.QueryRow(countQuery, claims.UserID).Scan(&totalCount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error counting jobs")
		return
	}

	// Get paginated results
	query := `SELECT id, user_id, company, position, status, applied_date, notes, 
			  contact_info, salary, location, created_at, updated_at 
			  FROM job_applications WHERE user_id = ? 
			  ORDER BY applied_date DESC LIMIT ? OFFSET ?`

	rows, err := s.DB.Query(query, claims.UserID, limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching jobs")
		return
	}
	defer rows.Close()

	jobs := []JobApplication{}
	for rows.Next() {
		var job JobApplication
		err := rows.Scan(&job.ID, &job.UserID, &job.Company, &job.Position, &job.Status,
			&job.AppliedDate, &job.Notes, &job.ContactInfo, &job.Salary,
			&job.Location, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning job")
			return
		}
		jobs = append(jobs, job)
	}

	totalPages := (totalCount + limit - 1) / limit
	response := PaginatedResponse{
		Data:       jobs,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// handleCreateJob creates a new job application
func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

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

	query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := s.DB.Exec(query, claims.UserID, req.Company, req.Position, req.Status, appliedDate,
		req.Notes, req.ContactInfo, req.Salary, req.Location, now, now)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating job application")
		return
	}

	id, _ := result.LastInsertId()
	job := JobApplication{
		ID:          int(id),
		UserID:      claims.UserID,
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
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	query := `SELECT id, user_id, company, position, status, applied_date, notes, 
			  contact_info, salary, location, created_at, updated_at 
			  FROM job_applications WHERE id = ? AND user_id = ?`

	var job JobApplication
	err = s.DB.QueryRow(query, id, claims.UserID).Scan(&job.ID, &job.UserID, &job.Company, &job.Position,
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
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

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

	query += " WHERE id = ? AND user_id = ?"
	args = append(args, id, claims.UserID)

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
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	query := "DELETE FROM job_applications WHERE id = ? AND user_id = ?"
	result, err := s.DB.Exec(query, id, claims.UserID)
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
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	status := vars["status"]

	query := `SELECT id, user_id, company, position, status, applied_date, notes, 
			  contact_info, salary, location, created_at, updated_at 
			  FROM job_applications WHERE status = ? AND user_id = ? ORDER BY applied_date DESC`

	rows, err := s.DB.Query(query, status, claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching jobs")
		return
	}
	defer rows.Close()

	jobs := []JobApplication{}
	for rows.Next() {
		var job JobApplication
		err := rows.Scan(&job.ID, &job.UserID, &job.Company, &job.Position, &job.Status,
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

// handleSearchJobs searches job applications with multiple filters
func (s *Server) handleSearchJobs(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse search parameters
	company := r.URL.Query().Get("company")
	position := r.URL.Query().Get("position")
	status := r.URL.Query().Get("status")
	location := r.URL.Query().Get("location")
	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Build dynamic search query
	baseQuery := "FROM job_applications WHERE user_id = ?"
	args := []interface{}{claims.UserID}

	if company != "" {
		baseQuery += " AND company LIKE ?"
		args = append(args, "%"+company+"%")
	}
	if position != "" {
		baseQuery += " AND position LIKE ?"
		args = append(args, "%"+position+"%")
	}
	if status != "" {
		baseQuery += " AND status = ?"
		args = append(args, status)
	}
	if location != "" {
		baseQuery += " AND location LIKE ?"
		args = append(args, "%"+location+"%")
	}

	// Get total count
	var totalCount int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := s.DB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error counting jobs")
		return
	}

	// Get paginated results
	searchQuery := `SELECT id, user_id, company, position, status, applied_date, notes, 
					contact_info, salary, location, created_at, updated_at ` +
		baseQuery + " ORDER BY applied_date DESC LIMIT ? OFFSET ?"
	
	searchArgs := append(args, limit, offset)
	rows, err := s.DB.Query(searchQuery, searchArgs...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error searching jobs")
		return
	}
	defer rows.Close()

	jobs := []JobApplication{}
	for rows.Next() {
		var job JobApplication
		err := rows.Scan(&job.ID, &job.UserID, &job.Company, &job.Position, &job.Status,
			&job.AppliedDate, &job.Notes, &job.ContactInfo, &job.Salary,
			&job.Location, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning job")
			return
		}
		jobs = append(jobs, job)
	}

	totalPages := (totalCount + limit - 1) / limit
	response := PaginatedResponse{
		Data:       jobs,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// parseIntQuery parses an integer from query parameters with a default value
func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
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
