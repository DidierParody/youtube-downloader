<script>
	import { goto } from '$app/navigation';
	import { user, token } from '$lib/stores.js';
	import { login } from '$lib/api.js';

	let email = '';
	let password = '';
	let error = '';

	async function handleSubmit() {
		const res = await login(email, password);
		if (res.token) {
			token.set(res.token);
			user.set(res.user);
			goto('/');
		} else {
			error = res.error || 'Login failed';
		}
	}
</script>

<svelte:head>
	<title>Login - YouTube Downloader</title>
</svelte:head>

<h2>Login</h2>

<form on:submit|preventDefault={handleSubmit} class="form">
	{#if error}
		<div class="alert error">{error}</div>
	{/if}
	<input type="email" bind:value={email} placeholder="Email" required />
	<input type="password" bind:value={password} placeholder="Password" required />
	<button type="submit" class="btn">Login</button>
	<p>Don't have an account? <a href="/register">Register</a></p>
</form>
