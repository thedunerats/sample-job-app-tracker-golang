import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { JobService } from '../../services/job.service';
import { JobApplication } from '../../models/job.model';

@Component({
  selector: 'app-job-detail',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './job-detail.component.html',
  styleUrl: './job-detail.component.scss'
})
export class JobDetailComponent implements OnInit {
  job: JobApplication | null = null;
  loading = false;
  errorMessage = '';

  constructor(
    private jobService: JobService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      const id = +params['id'];
      this.loadJob(id);
    });
  }

  loadJob(id: number): void {
    this.loading = true;
    this.jobService.getJob(id).subscribe({
      next: (job) => {
        this.job = job;
        this.loading = false;
      },
      error: (error) => {
        this.errorMessage = 'Failed to load job details.';
        this.loading = false;
      }
    });
  }

  deleteJob(): void {
    if (!this.job || !confirm('Are you sure you want to delete this job application?')) {
      return;
    }

    this.jobService.deleteJob(this.job.id!).subscribe({
      next: () => {
        this.router.navigate(['/jobs']);
      },
      error: (error) => {
        alert('Failed to delete job application.');
      }
    });
  }
}
