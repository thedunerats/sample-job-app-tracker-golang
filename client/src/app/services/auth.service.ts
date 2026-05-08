import { Injectable } from '@angular/core';
import { Observable, BehaviorSubject, tap } from 'rxjs';
import { ApiService } from './api.service';
import { LoginRequest, LoginResponse, RegisterRequest, User } from '../models/job.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  constructor(private api: ApiService) {
    // Check if user is already logged in
    const token = this.getToken();
    if (token) {
      this.getCurrentUser().subscribe();
    }
  }

  register(data: RegisterRequest): Observable<LoginResponse> {
    return this.api.post<LoginResponse>('/auth/register', data).pipe(
      tap(response => this.setSession(response))
    );
  }

  login(data: LoginRequest): Observable<LoginResponse> {
    return this.api.post<LoginResponse>('/auth/login', data).pipe(
      tap(response => this.setSession(response))
    );
  }

  logout(): void {
    localStorage.removeItem('token');
    this.currentUserSubject.next(null);
  }

  getCurrentUser(): Observable<User> {
    return this.api.get<User>('/auth/me').pipe(
      tap(user => this.currentUserSubject.next(user))
    );
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  getToken(): string | null {
    return localStorage.getItem('token');
  }

  private setSession(response: LoginResponse): void {
    localStorage.setItem('token', response.token);
    this.currentUserSubject.next(response.user);
  }
}
