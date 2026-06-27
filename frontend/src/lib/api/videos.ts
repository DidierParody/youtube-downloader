import { apiClient } from './client';
import type { Video } from '$lib/types';

export async function getVideo(id: string): Promise<Video> {
  return apiClient<Video>(`/api/v1/videos/${id}`);
}

export async function getVideos(): Promise<Video[]> {
  return apiClient<Video[]>('/api/v1/videos');
}

export async function searchVideos(query: string): Promise<Video[]> {
  const params = new URLSearchParams({ q: query });
  return apiClient<Video[]>(`/api/v1/videos?${params.toString()}`);
}