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
	// Check if users table exists
	var usersExists int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&usersExists)

	// Check if job_applications has user_id column
	var hasUserID int
	db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('job_applications') WHERE name='user_id'").Scan(&hasUserID)

	// If tables need migration, drop and recreate
	if usersExists == 0 || hasUserID == 0 {
		log.Println("Migrating database schema...")
		_, err := db.Exec("DROP TABLE IF EXISTS job_applications; DROP TABLE IF EXISTS users;")
		if err != nil {
			return err
		}
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS job_applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		company TEXT NOT NULL,
		position TEXT NOT NULL,
		status TEXT NOT NULL,
		applied_date DATETIME NOT NULL,
		notes TEXT,
		contact_info TEXT,
		salary TEXT,
		location TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_user_id ON job_applications(user_id);
	CREATE INDEX IF NOT EXISTS idx_company ON job_applications(company);
	CREATE INDEX IF NOT EXISTS idx_position ON job_applications(position);
	CREATE INDEX IF NOT EXISTS idx_status ON job_applications(status);
	CREATE INDEX IF NOT EXISTS idx_location ON job_applications(location);
	CREATE INDEX IF NOT EXISTS idx_applied_date ON job_applications(applied_date);
	CREATE INDEX IF NOT EXISTS idx_email ON users(email);
	`

	_, err := db.Exec(query)
	return err
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}
