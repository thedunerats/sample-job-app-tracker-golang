# Quick Start Guide

## Prerequisites

1. **Install Go** (version 1.21 or later)
   - Download from: https://go.dev/dl/
   - Follow the installation instructions for Windows
   - Verify installation: `go version`

2. **Install GCC** (required for SQLite)
   - Download TDM-GCC from: https://jmeubank.github.io/tdm-gcc/
   - Or install via MSYS2/MinGW

## Steps to Run

### 1. Install Dependencies
```powershell
cd "c:\Git Repos\sample-job-app-tracker-golang\server"
go mod download
```

### 2. Run Tests
```powershell
cd src
go test -v
```

You should see output like:
```
=== RUN   TestInitDB
--- PASS: TestInitDB (0.01s)
=== RUN   TestCreateJob
--- PASS: TestCreateJob (0.00s)
...
PASS
ok      job-tracker     0.123s
```

### 3. Start the Server
```powershell
cd src
go run .
```

You should see:
```
Database initialized successfully
Server starting on http://localhost:8080
API endpoints:
  GET    /api/jobs              - Get all job applications
  POST   /api/jobs              - Create a new job application
  ...
```

### 4. Test the API

Open a new PowerShell window and try:

**Create a job:**
```powershell
$body = @{
    company = "Google"
    position = "Software Engineer"
    status = "applied"
    applied_date = "2026-05-08T10:00:00Z"
    notes = "Great opportunity"
    salary = "$150k"
    location = "Remote"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Post -Body $body -ContentType "application/json"
```

**Get all jobs:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Get
```

**Update a job (replace 1 with actual ID):**
```powershell
$updateBody = @{
    status = "interview"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/1" -Method Put -Body $updateBody -ContentType "application/json"
```

**Delete a job (replace 1 with actual ID):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/1" -Method Delete
```

## Troubleshooting

### "gcc: command not found" Error
Install GCC compiler as mentioned in prerequisites.

### "go: command not found"
Make sure Go is installed and added to your PATH environment variable.

### Port Already in Use
Change the port by setting the PORT environment variable:
```powershell
$env:PORT = "3000"
cd src
go run .
```

## Next Steps

- Review the [README.md](README.md) for detailed API documentation
- Explore the code in `server/src/`
- Run individual tests: `go test -v -run TestName`
- Build an executable: `go build -o job-tracker.exe`
