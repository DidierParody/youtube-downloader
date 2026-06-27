export interface User {
  id: number;
  username: string;
  email: string;
}

export interface Video {
  id: string;
  youtubeId: string;
  title: string;
  description: string;
  thumbnailUrl: string;
  duration: number;
  uploader: string;
  createdAt: string;
  status: string;
}

export interface Download {
  id: string;
  videoId: string;
  url: string;
  status: 'pending' | 'downloading' | 'completed' | 'failed' | 'cancelled';
  progress: number;
  format: string;
  quality: string;
  filePath?: string;
  error?: string;
  createdAt: string;
  completedAt?: string;
  video?: Video;
}

export interface DownloadProgress {
  id: string;
  progress: number;
  status: string;
  speed: string;
  downloaded: number;
  total: number;
  eta: string;
}

export interface ApiError {
  message: string;
  status: number;
}