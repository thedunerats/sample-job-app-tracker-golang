# Job Application Tracker - Golang Backend

A RESTful API server built with Go and SQLite for tracking job applications.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete job applications
- **SQLite Database**: Lightweight, file-based database for data persistence
- **RESTful API**: Clean HTTP endpoints following REST principles
- **Comprehensive Testing**: Unit tests for all major functionality
- **Status Filtering**: Query jobs by application status

## Project Structure

```
server/
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
└── src/
    ├── main.go         # Application entry point
    ├── database.go     # Database initialization and connection
    ├── models.go       # Data structures and types
    ├── handlers.go     # HTTP request handlers
    ├── database_test.go    # Database tests
    └── handlers_test.go    # API handler tests
```

## Installation

1. Make sure you have Go 1.21 or later installed:
   ```bash
   go version
   ```

2. Navigate to the server directory:
   ```bash
   cd server
   ```

3. Download dependencies:
   ```bash
   go mod download
   ```

## Running the Server

From the `server` directory:

```bash
cd src
go run .
```

The server will start on `http://localhost:8080` by default.

### Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: Database file path (default: ./job_tracker.db)

Example:
```bash
set PORT=3000
set DB_PATH=./my_jobs.db
cd src
go run .
```

## API Endpoints

### Get All Jobs
```
GET /api/jobs
```

### Create a New Job
```
POST /api/jobs
Content-Type: application/json

{
  "company": "Google",
  "position": "Software Engineer",
  "status": "applied",
  "applied_date": "2026-05-08T10:00:00Z",
  "notes": "Exciting opportunity",
  "contact_info": "recruiter@google.com",
  "salary": "$150,000",
  "location": "Mountain View, CA"
}
```

### Get a Specific Job
```
GET /api/jobs/{id}
```

### Update a Job
```
PUT /api/jobs/{id}
Content-Type: application/json

{
  "status": "interview",
  "notes": "Phone screen scheduled for next week"
}
```

### Delete a Job
```
DELETE /api/jobs/{id}
```

### Get Jobs by Status
```
GET /api/jobs/status/{status}
```

Common status values: `applied`, `interview`, `offer`, `rejected`

## Running Tests

From the `server` directory:

```bash
cd src
go test -v
```

To run specific tests:
```bash
# Run only database tests
go test -v -run TestDatabase

# Run only handler tests
go test -v -run TestHandle
```

To see test coverage:
```bash
go test -cover
```

## Database Schema

The `job_applications` table includes:

| Column       | Type     | Description                    |
|-------------|----------|--------------------------------|
| id          | INTEGER  | Primary key (auto-increment)   |
| company     | TEXT     | Company name                   |
| position    | TEXT     | Job position/title             |
| status      | TEXT     | Application status             |
| applied_date| DATETIME | Date application was submitted |
| notes       | TEXT     | Additional notes               |
| contact_info| TEXT     | Recruiter/contact information  |
| salary      | TEXT     | Salary information             |
| location    | TEXT     | Job location                   |
| created_at  | DATETIME | Record creation timestamp      |
| updated_at  | DATETIME | Last update timestamp          |

## Example Usage with cURL

### Create a job application:
```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d "{\"company\":\"Microsoft\",\"position\":\"Cloud Engineer\",\"status\":\"applied\",\"applied_date\":\"2026-05-08T10:00:00Z\",\"salary\":\"$140k\",\"location\":\"Seattle\"}"
```

### Get all jobs:
```bash
curl http://localhost:8080/api/jobs
```

### Update a job:
```bash
curl -X PUT http://localhost:8080/api/jobs/1 \
  -H "Content-Type: application/json" \
  -d "{\"status\":\"interview\"}"
```

### Delete a job:
```bash
curl -X DELETE http://localhost:8080/api/jobs/1
```

## Testing Strategy

The test suite includes:

1. **Database Tests** (`database_test.go`):
   - Database initialization
   - CRUD operations at the database level
   - Query filtering by status

2. **Handler Tests** (`handlers_test.go`):
   - HTTP endpoint testing
   - Request/response validation
   - Error handling
   - Complete workflow testing

All tests use temporary SQLite databases that are automatically cleaned up after execution.

## Future Enhancements

Potential improvements:
- Authentication and authorization
- Pagination for large datasets
- Search functionality
- Export to CSV/PDF
- Email notifications for application deadlines
- Integration with job boards
- Frontend web interface

## Dependencies

- **github.com/mattn/go-sqlite3**: SQLite driver for Go
- **github.com/gorilla/mux**: HTTP router and URL matcher

## License

This is a sample project for educational purposes.
