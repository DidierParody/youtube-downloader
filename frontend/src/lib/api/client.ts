import { env } from '$env/dynamic/public';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { auth } from '$lib/stores/auth';

const API_BASE_URL = env.PUBLIC_API_BASE_URL || 'http://localhost:3000';

async function getToken(): Promise<string | null> {
  if (browser) {
    return localStorage.getItem('auth_token');
  }
  return null;
}

interface RequestOptions extends RequestInit {
  requireAuth?: boolean;
}

export async function apiClient<T>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> {
  const { requireAuth = true, ...rest } = options;
  const url = `${API_BASE_URL}${endpoint}`;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(rest.headers as Record<string, string>)
  };

  if (requireAuth) {
    const token = await getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  try {
    const response = await fetch(url, {
      ...rest,
      headers
    });

    if (response.status === 401) {
      auth.clearAuth();
      if (browser) {
        goto('/auth/login');
      }
      throw new Error('Unauthorized');
    }

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || `HTTP error! status: ${response.status}`);
    }

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return response.json() as Promise<T>;
    }
    return null as T;
  } catch (error) {
    console.error('API Client Error:', error);
    throw error;
  }
}

export function getApiBaseUrl(): string {
  return API_BASE_URL;
}