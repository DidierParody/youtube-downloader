<script lang="ts">
  import { goto } from '$app/navigation';
  import { register } from '$lib/api/auth';
  import { auth } from '$lib/stores/auth';

  let username = $state('');
  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let isSubmitting = $state(false);
  let error: string | null = $state(null);

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    isSubmitting = true;
    error = null;

    if (!username.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
      error = 'Please fill in all fields';
      isSubmitting = false;
      return;
    }

    if (password !== confirmPassword) {
      error = 'Passwords do not match';
      isSubmitting = false;
      return;
    }

    if (password.length < 6) {
      error = 'Password must be at least 6 characters';
      isSubmitting = false;
      return;
    }

    try {
      const response = await register({ username, email, password });
      auth.setToken(response.token, response.user);
      goto('/');
    } catch (err: any) {
      error = err.message || 'Registration failed. Please try again.';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="card p-6 sm:p-8">
  <div class="text-center">
    <h2 class="text-2xl font-bold text-gray-900">Create an account</h2>
    <p class="mt-2 text-sm text-gray-600">
      Already have an account?
      <a href="/auth/login" class="font-semibold text-primary-600 hover:text-primary-500">
        Sign in
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
        <label for="register-username" class="block text-sm font-medium text-gray-700">Username</label>
        <input
          id="register-username"
          type="text"
          bind:value={username}
          required
          class="input mt-1"
          placeholder="johndoe"
          disabled={isSubmitting}
        />
      </div>

      <div>
        <label for="register-email" class="block text-sm font-medium text-gray-700">Email</label>
        <input
          id="register-email"
          type="email"
          bind:value={email}
          required
          class="input mt-1"
          placeholder="you@example.com"
          disabled={isSubmitting}
        />
      </div>

      <div>
        <label for="register-password" class="block text-sm font-medium text-gray-700">Password</label>
        <input
          id="register-password"
          type="password"
          bind:value={password}
          required
          minlength="6"
          class="input mt-1"
          placeholder="••••••••"
          disabled={isSubmitting}
        />
      </div>

      <div>
        <label for="register-confirm" class="block text-sm font-medium text-gray-700">Confirm Password</label>
        <input
          id="register-confirm"
          type="password"
          bind:value={confirmPassword}
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
        Creating account...
      {:else}
        Create account
      {/if}
    </button>
  </form>
</div>