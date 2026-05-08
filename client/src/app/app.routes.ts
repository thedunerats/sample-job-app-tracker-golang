import { Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';
import { JobListComponent } from './components/job-list/job-list.component';
import { JobFormComponent } from './components/job-form/job-form.component';
import { JobDetailComponent } from './components/job-detail/job-detail.component';
import { authGuard } from './guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: '/jobs', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'jobs', component: JobListComponent, canActivate: [authGuard] },
  { path: 'jobs/new', component: JobFormComponent, canActivate: [authGuard] },
  { path: 'jobs/:id', component: JobDetailComponent, canActivate: [authGuard] },
  { path: 'jobs/:id/edit', component: JobFormComponent, canActivate: [authGuard] },
  { path: '**', redirectTo: '/jobs' }
];
