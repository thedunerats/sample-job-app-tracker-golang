import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { AuthService } from '../src/app/services/auth.service';
import { ApiService } from '../src/app/services/api.service';
import { LoginRequest, RegisterRequest, LoginResponse, User } from '../src/app/models/job.model';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;
  let apiService: ApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [AuthService, ApiService]
    });
    service = TestBed.inject(AuthService);
    apiService = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
    localStorage.clear();
  });

  afterEach(() => {
    httpMock.verify();
    localStorage.clear();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should register a new user', (done) => {
    const registerData: RegisterRequest = {
      email: 'test@example.com',
      password: 'password123',
      name: 'Test User'
    };

    const mockResponse: LoginResponse = {
      token: 'test-token',
      user: {
        id: 1,
        email: 'test@example.com',
        name: 'Test User',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
    };

    service.register(registerData).subscribe(response => {
      expect(response).toEqual(mockResponse);
      expect(localStorage.getItem('token')).toBe('test-token');
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/auth/register');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(registerData);
    req.flush(mockResponse);
  });

  it('should login a user', (done) => {
    const loginData: LoginRequest = {
      email: 'test@example.com',
      password: 'password123'
    };

    const mockResponse: LoginResponse = {
      token: 'test-token',
      user: {
        id: 1,
        email: 'test@example.com',
        name: 'Test User',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
    };

    service.login(loginData).subscribe(response => {
      expect(response).toEqual(mockResponse);
      expect(localStorage.getItem('token')).toBe('test-token');
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/auth/login');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(loginData);
    req.flush(mockResponse);
  });

  it('should get current user', (done) => {
    const mockUser: User = {
      id: 1,
      email: 'test@example.com',
      name: 'Test User',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };

    localStorage.setItem('token', 'test-token');

    service.getCurrentUser().subscribe(user => {
      expect(user).toEqual(mockUser);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/auth/me');
    expect(req.request.method).toBe('GET');
    expect(req.request.headers.get('Authorization')).toBe('Bearer test-token');
    req.flush(mockUser);
  });

  it('should logout user', () => {
    localStorage.setItem('token', 'test-token');
    service.logout();
    expect(localStorage.getItem('token')).toBeNull();
  });

  it('should check if user is authenticated', () => {
    expect(service.isAuthenticated()).toBe(false);
    
    localStorage.setItem('token', 'test-token');
    expect(service.isAuthenticated()).toBe(true);
  });

  it('should get token', () => {
    expect(service.getToken()).toBeNull();
    
    localStorage.setItem('token', 'test-token');
    expect(service.getToken()).toBe('test-token');
  });
});
