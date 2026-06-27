<script lang="ts">
  import type { Video } from '$lib/types';

  interface Props {
    video: Video;
    onClick?: (video: Video) => void;
  }

  let { video, onClick }: Props = $props();

  function formatDuration(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  function handleClick() {
    if (onClick) onClick(video);
  }
</script>

<div
  class="group cursor-pointer overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-gray-900/5 transition-all hover:shadow-md hover:ring-primary-300"
  onclick={handleClick}
  role="button"
  tabindex="0"
  onkeydown={(e: KeyboardEvent) => e.key === 'Enter' && handleClick()}
>
  <div class="relative aspect-video overflow-hidden bg-gray-100">
    <img
      src={video.thumbnailUrl || 'https://via.placeholder.com/640x360'}
      alt={video.title}
      class="h-full w-full object-cover transition-transform duration-200 group-hover:scale-105"
      loading="lazy"
    />
    <div class="absolute bottom-2 right-2 rounded bg-black/70 px-2 py-1 text-xs text-white">
      {formatDuration(video.duration || 0)}
    </div>
  </div>
  <div class="p-4">
    <h3 class="line-clamp-2 text-sm font-semibold text-gray-900">{video.title}</h3>
    <p class="mt-1 text-xs text-gray-500">{video.uploader || 'Unknown'}</p>
    <div class="mt-2 flex items-center gap-2 text-xs text-gray-400">
      <span>{new Date(video.createdAt).toLocaleDateString()}</span>
    </div>
  </div>
</div>