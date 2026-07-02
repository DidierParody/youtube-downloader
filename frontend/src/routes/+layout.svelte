<script>
	import '../app.css';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { user, token } from '$lib/stores.js';

	onMount(() => {
		if (typeof window !== 'undefined') {
			const savedToken = localStorage.getItem('token');
			const savedUser = localStorage.getItem('user');
			if (savedToken) token.set(savedToken);
			if (savedUser) user.set(JSON.parse(savedUser));
		}
	});

	function logout() {
		user.set(null);
		token.set(null);
		localStorage.removeItem('token');
		localStorage.removeItem('user');
		goto('/login');
	}
</script>

<header>
	<h1>🎬 YouTube Downloader</h1>
	{#if $user}
		<button class="btn" on:click={logout}>Logout ({$user.username})</button>
	{/if}
</header>

<main class="container">
	<slot />
</main>
