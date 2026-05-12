# Job Application Tracker

A full-stack application for tracking job applications with a **Go backend** and **Angular frontend**.

## Features

### Backend (Go + SQLite)
- **Authentication**: JWT-based authentication with password hashing (bcrypt)
- **Authorization**: User-specific data isolation
- **CRUD Operations**: Create, Read, Update, Delete job applications
- **Search & Filter**: Search by company, position, status, location
- **Pagination**: Efficient data loading with customizable page sizes
- **CORS Support**: Configured for cross-origin requests
- **SQLite Database**: Lightweight, file-based database with auto-migration
- **Comprehensive Testing**: 29 tests organized by module with full coverage
- **Code Quality**: Verified with `go fmt` and `go vet`
- **CI/CD**: GitHub Actions workflow for automated testing

### Frontend (Angular 17)
- **Modern UI**: Responsive design with professional styling
- **Authentication**: Login and registration with form validation
- **Job Management**: Create, edit, view, and delete job applications
- **Advanced Search**: Multi-field search with filters
- **Pagination**: Navigate through large datasets
- **Type Safety**: Full TypeScript type definitions
- **Unit Tests**: Comprehensive test coverage for all services

## Project Structure

```
.
├── .github/
│   └── workflows/
│       └── test.yml           # GitHub Actions CI/CD workflow
├── client/                    # Angular frontend
│   ├── src/
│   │   ├── app/
│   │   │   ├── components/   # UI components (Login, Register, Job List, etc.)
│   │   │   ├── guards/       # Route guards for authentication
│   │   │   ├── models/       # TypeScript interfaces
│   │   │   └── services/     # HTTP services (API, Auth, Job)
│   │   ├── environments/     # Environment configurations
│   │   └── styles.scss       # Global styles
│   ├── test/                 # Unit tests
│   ├── tsconfig.json         # TypeScript configuration
│   └── package.json          # npm dependencies
└── server/                    # Go backend
    ├── go.mod                # Go module dependencies
    ├── go.sum                # Dependency checksums
    └── src/
        ├── main.go           # Application entry point
        ├── database.go       # Database initialization with auto-migration
        ├── models.go         # Data structures and types
        ├── handlers.go       # Job CRUD operations and routing
        ├── auth.go           # JWT token management and password hashing
        ├── auth_handlers.go  # Authentication endpoints (register/login)
        ├── middleware.go     # Authentication & authorization middleware
        ├── cors.go           # CORS middleware
        ├── auth_test.go      # Auth utility tests (6 tests)
        ├── auth_handlers_test.go  # Auth endpoint tests (4 tests)
        ├── database_test.go  # Database layer tests (6 tests)
        ├── handlers_test.go  # Job handler tests (11 tests)
        └── middleware_test.go # Middleware tests (2 tests)
```

## Quick Start

### Option 1: Using Docker (Recommended)

The easiest way to run the application with zero dependency installation:

