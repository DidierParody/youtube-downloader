<script lang="ts">
  import DownloadForm from '$lib/components/DownloadForm.svelte';
  import { onMount } from 'svelte';
  import { getVideos } from '$lib/api/videos';
  import VideoCard from '$lib/components/VideoCard.svelte';
  import type { Video } from '$lib/types';

  let videos: Video[] = $state([]);
  let loading = $state(true);
  let error: string | null = $state(null);

  onMount(async () => {
    try {
      videos = await getVideos();
    } catch (err) {
      error = 'Failed to load video previews';
    } finally {
      loading = false;
    }
  });

  function handleVideoClick(video: Video) {
    window.location.href = `/videos/${video.id}`;
  }
</script>

<svelte:head>
  <title>YouTube Downloader</title>
</svelte:head>

<div class="mx-auto max-w-3xl text-center">
  <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
    Download YouTube Videos
  </h1>
  <p class="mt-4 text-lg text-gray-600">
    Fast, simple, and free. Just paste a YouTube URL to get started.
  </p>
</div>

<div class="mt-12">
  <DownloadForm />
</div>

{#if videos.length > 0}
  <div class="mt-12">
    <h2 class="mb-6 text-xl font-semibold text-gray-900">Recent Downloads</h2>
    <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {#each videos as video (video.id)}
        <VideoCard {video} onClick={handleVideoClick} />
      {/each}
    </div>
  </div>
{/if}