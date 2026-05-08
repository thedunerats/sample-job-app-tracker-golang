import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { JobService } from '../src/app/services/job.service';
import { ApiService } from '../src/app/services/api.service';
import { JobApplication, PaginatedResponse, SearchParams } from '../src/app/models/job.model';

describe('JobService', () => {
  let service: JobService;
  let httpMock: HttpTestingController;
  let apiService: ApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [JobService, ApiService]
    });
    service = TestBed.inject(JobService);
    apiService = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get all jobs with pagination', (done) => {
    const mockResponse: PaginatedResponse<JobApplication> = {
      data: [
        {
          id: 1,
          company: 'Google',
          position: 'Software Engineer',
          status: 'Applied',
          applied_date: new Date().toISOString()
        }
      ],
      page: 1,
      limit: 10,
      total_count: 1,
      total_pages: 1
    };

    service.getAllJobs(1, 10).subscribe(response => {
      expect(response).toEqual(mockResponse);
      expect(response.data.length).toBe(1);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/jobs?page=1&limit=10');
    expect(req.request.method).toBe('GET');
    req.flush(mockResponse);
  });

  it('should get a single job', (done) => {
    const mockJob: JobApplication = {
      id: 1,
      company: 'Google',
      position: 'Software Engineer',
      status: 'Applied',
      applied_date: new Date().toISOString()
    };

    service.getJob(1).subscribe(job => {
      expect(job).toEqual(mockJob);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/jobs/1');
    expect(req.request.method).toBe('GET');
    req.flush(mockJob);
  });

  it('should create a job', (done) => {
    const newJob: JobApplication = {
      company: 'Google',
      position: 'Software Engineer',
      status: 'Applied',
      applied_date: new Date().toISOString()
    };

    const createdJob: JobApplication = {
      id: 1,
      ...newJob
    };

    service.createJob(newJob).subscribe(job => {
      expect(job).toEqual(createdJob);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/jobs');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(newJob);
    req.flush(createdJob);
  });

  it('should update a job', (done) => {
    const updatedJob: JobApplication = {
      id: 1,
      company: 'Google',
      position: 'Senior Software Engineer',
      status: 'Interview',
      applied_date: new Date().toISOString()
    };

    service.updateJob(1, updatedJob).subscribe(job => {
      expect(job).toEqual(updatedJob);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/jobs/1');
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(updatedJob);
    req.flush(updatedJob);
  });

  it('should delete a job', (done) => {
    service.deleteJob(1).subscribe(() => {
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/jobs/1');
    expect(req.request.method).toBe('DELETE');
    req.flush(null);
  });

  it('should search jobs with filters', (done) => {
    const searchParams: SearchParams = {
      company: 'Google',
      status: 'Applied',
      page: 1,
      limit: 10
    };

    const mockResponse: PaginatedResponse<JobApplication> = {
      data: [
        {
          id: 1,
          company: 'Google',
          position: 'Software Engineer',
          status: 'Applied',
          applied_date: new Date().toISOString()
        }
      ],
      page: 1,
      limit: 10,
      total_count: 1,
      total_pages: 1
    };

    service.searchJobs(searchParams).subscribe(response => {
      expect(response).toEqual(mockResponse);
      done();
    });

    const req = httpMock.expectOne(req => 
      req.url === 'http://localhost:8080/api/jobs/search' &&
      req.params.get('company') === 'Google' &&
      req.params.get('status') === 'Applied'
    );
    expect(req.request.method).toBe('GET');
    req.flush(mockResponse);
  });

  it('should get jobs by status', (done) => {
    const mockJobs: JobApplication[] = [
      {
        id: 1,
        company: 'Google',
        position: 'Software Engineer',
        status: 'Applied',
        applied_date: new Date().toISOString()
      }
    ];

    service.getJobsByStatus('Applied').subscribe(jobs => {
      expect(jobs).toEqual(mockJobs);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/jobs/status/Applied');
    expect(req.request.method).toBe('GET');
    req.flush(mockJobs);
  });
});
