package main

import (
	"os"
	"testing"
	"time"
)

// setupTestDB creates a temporary test database
func setupTestDB(t *testing.T) *Database {
	dbPath := "./test_job_tracker.db"
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Clean up function
	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return db
}

// createTestUserDB creates a test user directly in the database and returns the user ID
func createTestUserDB(t *testing.T, db *Database) int {
	query := `INSERT INTO users (email, password, name, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := db.Exec(query, "test@example.com", "hashedpassword", "Test User", now, now)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	id, _ := result.LastInsertId()
	return int(id)
}

func TestInitDB(t *testing.T) {
	dbPath := "./test_init.db"
	defer os.Remove(dbPath)

	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Verify table was created
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name='job_applications'"
	var tableName string
	err = db.QueryRow(query).Scan(&tableName)
	if err != nil {
		t.Fatalf("Table job_applications was not created: %v", err)
	}

	if tableName != "job_applications" {
		t.Errorf("Expected table name 'job_applications', got '%s'", tableName)
	}
}

func TestCreateJob(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUserDB(t, db)

	// Insert a job application
	query := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
			  notes, contact_info, salary, location, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := db.Exec(query, userID, "Google", "Software Engineer", "applied", now,
		"Great opportunity", "hr@google.com", "$150k", "Remote", now, now)

	if err != nil {
		t.Fatalf("Failed to insert job application: %v", err)
	}

	id, _ := result.LastInsertId()
	if id == 0 {
		t.Error("Expected non-zero ID after insert")
	}
}

func TestGetJob(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUserDB(t, db)

	// Insert a job application first
	insertQuery := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
				   notes, contact_info, salary, location, created_at, updated_at) 
				   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, _ := db.Exec(insertQuery, userID, "Microsoft", "Cloud Engineer", "interview", now,
		"Second round", "recruiter@ms.com", "$140k", "Seattle", now, now)

	id, _ := result.LastInsertId()

	// Retrieve the job application
	selectQuery := `SELECT company, position, status FROM job_applications WHERE id = ?`
	var company, position, status string
	err := db.QueryRow(selectQuery, id).Scan(&company, &position, &status)

	if err != nil {
		t.Fatalf("Failed to retrieve job application: %v", err)
	}

	if company != "Microsoft" {
		t.Errorf("Expected company 'Microsoft', got '%s'", company)
	}
	if position != "Cloud Engineer" {
		t.Errorf("Expected position 'Cloud Engineer', got '%s'", position)
	}
	if status != "interview" {
		t.Errorf("Expected status 'interview', got '%s'", status)
	}
}

func TestUpdateJob(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUserDB(t, db)

	// Insert a job application
	insertQuery := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
				   notes, contact_info, salary, location, created_at, updated_at) 
				   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, _ := db.Exec(insertQuery, userID, "Amazon", "DevOps Engineer", "applied", now,
		"Pending response", "hiring@amazon.com", "$130k", "Austin", now, now)

	id, _ := result.LastInsertId()

	// Update the status
	updateQuery := "UPDATE job_applications SET status = ?, updated_at = ? WHERE id = ?"
	_, err := db.Exec(updateQuery, "offer", time.Now(), id)

	if err != nil {
		t.Fatalf("Failed to update job application: %v", err)
	}

	// Verify the update
	selectQuery := "SELECT status FROM job_applications WHERE id = ?"
	var status string
	db.QueryRow(selectQuery, id).Scan(&status)

	if status != "offer" {
		t.Errorf("Expected status 'offer', got '%s'", status)
	}
}

func TestDeleteJob(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUserDB(t, db)

	// Insert a job application
	insertQuery := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
				   notes, contact_info, salary, location, created_at, updated_at) 
				   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, _ := db.Exec(insertQuery, userID, "Apple", "iOS Developer", "rejected", now,
		"Not a good fit", "jobs@apple.com", "$145k", "Cupertino", now, now)

	id, _ := result.LastInsertId()

	// Delete the job application
	deleteQuery := "DELETE FROM job_applications WHERE id = ?"
	deleteResult, err := db.Exec(deleteQuery, id)

	if err != nil {
		t.Fatalf("Failed to delete job application: %v", err)
	}

	rowsAffected, _ := deleteResult.RowsAffected()
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	// Verify deletion
	selectQuery := "SELECT id FROM job_applications WHERE id = ?"
	var retrievedID int
	err = db.QueryRow(selectQuery, id).Scan(&retrievedID)

	if err == nil {
		t.Error("Expected error when retrieving deleted job, but got none")
	}
}

func TestDatabaseGetJobsByStatus(t *testing.T) {
	db := setupTestDB(t)
	userID := createTestUserDB(t, db)

	// Insert multiple job applications with different statuses
	insertQuery := `INSERT INTO job_applications (user_id, company, position, status, applied_date, 
				   notes, contact_info, salary, location, created_at, updated_at) 
				   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	db.Exec(insertQuery, userID, "Company A", "Engineer A", "applied", now, "", "", "", "", now, now)
	db.Exec(insertQuery, userID, "Company B", "Engineer B", "interview", now, "", "", "", "", now, now)
	db.Exec(insertQuery, userID, "Company C", "Engineer C", "applied", now, "", "", "", "", now, now)
	db.Exec(insertQuery, userID, "Company D", "Engineer D", "offer", now, "", "", "", "", now, now)

	// Query jobs with "applied" status
	selectQuery := "SELECT COUNT(*) FROM job_applications WHERE status = ?"
	var count int
	err := db.QueryRow(selectQuery, "applied").Scan(&count)

	if err != nil {
		t.Fatalf("Failed to count jobs by status: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 jobs with status 'applied', got %d", count)
	}
}
