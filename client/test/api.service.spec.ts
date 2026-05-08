import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ApiService } from '../src/app/services/api.service';

describe('ApiService', () => {
  let service: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [ApiService]
    });
    service = TestBed.inject(ApiService);
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

  it('should make GET request', (done) => {
    const mockData = { test: 'data' };

    service.get('/test').subscribe(data => {
      expect(data).toEqual(mockData);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/test');
    expect(req.request.method).toBe('GET');
    req.flush(mockData);
  });

  it('should make POST request', (done) => {
    const postData = { name: 'test' };
    const mockResponse = { id: 1, name: 'test' };

    service.post('/test', postData).subscribe(data => {
      expect(data).toEqual(mockResponse);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/test');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(postData);
    req.flush(mockResponse);
  });

  it('should make PUT request', (done) => {
    const putData = { name: 'updated' };
    const mockResponse = { id: 1, name: 'updated' };

    service.put('/test', putData).subscribe(data => {
      expect(data).toEqual(mockResponse);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/test');
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(putData);
    req.flush(mockResponse);
  });

  it('should make DELETE request', (done) => {
    service.delete('/test').subscribe(() => {
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/test');
    expect(req.request.method).toBe('DELETE');
    req.flush(null);
  });

  it('should include Authorization header when token exists', (done) => {
    localStorage.setItem('token', 'test-token');
    const mockData = { test: 'data' };

    service.get('/test').subscribe(data => {
      expect(data).toEqual(mockData);
      done();
    });

    const req = httpMock.expectOne('http://localhost:8080/api/test');
    expect(req.request.headers.get('Authorization')).toBe('Bearer test-token');
    expect(req.request.headers.get('Content-Type')).toBe('application/json');
    req.flush(mockData);
  });

  it('should include query parameters', (done) => {
    const params = { page: 1, limit: 10, search: 'test' };
    const mockData = { results: [] };

    service.get('/test', params).subscribe(data => {
      expect(data).toEqual(mockData);
      done();
    });

    const req = httpMock.expectOne(req => 
      req.url === 'http://localhost:8080/api/test' &&
      req.params.get('page') === '1' &&
      req.params.get('limit') === '10' &&
      req.params.get('search') === 'test'
    );
    expect(req.request.method).toBe('GET');
    req.flush(mockData);
  });

  it('should skip empty or null parameters', (done) => {
    const params = { page: 1, search: '', status: null, filter: undefined };
    const mockData = { results: [] };

    service.get('/test', params).subscribe(data => {
      expect(data).toEqual(mockData);
      done();
    });

    const req = httpMock.expectOne(req => 
      req.url === 'http://localhost:8080/api/test' &&
      req.params.get('page') === '1' &&
      req.params.has('search') === false &&
      req.params.has('status') === false &&
      req.params.has('filter') === false
    );
    expect(req.request.method).toBe('GET');
    req.flush(mockData);
  });
});
