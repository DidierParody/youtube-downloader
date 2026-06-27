import { apiClient } from './client';
import type { Download, DownloadProgress } from '$lib/types';

export interface CreateDownloadData {
  url: string;
  format?: string;
  quality?: string;
}

export async function requestDownload(data: CreateDownloadData): Promise<Download> {
  return apiClient<Download>('/api/v1/downloads', {
    method: 'POST',
    body: JSON.stringify(data)
  });
}

export async function getDownloads(): Promise<Download[]> {
  return apiClient<Download[]>('/api/v1/downloads');
}

export async function getDownload(id: string): Promise<Download> {
  return apiClient<Download>(`/api/v1/downloads/${id}`);
}

export function createProgressConnection(
  id: string,
  onProgress: (progress: DownloadProgress) => void,
  onError?: (error: Event) => void
): EventSource {
  const token = localStorage.getItem('auth_token');
  const url = new URL(`/api/v1/downloads/${id}/progress`, 'http://localhost:3000');
  if (token) {
    url.searchParams.set('token', token);
  }
  
  const eventSource = new EventSource(url.toString());
  
  eventSource.onmessage = (event) => {
    try {
      const progress = JSON.parse(event.data) as DownloadProgress;
      onProgress(progress);
    } catch (err) {
      console.error('Failed to parse progress data:', err);
    }
  };
  
  eventSource.onerror = (error) => {
    console.error('SSE Error:', error);
    if (onError) onError(error);
  };
  
  return eventSource;
}