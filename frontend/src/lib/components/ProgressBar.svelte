<script lang="ts">
  interface Props {
    progress: number; // 0 to 100
    status?: string;
    size?: string;
  }

  let { progress, status, size = 'md' }: Props = $props();

  const percentage = Math.min(Math.max(progress, 0), 100);
  
  const sizeClasses = {
    sm: 'h-1.5',
    md: 'h-2',
    lg: 'h-3'
  };
</script>

<div class="w-full">
  <div class={`w-full rounded-full bg-gray-200 ${sizeClasses[size as keyof typeof sizeClasses] || sizeClasses.md}`}>
    <div
      class="progress-bar h-full rounded-full bg-primary-600 transition-all duration-300 ease-out"
      style="width: {percentage}%"
      role="progressbar"
      aria-valuenow={percentage}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-label={status || `Progress: ${percentage.toFixed(0)}%`}
    ></div>
  </div>
  {#if status}
    <div class="mt-1 flex justify-between text-xs text-gray-500">
      <span>{status}</span>
      <span>{percentage.toFixed(1)}%</span>
    </div>
  {/if}
</div>