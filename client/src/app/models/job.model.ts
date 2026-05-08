export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface JobApplication {
  id?: number;
  user_id?: number;
  company: string;
  position: string;
  status: string;
  applied_date: string;
  notes?: string;
  contact_info?: string;
  salary?: string;
  location?: string;
  created_at?: string;
  updated_at?: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total_count: number;
  total_pages: number;
}

export interface SearchParams {
  company?: string;
  position?: string;
  status?: string;
  location?: string;
  page?: number;
  limit?: number;
}
