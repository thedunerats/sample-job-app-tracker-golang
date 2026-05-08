package main

import (
	"time"
)

// User represents a user account
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never include password in JSON responses
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JobApplication represents a job application entry
type JobApplication struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	Status      string    `json:"status"` // e.g., "applied", "interview", "offer", "rejected"
	AppliedDate time.Time `json:"applied_date"`
	Notes       string    `json:"notes"`
	ContactInfo string    `json:"contact_info"`
	Salary      string    `json:"salary"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateJobApplicationRequest represents the request body for creating a job application
type CreateJobApplicationRequest struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	Status      string `json:"status"`
	AppliedDate string `json:"applied_date"` // ISO 8601 format
	Notes       string `json:"notes"`
	ContactInfo string `json:"contact_info"`
	Salary      string `json:"salary"`
	Location    string `json:"location"`
}

// UpdateJobApplicationRequest represents the request body for updating a job application
type UpdateJobApplicationRequest struct {
	Company     *string `json:"company,omitempty"`
	Position    *string `json:"position,omitempty"`
	Status      *string `json:"status,omitempty"`
	AppliedDate *string `json:"applied_date,omitempty"`
	Notes       *string `json:"notes,omitempty"`
	ContactInfo *string `json:"contact_info,omitempty"`
	Salary      *string `json:"salary,omitempty"`
	Location    *string `json:"location,omitempty"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// SearchParams represents search query parameters
type SearchParams struct {
	Company  string
	Position string
	Status   string
	Location string
	Page     int
	Limit    int
}
