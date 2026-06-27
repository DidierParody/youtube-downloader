<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth';
  import type { Video } from '$lib/types';
  import { getVideo } from '$lib/api/videos';
  import { onMount } from 'svelte';

  const videoId: string = $page.params.id;
  let video: Video | null = $state(null);
  let loading = $state(true);
  let error: string | null = $state(null);

  onMount(async () => {
    if (!$isAuthenticated) {
      goto('/auth/login');
      return;
    }

    try {
      video = await getVideo(videoId);
    } catch (err) {
      error = 'Failed to load video details';
    } finally {
      loading = false;
    }
  });

  function formatDuration(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }
</script>

<svelte:head>
  <title>{video?.title || 'Video'} - YouTube Downloader</title>
</svelte:head>

<div class="mx-auto max-w-4xl">
  <a
    href="/videos"
    class="mb-6 inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-700"
  >
    <svg class="mr-1 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
    </svg>
    Back to Videos
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
  {:else if video}
    <div class="space-y-6">
      <!-- Video Player / Thumbnail -->
      <div class="card overflow-hidden">
        <div class="relative aspect-video bg-gray-900">
          {#if video.thumbnailUrl}
            <img
              src={video.thumbnailUrl}
              alt={video.title}
              class="h-full w-full object-cover"
            />
          {:else}
            <div class="flex h-full w-full items-center justify-center">
              <svg class="h-16 w-16 text-gray-600" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-18-9m0 0l18 9m-18-9v18m0 0l18-9m-18 9l18 9" />
              </svg>
            </div>
          {/if}
          <div class="absolute bottom-2 right-2 rounded bg-black/70 px-2 py-1 text-xs text-white">
            {formatDuration(video.duration)}
          </div>
        </div>
      </div>

      <!-- Video Info -->
      <div class="card p-6">
        <h1 class="text-2xl font-bold text-gray-900">{video.title}</h1>
        <p class="mt-2 text-sm text-gray-500">{video.uploader}</p>
        
        <div class="mt-4 flex items-center gap-4 text-sm text-gray-500">
          <span>{new Date(video.createdAt).toLocaleDateString()}</span>
          <span>•</span>
          <span>{formatDuration(video.duration)}</span>
        </div>

        {#if video.description}
          <div class="mt-6">
            <h2 class="text-sm font-medium text-gray-900">Description</h2>
            <p class="mt-2 text-sm text-gray-600 whitespace-pre-wrap">{video.description}</p>
          </div>
        {/if}
      </div>

      <!-- Actions -->
      <div class="card p-6">
        <h2 class="text-sm font-medium text-gray-900 mb-4">Actions</h2>
        <div class="flex gap-4">
          <a
            href="/"
            class="btn btn-primary"
          >
            Download
          </a>
        </div>
      </div>
    </div>
  {/if}
</div>