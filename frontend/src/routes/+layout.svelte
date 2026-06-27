<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { browser } from '$app/environment'
  import { goto } from '$app/navigation';
  import Navbar from '$lib/components/Navbar.svelte';
  import { isAuthenticated } from '$lib/stores/auth';

  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();

  $effect(() => {
    if (browser) {
      const path = $page.url.pathname;
      const requiresAuth = ['/downloads', '/videos'].some((prefix) =>
        path.startsWith(prefix)
      );
      const isAuthPage = path.startsWith('/auth/');

      if (requiresAuth && !$isAuthenticated) {
        goto('/auth/login');
      } else if (isAuthPage && $isAuthenticated) {
        goto('/');
      }
    }
  });
</script>

<div class="min-h-screen bg-gray-50">
  <Navbar />
  <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
    {#if children}
      {@render children()}
    {/if}
  </main>
</div>