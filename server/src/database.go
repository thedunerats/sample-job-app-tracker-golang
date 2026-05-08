package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Database wraps the SQL database connection
type Database struct {
	*sql.DB
}

// InitDB initializes the SQLite database and creates the necessary tables
func InitDB(filepath string) (*Database, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	database := &Database{db}
	if err = database.createTables(); err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// createTables creates the necessary database tables
func (db *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS job_applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		company TEXT NOT NULL,
		position TEXT NOT NULL,
		status TEXT NOT NULL,
		applied_date DATETIME NOT NULL,
		notes TEXT,
		contact_info TEXT,
		salary TEXT,
		location TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_company ON job_applications(company);
	CREATE INDEX IF NOT EXISTS idx_status ON job_applications(status);
	CREATE INDEX IF NOT EXISTS idx_applied_date ON job_applications(applied_date);
	`

	_, err := db.Exec(query)
	return err
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}