**Prerequisites:**
- Docker and Docker Compose ([Download](https://www.docker.com/get-started))

**Run the entire stack:**
```bash
# Clone the repository
git clone <your-repo-url>
cd sample-job-app-tracker-golang

# Start both backend and frontend
docker compose up -d

# View logs
docker compose logs -f

# Stop the services
docker compose down
```

The application will be available at:
- **Frontend:** http://localhost (port 80)
- **Backend API:** http://localhost:8080

**What Docker does for you:**
- ✅ No Go installation required
- ✅ No Node.js/npm installation required
- ✅ No dependency conflicts
- ✅ Consistent environment across all machines
- ✅ Database persistence in Docker volumes
- ✅ Production-like environment

**Development mode with live reload:**
```bash
# Use the development compose file
docker compose -f docker-compose.dev.yml up

# Frontend: http://localhost:4200
# Backend: http://localhost:8080
# Code changes automatically reload
```

### Option 2: Manual Installation

**Prerequisites:**

**Backend:**
- Go 1.21 or later ([Download](https://go.dev/dl/))

**Frontend:**
- Node.js 18+ and npm ([Download](https://nodejs.org/))
- Angular CLI: `npm install -g @angular/cli`

### Installation & Running

#### 1. Start the Backend Server

```powershell
# Navigate to server directory
cd server

# Download Go dependencies
go mod download

# Run tests to verify everything works
cd src
go test -v

# Start the server
go run .
```

Server will start on `http://localhost:8080`

#### 2. Start the Frontend Application

```powershell
# Navigate to client directory (in a new terminal)
cd client

# Install npm dependencies
npm install

# Start the development server
ng serve
```

Frontend will be available at `http://localhost:4200`

### First Time Setup

1. Open your browser to `http://localhost:4200`
2. Click "Register" to create a new account
3. Fill in name, email, and password (minimum 6 characters)
4. Start tracking your job applications!

## API Documentation

### Authentication Endpoints

#### Register a New User
```http
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe"
  }
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secure123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe"
  }
}
```

#### Get Current User
```http
GET /api/auth/me
Authorization: Bearer {token}

Response:
{
  "id": 1,
  "email": "john@example.com",
  "name": "John Doe",
  "created_at": "2026-05-08T10:00:00Z"
}
```

### Job Application Endpoints

**Note:** All job endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer {your-token-here}
```

#### Get All Jobs (with Pagination)
```http
GET /api/jobs?page=1&limit=10
Authorization: Bearer {token}

Response:
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "company": "Google",
      "position": "Software Engineer",
      "status": "applied",
      "applied_date": "2026-05-08T10:00:00Z",
      "location": "Mountain View, CA",
      "salary": "$150,000",
      "contact_info": "recruiter@google.com",
      "notes": "Great opportunity",
      "created_at": "2026-05-08T10:00:00Z",
      "updated_at": "2026-05-08T10:00:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "limit": 10,
  "total_pages": 3
}
```

#### Search Jobs
```http
GET /api/jobs/search?company=Google&status=applied&position=Engineer&location=Remote
Authorization: Bearer {token}

# Supports filters: company, position, status, location
# All filters are optional and use partial matching (LIKE)
```

#### Create a New Job
```http
POST /api/jobs
Authorization: Bearer {token}
Content-Type: application/json

{
  "company": "Microsoft",
  "position": "Cloud Engineer",
  "status": "applied",
  "applied_date": "2026-05-08T10:00:00Z",
  "location": "Seattle, WA",
  "salary": "$140,000",
  "contact_info": "hiring@microsoft.com",
  "notes": "Excited about this role"
}
```

#### Get a Specific Job
```http
GET /api/jobs/{id}
Authorization: Bearer {token}
```

#### Update a Job
```http
PUT /api/jobs/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": "interview",
  "notes": "Phone screen completed, waiting for next round"
}
```

#### Delete a Job
```http
DELETE /api/jobs/{id}
Authorization: Bearer {token}
```

#### Get Jobs by Status
```http
GET /api/jobs/status/{status}
Authorization: Bearer {token}

# Common status values: applied, interview, offer, rejected, withdrawn, accepted
```

## API Examples

### Using PowerShell (Windows)

#### Complete Workflow Example
```powershell
# 1. Register a new user
$registerData = @{
    name = "Jane Smith"
    email = "jane@example.com"
    password = "secure456"
} | ConvertTo-Json

$authResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/register" `
    -Method Post -Body $registerData -ContentType "application/json"

$token = $authResponse.token
Write-Host "Registered! Token: $token"

# 2. Create a job application
$jobData = @{
    company = "Google"
    position = "Software Engineer"
    status = "applied"
    applied_date = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")
    location = "Mountain View, CA"
    salary = "$150,000 - $180,000"
    contact_info = "recruiter@google.com"
    notes = "Applied through referral"
} | ConvertTo-Json

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$job = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" `
    -Method Post -Body $jobData -Headers $headers

Write-Host "Created job with ID: $($job.id)"

# 3. Get all jobs (with pagination)
$allJobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs?page=1&limit=10" `
    -Method Get -Headers $headers

Write-Host "Total jobs: $($allJobs.total)"
$allJobs.data | Format-Table -Property id, company, position, status

# 4. Search for jobs
$searchResults = Invoke-RestMethod `
    -Uri "http://localhost:8080/api/jobs/search?company=Google&status=applied" `
    -Method Get -Headers $headers

Write-Host "Search results: $($searchResults.total)"

# 5. Update job status
$updateData = @{
    status = "interview"
    notes = "Phone screen scheduled for next week"
} | ConvertTo-Json

$updated = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/$($job.id)" `
    -Method Put -Body $updateData -Headers $headers

Write-Host "Updated status to: $($updated.status)"

# 6. Delete a job (optional)
# Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/$($job.id)" `
#     -Method Delete -Headers $headers
```

### Using cURL (Linux/Mac/Git Bash)

```bash
# Register and save token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secure123"}' \
  | jq -r '.token')

echo "Token: $TOKEN"

# Create a job application
curl -X POST http://localhost:8080/api/jobs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Amazon",
    "position": "DevOps Engineer",
    "status": "applied",
    "applied_date": "2026-05-08T10:00:00Z",
    "location": "Austin, TX",
    "salary": "$145,000",
    "notes": "Interesting cloud role"
  }'

# Get all jobs
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/jobs?page=1&limit=10

# Search jobs
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/jobs/search?company=Amazon&status=applied"

# Update a job (replace 1 with actual ID)
curl -X PUT http://localhost:8080/api/jobs/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"interview","notes":"Technical interview next Tuesday"}'

# Delete a job (replace 1 with actual ID)
curl -X DELETE http://localhost:8080/api/jobs/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Database Schema

### Users Table
| Column       | Type     | Description                    |
|-------------|----------|--------------------------------|
| id          | INTEGER  | Primary key (auto-increment)   |
| email       | TEXT     | User email (unique)            |
| password    | TEXT     | Bcrypt hashed password         |
| name        | TEXT     | User's full name               |
| created_at  | DATETIME | Account creation timestamp     |
| updated_at  | DATETIME | Last update timestamp          |

### Job Applications Table
| Column       | Type     | Description                    |
|-------------|----------|--------------------------------|
| id          | INTEGER  | Primary key (auto-increment)   |
| user_id     | INTEGER  | Foreign key to users table     |
| company     | TEXT     | Company name                   |
| position    | TEXT     | Job position/title             |
| status      | TEXT     | Application status             |
| applied_date| DATETIME | Date application was submitted |
| location    | TEXT     | Job location                   |
| salary      | TEXT     | Salary information             |
| contact_info| TEXT     | Recruiter/contact information  |
| notes       | TEXT     | Additional notes               |
| created_at  | DATETIME | Record creation timestamp      |
| updated_at  | DATETIME | Last update timestamp          |

**Data Isolation:** Each user can only access their own job applications. The foreign key constraint ensures data integrity, and `ON DELETE CASCADE` automatically removes jobs when a user is deleted.

**Auto-Migration:** The database automatically migrates to the latest schema on startup, checking for required tables and columns.

## Testing

### Backend Tests

The Go backend has comprehensive test coverage with **29 tests** organized by module:

**Authentication Utility Tests** (`auth_test.go` - 6 tests):
- Password hashing with bcrypt
- Password verification
- JWT token generation
- JWT token validation
- Token expiration handling
- Multiple token generation

**Authentication Handler Tests** (`auth_handlers_test.go` - 4 tests):
- User registration success
- User registration validation errors
- User login success
- User login with invalid credentials

**Database Tests** (`database_test.go` - 6 tests):
- Database initialization
- Job creation (SQL level)
- Job retrieval (SQL level)
- Job updates (SQL level)
- Job deletion (SQL level)
- Filtering jobs by status (SQL level)

**Handler Tests** (`handlers_test.go` - 11 tests):
- Get all jobs with pagination
- Create job with authentication
- Create job validation (missing required fields)
- Get single job by ID
- Get job not found (404 handling)
- Update job (partial updates)
- Delete job
- Get jobs by status filter
- Pagination (multiple page sizes and offsets)
- Search functionality (company, position, status, location, combined)
- Complete workflow (create → read → update → search → delete)

**Middleware Tests** (`middleware_test.go` - 2 tests):
- Authentication requirement enforcement
- User data isolation (users can only access their own data)

**Test Organization:**
All test files follow Go conventions:
- Named `*_test.go` to match the source file they test
- Located in the same package/directory as the source code
- Example: `handlers.go` is tested by `handlers_test.go`

**Code Quality:**
- All code is formatted with `go fmt`
- All code passes `go vet` static analysis
- Zero compilation warnings or errors

Run all tests:
```powershell
cd server/src
go test -v
```

Run with coverage:
```powershell
go test -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Run with race detection:
```powershell
go test -v -race
```

Format code:
```powershell
go fmt ./...
```

Run static analysis:
```powershell
go vet ./...
```

### Frontend Tests

The Angular frontend includes unit tests for all services:

**API Service Tests** (`client/test/api.service.spec.ts`):
- HTTP GET, POST, PUT, DELETE methods
- Authorization header inclusion
- Query parameter handling

**Auth Service Tests** (`client/test/auth.service.spec.ts`):
- User registration and login
- Token management
- Current user fetching
- Logout functionality

**Job Service Tests** (`client/test/job.service.spec.ts`):
- Job CRUD operations
- Pagination
- Search with filters
- Status filtering

Run frontend tests:
```powershell
cd client
ng test
```

### CI/CD Pipeline

GitHub Actions automatically builds and tests both backend and frontend on every pull request and push to main/master branches. The workflow runs three parallel jobs:

**Backend Job:**
1. Sets up Go 1.21 environment
2. Caches Go modules for faster builds
3. Downloads dependencies
4. **Builds the backend binary** (catches compilation errors)
5. Runs all tests with race detection and coverage
6. Uploads coverage reports to Codecov

**Frontend Job:**
1. Sets up Node.js 18 environment
2. Caches npm packages for faster builds
3. Installs dependencies
4. **Builds the Angular app for production** (catches build errors)
5. Runs all unit tests in headless Chrome
6. Uploads coverage reports to Codecov

**Docker Job:**
1. Sets up Docker Buildx
2. **Builds backend Docker image** (verifies Dockerfile and dependencies)
3. **Builds frontend Docker image** (verifies nginx configuration)
4. **Tests docker-compose.yml** configuration
5. Starts both services in containers
6. **Health checks both services** (verifies they actually run)
7. Tests API connectivity
8. Shows logs if any step fails

Both jobs must pass before a pull request can be merged. This ensures:
- ✅ Code compiles successfully
- ✅ All tests pass
- ✅ No breaking changes introduced
- ✅ Production builds work

View the workflow: [.github/workflows/test.yml](.github/workflows/test.yml)

## Configuration

### Backend Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_PATH` | `./job_tracker.db` | SQLite database file path |

Example:
```powershell
$env:PORT = "3000"
$env:DB_PATH = "C:\data\jobs.db"
cd server/src
go run .
```

### Frontend Environment Configuration

**Development** (`client/src/environments/environment.ts`):
```typescript
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api'
};
```

**Production** (`client/src/environments/environment.prod.ts`):
```typescript
export const environment = {
  production: true,
  apiUrl: '/api'  // Assumes backend served from same domain
};
```

Build for production:
```powershell
cd client
ng build --configuration production
```

### Authentication Configuration

**JWT Secret** (⚠️ Security Note):
The JWT secret is currently hardcoded in `server/src/auth.go` as `"your-secret-key-change-this-in-production"`.

**For production, update this to use an environment variable:**
```go
// In auth.go
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
```

**Token Expiration:** JWT tokens expire after 24 hours. Users must log in again after expiration.

**Password Requirements:**
- Minimum length: 6 characters
- Hashing: bcrypt with default cost (10)

### CORS Configuration

CORS is configured in `server/src/cors.go` to allow all origins (`Access-Control-Allow-Origin: *`).

**For production, restrict to specific origins:**
```go
// In cors.go, replace:
w.Header().Set("Access-Control-Allow-Origin", "*")

// With:
origin := r.Header.Get("Origin")
if origin == "https://yourdomain.com" {
    w.Header().Set("Access-Control-Allow-Origin", origin)
}
```

## Deployment

### Option 1: Docker Deployment (Recommended)

Docker provides the easiest and most consistent deployment method.

**Production deployment:**
```bash
# Pull the latest code
git pull

# Build and start containers
docker compose up -d --build

# Check status
docker compose ps

# View logs
docker compose logs -f

# Stop services
docker compose down
```

**Environment variables for production:**
Create a `.env` file in the project root:
```env
# Backend
PORT=8080
DB_PATH=/home/appuser/data/job_tracker.db
JWT_SECRET=your-very-secure-random-secret-key-here-min-32-chars

# Frontend (if needed)
NODE_ENV=production
```

**Deploy to cloud platforms:**

*AWS EC2 / Google Cloud Compute / Azure VM:*
1. Install Docker on the VM
2. Clone repository
3. Run `docker compose up -d`
4. Configure firewall to allow ports 80 and 8080

*Docker Swarm / Kubernetes:*
Use the Dockerfiles as base images and create appropriate manifests.

**Database persistence:**
The database is automatically persisted in a Docker volume named `backend-data`. To backup:
```bash
# Backup database
docker compose exec backend tar -czf /tmp/backup.tar.gz /home/appuser/data
docker compose cp backend:/tmp/backup.tar.gz ./backup.tar.gz

# Restore database
docker compose cp ./backup.tar.gz backend:/tmp/backup.tar.gz
docker compose exec backend tar -xzf /tmp/backup.tar.gz -C /
docker compose restart backend
```

### Option 2: Manual Deployment

#### Backend Deployment

**Build executable:**
```bash
cd server/src
# Linux/Mac
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o job-tracker .

# Windows
go build -o job-tracker.exe .
```

**Run in production:**
```bash
export PORT=8080
export DB_PATH=/var/lib/job-tracker/jobs.db
export JWT_SECRET=your-very-secure-random-secret-key-here
./job-tracker
```

**Create systemd service (Linux):**
```ini
# /etc/systemd/system/job-tracker.service
[Unit]
Description=Job Tracker API Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/job-tracker
Environment="PORT=8080"
Environment="DB_PATH=/var/lib/job-tracker/jobs.db"
Environment="JWT_SECRET=your-secret-key"
ExecStart=/opt/job-tracker/job-tracker
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable job-tracker
sudo systemctl start job-tracker
```

#### Frontend Deployment

**Build for production:**
```bash
cd client
ng build --configuration production
```

The build artifacts will be in `client/dist/client/`. Serve these static files with any web server (nginx, Apache, IIS, etc.).

**Important:** Configure your web server to:
1. Serve `index.html` for all routes (for Angular routing)
2. Proxy `/api/*` requests to the Go backend

Example nginx configuration:
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    root /var/www/job-tracker/client/dist/client;
    index index.html;
    
    # Serve Angular app
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # Proxy API requests to Go backend
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Docker Image Management

**Build images individually:**
```bash
# Backend
docker build -t job-tracker-backend:latest ./server

# Frontend
docker build -t job-tracker-frontend:latest ./client
```

**Push to Docker registry:**
```bash
# Tag images
docker tag job-tracker-backend:latest your-registry/job-tracker-backend:latest
docker tag job-tracker-frontend:latest your-registry/job-tracker-frontend:latest

# Push
docker push your-registry/job-tracker-backend:latest
docker push your-registry/job-tracker-frontend:latest
```

**Image sizes:**
- Backend: ~15-20 MB (multi-stage build with Alpine Linux)
- Frontend: ~25-30 MB (nginx-alpine serving static files)

## Troubleshooting

### Docker Issues

**Container won't start:**
```bash
# Check logs
docker compose logs backend
docker compose logs frontend

# Check container status
docker compose ps

# Rebuild from scratch
docker compose down -v
docker compose build --no-cache
docker compose up -d
```

**Port already in use:**
```bash
# Check what's using the port
sudo lsof -i :8080  # Linux/Mac
netstat -ano | findstr :8080  # Windows

# Stop existing containers
docker compose down
```

**Database volume issues:**
```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect sample-job-app-tracker-golang_backend-data

# Remove volume (WARNING: deletes all data)
docker compose down -v
```

### Backend Issues

**Port already in use:**
```powershell
# Windows: Kill process on port 8080
Get-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess | Stop-Process -Force

# Then restart server
cd server/src
go run .
```

**Database locked error:**
- Close any other applications accessing the database file
- Ensure only one server instance is running
- Check file permissions on the database file

**Tests failing:**
```powershell
# Clean build cache and rerun
go clean -testcache
go test -v
```

### Frontend Issues

**Registration/Login fails with CORS error:**
- Ensure backend server is running on port 8080
- Check browser console for specific CORS errors
- Verify CORS middleware is enabled in backend

**"Module not found" errors:**
```powershell
# Delete node_modules and reinstall
rm -r node_modules
rm package-lock.json
npm install
```

**Angular build errors:**
```powershell
# Clear Angular cache
ng cache clean

# Rebuild
ng build
```

**Port 4200 already in use:**
```powershell
# Run on different port
ng serve --port 4201
```

### Authentication Issues

**"Invalid token" or "Unauthorized" errors:**
- Token may have expired (24-hour lifetime)
- Log out and log back in to get a new token
- Check that token is being included in Authorization header

**Password validation fails:**
- Passwords must be at least 6 characters
- Check for whitespace or special characters if issues persist

### Database Issues

**"No such table" errors:**
- Delete the database file and restart the server (auto-migration will recreate)
```powershell
rm server/src/job_tracker.db
cd server/src
go run .
```

**Schema migration issues:**
- The app auto-migrates on startup
- Check logs for migration messages
- If problems persist, backup data and delete database to start fresh

## Version Control

### What to Commit

**DO commit:**
- `go.mod` and `go.sum` (required for reproducible builds)
- `package.json` and `package-lock.json`
- All source code files
- Configuration files (tsconfig.json, angular.json, etc.)
- **Docker files** (Dockerfile, docker-compose.yml, .dockerignore, nginx.conf)
- GitHub Actions workflows

**DO NOT commit:**
- `*.db` files (database files with user data)
- `node_modules/` directory
- `dist/` build outputs
- Docker volumes and containers
- IDE-specific files (.vscode/, .idea/)
- Environment-specific secrets (.env files)

See `.gitignore` and `.dockerignore` for complete exclusion lists.

## Technology Stack

### Backend
- **Language:** Go 1.21
- **Database:** SQLite (modernc.org/sqlite - pure Go driver, no CGO required)
- **Router:** gorilla/mux v1.8.1
- **Authentication:** golang-jwt/jwt/v5 v5.2.0
- **Password Hashing:** golang.org/x/crypto (bcrypt)

### Frontend
- **Framework:** Angular 17.3.14
- **Language:** TypeScript
- **Styling:** SCSS
- **Architecture:** Standalone Components
- **HTTP Client:** Angular HttpClient
- **Forms:** Reactive Forms
- **Web Server:** nginx-alpine (production)

### Infrastructure & DevOps
- **Containerization:** Docker & Docker Compose
- **Base Images:** Alpine Linux (minimal size)
- **CI/CD:** GitHub Actions
- **Testing:** Go testing package, Jasmine/Karma
- **Code Coverage:** Codecov

### Key Features
- **Zero External Dependencies:** Pure Go SQLite driver means no C compiler needed
- **Multi-stage Docker Builds:** Optimized images (~15MB backend, ~25MB frontend)
- **Container Health Checks:** Automated health monitoring
- **Database Persistence:** Docker volumes for data safety
- **Production Ready:** nginx with proper caching and security headers

## License

This project is provided as-is for educational and personal use.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

All pull requests will trigger the automated test suite via GitHub Actions.

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Dependencies

- **github.com/mattn/go-sqlite3**: SQLite driver for Go
- **github.com/gorilla/mux**: HTTP router and URL matcher

## License

This is a sample project for educational purposes.
