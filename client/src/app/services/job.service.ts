import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import { JobApplication, PaginatedResponse, SearchParams } from '../models/job.model';

@Injectable({
  providedIn: 'root'
})
export class JobService {
  constructor(private api: ApiService) { }

  getAllJobs(page: number = 1, limit: number = 10): Observable<PaginatedResponse<JobApplication>> {
    return this.api.get<PaginatedResponse<JobApplication>>('/jobs', { page, limit });
  }

  getJob(id: number): Observable<JobApplication> {
    return this.api.get<JobApplication>(`/jobs/${id}`);
  }

  createJob(job: JobApplication): Observable<JobApplication> {
    return this.api.post<JobApplication>('/jobs', job);
  }

  updateJob(id: number, job: JobApplication): Observable<JobApplication> {
    return this.api.put<JobApplication>(`/jobs/${id}`, job);
  }

  deleteJob(id: number): Observable<void> {
    return this.api.delete<void>(`/jobs/${id}`);
  }

  getJobsByStatus(status: string): Observable<JobApplication[]> {
    return this.api.get<JobApplication[]>(`/jobs/status/${status}`);
  }

  searchJobs(params: SearchParams): Observable<PaginatedResponse<JobApplication>> {
    return this.api.get<PaginatedResponse<JobApplication>>('/jobs/search', params);
  }
}
