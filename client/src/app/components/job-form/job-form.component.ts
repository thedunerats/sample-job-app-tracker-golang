import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { JobService } from '../../services/job.service';
import { JobApplication } from '../../models/job.model';

@Component({
  selector: 'app-job-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  templateUrl: './job-form.component.html',
  styleUrl: './job-form.component.scss'
})
export class JobFormComponent implements OnInit {
  jobForm!: FormGroup;
  isEditMode = false;
  jobId: number | null = null;
  loading = false;
  errorMessage = '';

  constructor(
    private fb: FormBuilder,
    private jobService: JobService,
    private router: Router,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.jobForm = this.fb.group({
      company: ['', Validators.required],
      position: ['', Validators.required],
      status: ['Applied', Validators.required],
      applied_date: ['', Validators.required],
      location: [''],
      salary: [''],
      contact_info: [''],
      notes: ['']
    });

    this.route.params.subscribe(params => {
      if (params['id']) {
        this.jobId = +params['id'];
        this.isEditMode = true;
        this.loadJob(this.jobId);
      }
    });
  }

  loadJob(id: number): void {
    this.loading = true;
    this.jobService.getJob(id).subscribe({
      next: (job) => {
        this.jobForm.patchValue({
          company: job.company,
          position: job.position,
          status: job.status,
          applied_date: job.applied_date.split('T')[0],
          location: job.location,
          salary: job.salary,
          contact_info: job.contact_info,
          notes: job.notes
        });
        this.loading = false;
      },
      error: (error) => {
        this.errorMessage = 'Failed to load job details.';
        this.loading = false;
      }
    });
  }

  onSubmit(): void {
    if (this.jobForm.invalid) {
      return;
    }

    this.loading = true;
    this.errorMessage = '';

    const jobData: JobApplication = {
      ...this.jobForm.value,
      applied_date: new Date(this.jobForm.value.applied_date).toISOString()
    };

    const request = this.isEditMode && this.jobId
      ? this.jobService.updateJob(this.jobId, jobData)
      : this.jobService.createJob(jobData);

    request.subscribe({
      next: () => {
        this.router.navigate(['/jobs']);
      },
      error: (error) => {
        this.errorMessage = this.isEditMode 
          ? 'Failed to update job application.'
          : 'Failed to create job application.';
        this.loading = false;
      }
    });
  }
}
