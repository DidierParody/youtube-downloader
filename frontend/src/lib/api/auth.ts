import { apiClient, getApiBaseUrl } from './client';
import type { User } from '$lib/types';

interface RegisterData {
  username: string;
  email: string;
  password: string;
}

interface LoginData {
  email: string;
  password: string;
}

interface AuthResponse {
  token: string;
  user: User;
}

export async function register(data: RegisterData): Promise<AuthResponse> {
  return apiClient<AuthResponse>('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify(data),
    requireAuth: false
  });
}

export async function login(data: LoginData): Promise<AuthResponse> {
  return apiClient<AuthResponse>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify(data),
    requireAuth: false
  });
}

export async function me(): Promise<User> {
  return apiClient<User>('/api/v1/auth/me');
}

export function getSSEUrl(endpoint: string): string {
  const baseUrl = getApiBaseUrl();
  return `${baseUrl}${endpoint}`;
}