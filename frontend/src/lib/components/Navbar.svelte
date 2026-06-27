<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { isAuthenticated } from '$lib/stores/auth';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  let menuOpen = $state(false);

  function toggleMenu() {
    menuOpen = !menuOpen;
  }

  function logout() {
    auth.clearAuth();
    goto('/auth/login');
  }

  const navLinks = [
    { href: '/', label: 'Home' },
    { href: '/downloads', label: 'Downloads' },
    { href: '/videos', label: 'Videos' }
  ];
</script>

<nav class="bg-white shadow-sm">
  <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
    <div class="flex h-16 justify-between">
      <!-- Logo -->
      <div class="flex shrink-0 items-center">
        <a href="/" class="text-xl font-bold text-primary-600">YouTube Downloader</a>
      </div>

      <!-- Desktop navigation -->
      <div class="hidden items-center space-x-8 md:flex">
        {#each navLinks as link}
          <a
            href={link.href}
            class="text-sm font-medium transition-colors {$page.url.pathname === link.href
              ? 'text-primary-600'
              : 'text-gray-500 hover:text-gray-900'}"
          >
            {link.label}
          </a>
        {/each}

        {#if $isAuthenticated}
          <button
            onclick={logout}
            class="rounded-md bg-gray-100 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200"
          >
            Logout
          </button>
        {:else}
          <a
            href="/auth/login"
            class="rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700"
          >
            Login
          </a>
        {/if}
      </div>

      <!-- Mobile menu button -->
      <div class="flex items-center md:hidden">
        <button
          onclick={toggleMenu}
          class="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
          aria-expanded={menuOpen}
        >
          <span class="sr-only">Open main menu</span>
          {#if menuOpen}
            <svg class="block h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          {:else}
            <svg class="block h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
            </svg>
          {/if}
        </button>
      </div>
    </div>
  </div>

  <!-- Mobile menu -->
  {#if menuOpen}
    <div class="md:hidden">
      <div class="space-y-1 px-2 pb-3 pt-2">
        {#each navLinks as link}
          <a
            href={link.href}
            class="block rounded-md px-3 py-2 text-base font-medium {$page.url.pathname === link.href
              ? 'bg-primary-50 text-primary-700'
              : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'}"
            onclick={() => (menuOpen = false)}
          >
            {link.label}
          </a>
        {/each}
        {#if $isAuthenticated}
          <button
            onclick={logout}
            class="block w-full rounded-md px-3 py-2 text-left text-base font-medium text-gray-700 hover:bg-gray-50"
          >
            Logout
          </button>
        {:else}
          <a
            href="/auth/login"
            class="block rounded-md px-3 py-2 text-base font-medium text-primary-600 hover:bg-primary-50"
          >
            Login
          </a>
        {/if}
      </div>
    </div>
  {/if}
</nav>