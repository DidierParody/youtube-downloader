<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { user, token } from '$lib/stores.js';
	import { getDownloads, createDownload } from '$lib/api.js';

	let downloads = [];
	let url = '';
	let message = '';

	onMount(async () => {
		if (!$token) {
			goto('/login');
			return;
		}
		const data = await getDownloads($token);
		if (data.downloads) {
			downloads = data.downloads;
		}
	});

	async function handleSubmit() {
		if (!url) return;
		const res = await createDownload(url, $token);
		if (res.message) {
			message = res.message;
			const data = await getDownloads($token);
			if (data.downloads) downloads = data.downloads;
			url = '';
		}
	}
</script>

<svelte:head>
	<title>Dashboard - YouTube Downloader</title>
</svelte:head>

<h2>Your Downloads</h2>

<form on:submit|preventDefault={handleSubmit} class="form">
	<h3>Add New Download</h3>
	<input type="text" bind:value={url} placeholder="YouTube URL" required />
	<button type="submit" class="btn">Download</button>
</form>

{#if message}
	<div class="alert success">{message}</div>
{/if}

<div>
	{#each downloads as dl}
		<div class="card">
			<h4>{dl.title || 'Untitled'}</h4>
			<p>URL: {dl.url}</p>
			<p>Status: <strong>{dl.status}</strong></p>
			<p>Created: {new Date(dl.created_at).toLocaleString()}</p>
		</div>
	{:else}
		<p>No downloads yet.</p>
	{/each}
</div>
