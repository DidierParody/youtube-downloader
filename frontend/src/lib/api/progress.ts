import { getApiBaseUrl } from './client';
import type { DownloadProgress } from '$lib/types';

export interface ProgressConnection {
  close: () => void;
  reconnect: () => void;
}

export function createProgressConnection(
  downloadId: string,
  onProgress: (progress: DownloadProgress) => void,
  onError?: (error: Event) => void
): ProgressConnection {
  let eventSource: EventSource | null = null;
  let isClosed = false;

  const baseUrl = getApiBaseUrl();
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  const url = new URL(`/api/v1/downloads/${downloadId}/progress`, baseUrl);
  if (token) {
    url.searchParams.set('token', token);
  }

  function connect() {
    if (isClosed) return;
    
    eventSource = new EventSource(url.toString());

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
      // Auto-reconnect on error after 3 seconds
      setTimeout(() => {
        if (!isClosed) connect();
      }, 3000);
    };
  }

  connect();

  return {
    close: () => {
      isClosed = true;
      if (eventSource) {
        eventSource.close();
        eventSource = null;
      }
    },
    reconnect: () => {
      if (eventSource) {
        eventSource.close();
      }
      connect();
    }
  };
}