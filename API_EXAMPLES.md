# API Testing Examples

This file contains sample requests you can use to test the Job Tracker API.

## Using PowerShell (Invoke-RestMethod)

### 1. Create Multiple Job Applications

```powershell
# Job 1 - Google
$job1 = @{
    company = "Google"
    position = "Software Engineer"
    status = "applied"
    applied_date = "2026-05-01T10:00:00Z"
    notes = "Applied through referral"
    contact_info = "recruiter@google.com"
    salary = "$150,000 - $180,000"
    location = "Mountain View, CA"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Post -Body $job1 -ContentType "application/json"

# Job 2 - Microsoft
$job2 = @{
    company = "Microsoft"
    position = "Cloud Engineer"
    status = "interview"
    applied_date = "2026-04-28T09:00:00Z"
    notes = "First round completed, waiting for feedback"
    contact_info = "hiring@microsoft.com"
    salary = "$140,000 - $170,000"
    location = "Seattle, WA"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Post -Body $job2 -ContentType "application/json"

# Job 3 - Amazon
$job3 = @{
    company = "Amazon"
    position = "DevOps Engineer"
    status = "offer"
    applied_date = "2026-04-15T14:30:00Z"
    notes = "Received offer, negotiating salary"
    contact_info = "talent@amazon.com"
    salary = "$145,000 + stock options"
    location = "Austin, TX"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Post -Body $job3 -ContentType "application/json"

# Job 4 - Netflix
$job4 = @{
    company = "Netflix"
    position = "Backend Developer"
    status = "rejected"
    applied_date = "2026-04-10T11:00:00Z"
    notes = "Not selected after technical interview"
    contact_info = "careers@netflix.com"
    salary = "$160,000"
    location = "Los Gatos, CA"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Post -Body $job4 -ContentType "application/json"
```

### 2. Get All Jobs

```powershell
$jobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Get
$jobs | Format-Table -Property id, company, position, status, salary
```

### 3. Get a Specific Job (replace {id} with actual ID)

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/1" -Method Get
```

### 4. Get Jobs by Status

```powershell
# Get all jobs with "applied" status
$appliedJobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/status/applied" -Method Get
$appliedJobs | Format-Table

# Get all jobs with "interview" status
$interviewJobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/status/interview" -Method Get
$interviewJobs | Format-Table

# Get all jobs with "offer" status
$offerJobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/status/offer" -Method Get
$offerJobs | Format-Table
```

### 5. Update a Job (replace {id} with actual ID)

```powershell
# Update status to interview
$update1 = @{
    status = "interview"
    notes = "Phone screen scheduled for next Tuesday"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/1" -Method Put -Body $update1 -ContentType "application/json"

# Update multiple fields
$update2 = @{
    status = "offer"
    notes = "Received offer letter"
    salary = "$155,000 + 10k sign-on bonus"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/2" -Method Put -Body $update2 -ContentType "application/json"
```

### 6. Delete a Job (replace {id} with actual ID)

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/4" -Method Delete
```

---

## Using cURL

### Create a Job

```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Tesla",
    "position": "Electrical Engineer",
    "status": "applied",
    "applied_date": "2026-05-08T10:00:00Z",
    "notes": "Exciting opportunity in automotive tech",
    "contact_info": "hr@tesla.com",
    "salary": "$120,000 - $150,000",
    "location": "Palo Alto, CA"
  }'
```

### Get All Jobs

```bash
curl http://localhost:8080/api/jobs
```

### Get a Specific Job

```bash
curl http://localhost:8080/api/jobs/1
```

### Get Jobs by Status

```bash
curl http://localhost:8080/api/jobs/status/applied
```

### Update a Job

```bash
curl -X PUT http://localhost:8080/api/jobs/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "interview",
    "notes": "Technical interview scheduled"
  }'
```

### Delete a Job

```bash
curl -X DELETE http://localhost:8080/api/jobs/1
```

---

## Status Values

Common status values to use:
- `applied` - Initial application submitted
- `interview` - In interview process
- `offer` - Offer received
- `rejected` - Application rejected
- `withdrawn` - Application withdrawn
- `accepted` - Offer accepted

---

## Complete Testing Workflow

```powershell
# 1. Create a new job application
$newJob = @{
    company = "Stripe"
    position = "Senior Backend Engineer"
    status = "applied"
    applied_date = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")
    notes = "Applied through company website"
    contact_info = "jobs@stripe.com"
    salary = "$165,000"
    location = "Remote"
} | ConvertTo-Json

$created = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Post -Body $newJob -ContentType "application/json"
Write-Host "Created job with ID: $($created.id)"

# 2. Get all jobs to verify creation
$allJobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method Get
Write-Host "Total jobs: $($allJobs.Count)"

# 3. Update the job status
$update = @{
    status = "interview"
    notes = "Phone screen completed successfully"
} | ConvertTo-Json

$updated = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/$($created.id)" -Method Put -Body $update -ContentType "application/json"
Write-Host "Updated status to: $($updated.status)"

# 4. Get jobs by status
$interviewJobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/status/interview" -Method Get
Write-Host "Jobs in interview status: $($interviewJobs.Count)"

# 5. Optionally delete
# Invoke-RestMethod -Uri "http://localhost:8080/api/jobs/$($created.id)" -Method Delete
```
