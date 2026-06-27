<script lang="ts">
  import { goto } from '$app/navigation';
  import type { Download } from '$lib/types';

  interface Props {
    downloads: Download[];
    onRefresh?: () => void;
  }

  let { downloads, onRefresh }: Props = $props();

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
      case 'cancelled':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }

  function navigateToDownload(id: string) {
    goto(`/downloads/${id}`);
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleString();
  }
</script>

<div class="overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-gray-900/5">
  <div class="flex items-center justify-between border-b border-gray-200 px-6 py-4">
    <h3 class="text-base font-semibold text-gray-900">Download History</h3>
    {#if onRefresh}
      <button
        onclick={onRefresh}
        class="inline-flex items-center text-sm text-primary-600 hover:text-primary-700"
      >
        <svg class="mr-1 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M9.334 20.538A8.953 8.953 0 0112 21c4.968 0 9-3.582 9-8s-4.032-8 9-8 9 4.032 9 9-4.032 9-9 9m0 0l-3.469-3.469M12 21l3.469-3.469" />
        </svg>
        Refresh
      </button>
    {/if}
  </div>
  
  {#if downloads.length === 0}
    <div class="px-6 py-12 text-center">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M3 16.5l6.75-6.75M3 16.5l6.75 6.75m12-15L15 6m0 0l6.75-6.75M15 6v6.75" />
      </svg>
      <p class="mt-2 text-sm text-gray-500">No downloads yet. Start by pasting a YouTube URL above.</p>
    </div>
  {:else}
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Video</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Status</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Progress</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Date</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white">
          {#each downloads as download (download.id)}
            <tr
              class="cursor-pointer transition-colors hover:bg-gray-50"
              onclick={() => navigateToDownload(download.id)}
            >
              <td class="px-6 py-4">
                <div class="flex items-center">
                  <div class="h-10 w-10 flex-shrink-0 overflow-hidden rounded bg-gray-100">
                    {#if download.video?.thumbnailUrl}
                      <img
                        src={download.video.thumbnailUrl}
                        alt=""
                        class="h-full w-full object-cover"
                      />
                    {:else}
                      <div class="flex h-full w-full items-center justify-center">
                        <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3" />
                        </svg>
                      </div>
                    {/if}
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">
                      {#if download.video}
                        {download.video.title}
                      {:else}
                        {download.url}
                      {/if}
                    </div>
                    <div class="text-sm text-gray-500">{download.format} • {download.quality}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${getStatusColor(download.status)}`}>
                  {download.status}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="w-32">
                  <div class="h-2 w-full rounded-full bg-gray-200">
                    <div
                      class="h-2 rounded-full bg-primary-600 transition-all"
                      style="width: {download.progress}%"
                    ></div>
                  </div>
                  <span class="mt-1 text-xs text-gray-500">{download.progress.toFixed(1)}%</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {formatDate(download.createdAt)}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>