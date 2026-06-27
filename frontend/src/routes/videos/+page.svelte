<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth';
  import { getVideos } from '$lib/api/videos';
  import type { Video } from '$lib/types';
  import VideoCard from '$lib/components/VideoCard.svelte';

  let videos: Video[] = $state([]);
  let searchQuery = $state('');
  let loading = $state(true);
  let error: string | null = $state(null);

  let filteredVideos = $derived(
    searchQuery.trim()
      ? videos.filter((v: Video) =>
          v.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          v.uploader.toLowerCase().includes(searchQuery.toLowerCase())
        )
      : videos
  );

  onMount(async () => {
    if (!$isAuthenticated) {
      goto('/auth/login');
      return;
    }

    try {
      videos = await getVideos();
    } catch (err) {
      error = 'Failed to load videos';
      console.error(err);
    } finally {
      loading = false;
    }
  });

  function handleVideoClick(video: Video) {
    goto(`/videos/${video.id}`);
  }
</script>

<svelte:head>
  <title>Videos - YouTube Downloader</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
    <h1 class="text-2xl font-bold text-gray-900">Videos</h1>
    <div class="relative max-w-md flex-1">
      <input
        type="search"
        placeholder="Search videos..."
        bind:value={searchQuery}
        class="input pl-10"
      />
      <svg
        class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
        />
      </svg>
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <svg class="h-8 w-8 animate-spin text-primary-600" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
      </svg>
    </div>
  {:else if error}
    <div class="rounded-md bg-red-50 p-4 text-sm text-red-700">
      {error}
    </div>
  {:else if filteredVideos.length === 0}
    <div class="py-12 text-center">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
      </svg>
      <p class="mt-2 text-sm text-gray-500">
        {#if searchQuery}
          No videos found for "{searchQuery}"
        {:else}
          No videos available yet.
        {/if}
      </p>
    </div>
  {:else}
    <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {#each filteredVideos as video (video.id)}
        <VideoCard {video} onClick={handleVideoClick} />
      {/each}
    </div>
  {/if}
</div>