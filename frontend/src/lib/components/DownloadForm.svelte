<script lang="ts">
  import { goto } from '$app/navigation';
  import { requestDownload } from '$lib/api/downloads';
  import { downloadStore } from '$lib/stores/downloads';

  let url = $state('');
  let isSubmitting = $state(false);
  let error: string | null = $state(null);
  let success = $state(false);

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    
    if (!url.trim()) {
      error = 'Please enter a YouTube URL';
      return;
    }

    // Basic URL validation
    const youtubeRegex = /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+$/;
    if (!youtubeRegex.test(url.trim())) {
      error = 'Please enter a valid YouTube URL';
      return;
    }

    isSubmitting = true;
    error = null;
    success = false;

    try {
      const download = await requestDownload({
        url: url.trim(),
        format: 'mp4',
        quality: '720p'
      });
      
      downloadStore.addDownload(download);
      success = true;
      url = '';
      
      // Redirect to the download detail page
      setTimeout(() => {
        goto(`/downloads/${download.id}`);
      }, 1000);
    } catch (err: any) {
      error = err.message || 'Failed to start download. Please try again.';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="card p-6">
  <h2 class="mb-4 text-lg font-semibold text-gray-900">Download Video</h2>
  
  <form onsubmit={handleSubmit} class="space-y-4">
    <div>
      <label for="youtube-url" class="mb-1 block text-sm font-medium text-gray-700">
        YouTube URL
      </label>
      <input
        id="youtube-url"
        type="url"
        placeholder="https://youtube.com/watch?v=..."
        bind:value={url}
        class="input"
        required
        disabled={isSubmitting}
      />
    </div>

    {#if error}
      <div class="rounded-md bg-red-50 p-3 text-sm text-red-700">
        {error}
      </div>
    {/if}

    {#if success}
      <div class="rounded-md bg-green-50 p-3 text-sm text-green-700">
        Download started successfully! Redirecting...
      </div>
    {/if}

    <button
      type="submit"
      disabled={isSubmitting}
      class="btn btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {#if isSubmitting}
        <svg class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Starting Download...
      {:else}
        Download
      {/if}
    </button>
  </form>
</div>