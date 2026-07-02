<script>
	import { goto } from '$app/navigation';
	import { user, token } from '$lib/stores.js';
	import { register } from '$lib/api.js';

	let username = '';
	let email = '';
	let password = '';
	let error = '';
	let success = '';

	async function handleSubmit() {
		const res = await register(username, email, password);
		if (res.message) {
			success = res.message;
			setTimeout(() => goto('/login'), 1500);
		} else {
			error = res.error || 'Registration failed';
		}
	}
</script>

<svelte:head>
	<title>Register - YouTube Downloader</title>
</svelte:head>

<h2>Register</h2>

<form on:submit|preventDefault={handleSubmit} class="form">
	{#if success}
		<div class="alert success">{success}</div>
	{/if}
	{#if error}
		<div class="alert error">{error}</div>
	{/if}
	<input type="text" bind:value={username} placeholder="Username" required />
	<input type="email" bind:value={email} placeholder="Email" required />
	<input type="password" bind:value={password} placeholder="Password" required />
	<button type="submit" class="btn">Register</button>
	<p>Already have an account? <a href="/login">Login</a></p>
</form>
