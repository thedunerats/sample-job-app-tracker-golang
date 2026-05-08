package main

import (
	"time"
)

// JobApplication represents a job application entry
type JobApplication struct {
	ID          int       `json:"id"`
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
