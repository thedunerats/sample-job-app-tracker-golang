import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { JobService } from '../../services/job.service';
import { AuthService } from '../../services/auth.service';
import { JobApplication, SearchParams, User } from '../../models/job.model';

@Component({
  selector: 'app-job-list',
  standalone: true,
  imports: [CommonModule, RouterModule, ReactiveFormsModule],
  templateUrl: './job-list.component.html',
  styleUrl: './job-list.component.scss'
})
export class JobListComponent implements OnInit {
  jobs: JobApplication[] = [];
  currentUser: User | null = null;
  searchForm!: FormGroup;
  
  // Pagination
  currentPage = 1;
  pageSize = 10;
  totalCount = 0;
  totalPages = 0;
  
  loading = false;
  errorMessage = '';

  constructor(
    private jobService: JobService,
    private authService: AuthService,
    private router: Router,
    private fb: FormBuilder
  ) {}

  ngOnInit(): void {
    this.searchForm = this.fb.group({
      company: [''],
      position: [''],
      status: [''],
      location: ['']
    });

    this.authService.currentUser$.subscribe(user => {
      this.currentUser = user;
    });

    this.loadJobs();
  }

  loadJobs(): void {
    this.loading = true;
    this.errorMessage = '';
    
    const searchParams = this.getSearchParams();
    const hasFilters = Object.keys(searchParams).some(key => 
      key !== 'page' && key !== 'limit' && searchParams[key as keyof SearchParams]
    );

    const request = hasFilters 
      ? this.jobService.searchJobs(searchParams)
      : this.jobService.getAllJobs(this.currentPage, this.pageSize);

    request.subscribe({
      next: (response) => {
        this.jobs = response.data;
        this.currentPage = response.page;
        this.pageSize = response.limit;
        this.totalCount = response.total_count;
        this.totalPages = response.total_pages;
        this.loading = false;
      },
      error: (error) => {
        this.errorMessage = 'Failed to load jobs. Please try again.';
        this.loading = false;
      }
    });
  }

  getSearchParams(): SearchParams {
    const formValue = this.searchForm.value;
    return {
      company: formValue.company || undefined,
      position: formValue.position || undefined,
      status: formValue.status || undefined,
      location: formValue.location || undefined,
      page: this.currentPage,
      limit: this.pageSize
    };
  }

  onSearch(): void {
    this.currentPage = 1;
    this.loadJobs();
  }

  onClearSearch(): void {
    this.searchForm.reset();
    this.currentPage = 1;
    this.loadJobs();
  }

  onPageChange(page: number): void {
    if (page < 1 || page > this.totalPages) return;
    this.currentPage = page;
    this.loadJobs();
  }

  deleteJob(id: number): void {
    if (!confirm('Are you sure you want to delete this job application?')) {
      return;
    }

    this.jobService.deleteJob(id).subscribe({
      next: () => {
        this.loadJobs();
      },
      error: (error) => {
        alert('Failed to delete job application.');
      }
    });
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  get pages(): number[] {
    return Array.from({ length: this.totalPages }, (_, i) => i + 1);
  }
}
