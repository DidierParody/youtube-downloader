import { writable, derived } from 'svelte/store';
import type { Download, DownloadProgress } from '$lib/types';

interface DownloadMap {
  [id: string]: Download;
}

interface ProgressMap {
  [id: string]: DownloadProgress;
}

function createDownloadStore() {
  const { subscribe, update } = writable<DownloadMap>({});
  const progressStore = writable<ProgressMap>({});

  return {
    subscribe,
    progress: { subscribe: progressStore.subscribe },
    addDownload: (download: Download) => {
      update(state => ({ ...state, [download.id]: download }));
    },
    updateDownload: (id: string, data: Partial<Download>) => {
      update(state => ({
        ...state,
        [id]: { ...state[id], ...data }
      }));
    },
    setDownloads: (downloads: Download[]) => {
      const map: DownloadMap = {};
      for (const d of downloads) {
        map[d.id] = d;
      }
      update(() => map);
    },
    removeDownload: (id: string) => {
      update(state => {
        const { [id]: _, ...rest } = state;
        return rest;
      });
    },
    updateProgress: (id: string, progress: DownloadProgress) => {
      progressStore.update(state => ({ ...state, [id]: progress }));
    },
    clearProgress: (id: string) => {
      progressStore.update(state => {
        const { [id]: _, ...rest } = state;
        return rest;
      });
    }
  };
}

export const downloadStore = createDownloadStore();
export const activeDownloads = derived(downloadStore, $store =>
  Object.values($store).filter(d => d.status === 'downloading' || d.status === 'pending')
);
export const completedDownloads = derived(downloadStore, $store =>
  Object.values($store).filter(d => d.status === 'completed')
);