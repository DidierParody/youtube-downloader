<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth';
  import type { Download, DownloadProgress } from '$lib/types';
  import { getDownload } from '$lib/api/downloads';
  import { createProgressConnection, type ProgressConnection } from '$lib/api/progress';
  import ProgressBar from '$lib/components/ProgressBar.svelte';

  const downloadId: string = $page.params.id;
  let download: Download | null = $state(null);
  let progress: DownloadProgress | null = $state(null);
  let loading = $state(true);
  let error: string | null = $state(null);
  let connection: ProgressConnection | null = $state(null);

  onMount(async () => {
    if (!$isAuthenticated) {
      goto('/auth/login');
      return;
    }

    try {
      download = await getDownload(downloadId);
    } catch (err) {
      error = 'Failed to load download details';
      console.error(err);
    } finally {
      loading = false;
    }

    // Setup SSE for real-time progress
    if (download?.status === 'downloading') {
      connection = createProgressConnection(
        downloadId,
        (p: DownloadProgress) => {
          progress = p;
        },
        (err: Event) => {
          console.error('Progress connection error:', err);
        }
      );
    }
  });

  onDestroy(() => {
    if (connection) {
      connection.close();
    }
  });

  function getStatusColor(status: string): string {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'downloading':
        return 'bg-blue-100 text-blue-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }
</script>

<svelte:head>
  <title>Download Details - YouTube Downloader</title>
</svelte:head>

<div class="mx-auto max-w-3xl">
  <a
    href="/downloads"
    class="mb-6 inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-700"
  >
    <svg class="mr-1 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
    </svg>
    Back to Downloads
  </a>

  {#if loading}
    <div class="flex h-48 items-center justify-center">
      <svg class="h-8 w-8 animate-spin text-primary-600" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
      </svg>
    </div>
  {:else if error}
    <div class="rounded-md bg-red-50 p-4 text-sm text-red-700">
      {error}
    </div>
  {:else if download}
    <div class="space-y-6">
      <!-- Video Info -->
      <div class="card p-6">
        <div class="flex items-start gap-4">
          {#if download.video?.thumbnailUrl}
            <img
              src={download.video.thumbnailUrl}
              alt=""
              class="h-32 w-48 rounded-lg object-cover"
            />
          {/if}
          <div class="flex-1">
            <h1 class="text-xl font-bold text-gray-900">
              {download.video?.title || 'Unknown Video'}
            </h1>
            <p class="mt-1 text-sm text-gray-500">{download.video?.uploader || 'Unknown uploader'}</p>
            <div class="mt-2 inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium {getStatusColor(download.status)}">
              {download.status}
            </div>
          </div>
        </div>
      </div>

      <!-- Progress -->
      {#if download.status === 'downloading' && progress}
        <div class="card p-6">
          <h2 class="mb-4 text-lg font-semibold text-gray-900">Download Progress</h2>
          <ProgressBar progress={progress.progress} status={progress.status} size="lg" />
          <div class="mt-4 grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-gray-500">Speed:</span>
              <span class="ml-1 font-mono">{progress.speed}</span>
            </div>
            <div>
              <span class="text-gray-500">ETA:</span>
              <span class="ml-1 font-mono">{progress.eta}</span>
            </div>
            <div>
              <span class="text-gray-500">Downloaded:</span>
              <span class="ml-1 font-mono">{progress.downloaded} MB</span>
            </div>
            <div>
              <span class="text-gray-500">Total:</span>
              <span class="ml-1 font-mono">{progress.total} MB</span>
            </div>
          </div>
        </div>
      {/if}

      <!-- Download details -->
      <div class="card p-6">
        <h2 class="mb-4 text-lg font-semibold text-gray-900">Details</h2>
        <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <dt class="text-sm font-medium text-gray-500">Format</dt>
            <dd class="mt-1 text-sm text-gray-900">{download.format}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500">Quality</dt>
            <dd class="mt-1 text-sm text-gray-900">{download.quality}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500">Started</dt>
            <dd class="mt-1 text-sm text-gray-900">{new Date(download.createdAt).toLocaleString()}</dd>
          </div>
          {#if download.completedAt}
            <div>
              <dt class="text-sm font-medium text-gray-500">Completed</dt>
              <dd class="mt-1 text-sm text-gray-900">{new Date(download.completedAt).toLocaleString()}</dd>
            </div>
          {/if}
          {#if download.error}
            <div class="sm:col-span-2">
              <dt class="text-sm font-medium text-red-500">Error</dt>
              <dd class="mt-1 text-sm text-red-700">{download.error}</dd>
            </div>
          {/if}
        </dl>
      </div>
    </div>
  {/if}
</div>