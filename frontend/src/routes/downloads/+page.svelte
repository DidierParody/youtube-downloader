<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { isAuthenticated } from '$lib/stores/auth';
  import { downloadStore } from '$lib/stores/downloads';
  import { getDownloads } from '$lib/api/downloads';
  import type { Download } from '$lib/types';
  import DownloadList from '$lib/components/DownloadList.svelte';

  let downloads: Download[] = $state([]);
  let loading = $state(true);
  let error: string | null = $state(null);

  async function loadDownloads() {
    loading = true;
    error = null;
    try {
      downloads = await getDownloads();
      downloadStore.setDownloads(downloads);
    } catch (err) {
      error = 'Failed to load downloads';
      console.error(err);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    if (!$isAuthenticated) {
      goto('/auth/login');
      return;
    }
    loadDownloads();
  });
</script>

<svelte:head>
  <title>Downloads - YouTube Downloader</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-gray-900">Downloads</h1>
    <a
      href="/"
      class="btn btn-primary text-sm"
    >
      New Download
    </a>
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
  {:else}
    <DownloadList {downloads} onRefresh={loadDownloads} />
  {/if}
</div>