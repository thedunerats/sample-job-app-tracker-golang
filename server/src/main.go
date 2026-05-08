package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./job_tracker.db"
	}

	// Initialize database
	db, err := InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create server
	server := NewServer(db)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("API endpoints:")
	log.Printf("  GET    /health                - Health check")
	log.Printf("  POST   /api/auth/register     - Register new user")
	log.Printf("  POST   /api/auth/login        - Login user")
	log.Printf("  GET    /api/jobs              - Get all job applications")
	log.Printf("  POST   /api/jobs              - Create a new job application")
	log.Printf("  GET    /api/jobs/{id}         - Get a specific job application")
	log.Printf("  PUT    /api/jobs/{id}         - Update a job application")
	log.Printf("  DELETE /api/jobs/{id}         - Delete a job application")
	log.Printf("  GET    /api/jobs/status/{status} - Get jobs by status")

	if err := http.ListenAndServe(addr, server.Router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
