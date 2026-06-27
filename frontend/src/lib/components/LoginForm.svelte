<script lang="ts">
  import { goto } from '$app/navigation';
  import { login } from '$lib/api/auth';
  import { auth } from '$lib/stores/auth';

  let email = $state('');
  let password = $state('');
  let isSubmitting = $state(false);
  let error: string | null = $state(null);

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    isSubmitting = true;
    error = null;

    // Validation
    if (!email.trim() || !password.trim()) {
      error = 'Please fill in all fields';
      isSubmitting = false;
      return;
    }

    try {
      const response = await login({ email, password });
      auth.setToken(response.token, response.user);
      goto('/');
    } catch (err: any) {
      error = err.message || 'Invalid email or password';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="card p-6 sm:p-8">
  <div class="text-center">
    <h2 class="text-2xl font-bold text-gray-900">Sign in to your account</h2>
    <p class="mt-2 text-sm text-gray-600">
      Don't have an account?
      <a href="/auth/register" class="font-semibold text-primary-600 hover:text-primary-500">
        Register
      </a>
    </p>
  </div>

  <form onsubmit={handleSubmit} class="mt-8 space-y-6">
    {#if error}
      <div class="rounded-md bg-red-50 p-4 text-sm text-red-700" role="alert">
        {error}
      </div>
    {/if}

    <div class="space-y-4">
      <div>
        <label for="login-email" class="block text-sm font-medium text-gray-700">Email</label>
        <input
          id="login-email"
          type="email"
          bind:value={email}
          required
          class="input mt-1"
          placeholder="you@example.com"
          disabled={isSubmitting}
        />
      </div>

      <div>
        <label for="login-password" class="block text-sm font-medium text-gray-700">Password</label>
        <input
          id="login-password"
          type="password"
          bind:value={password}
          required
          minlength="6"
          class="input mt-1"
          placeholder="••••••••"
          disabled={isSubmitting}
        />
      </div>
    </div>

    <button
      type="submit"
      disabled={isSubmitting}
      class="btn btn-primary w-full disabled:opacity-50"
    >
      {#if isSubmitting}
        <svg class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        Signing in...
      {:else}
        Sign in
      {/if}
    </button>
  </form>
</div>