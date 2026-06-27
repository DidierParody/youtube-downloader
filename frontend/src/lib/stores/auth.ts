import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

export interface User {
  id: number;
  username: string;
  email: string;
}

export interface AuthState {
  token: string | null;
  user: User | null;
  isLoading: boolean;
}

const STORAGE_KEY = 'auth_token';

function getInitialState(): AuthState {
  let token: string | null = null;
  if (browser) {
    token = localStorage.getItem(STORAGE_KEY);
  }
  return {
    token,
    user: null,
    isLoading: true
  };
}

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>(getInitialState());

  return {
    subscribe,
    setToken: (token: string, user?: User) => {
      if (browser) {
        localStorage.setItem(STORAGE_KEY, token);
      }
      update(state => ({ ...state, token, user: user ?? state.user }));
    },
    setUser: (user: User) => {
      update(state => ({ ...state, user }));
    },
    clearAuth: () => {
      if (browser) {
        localStorage.removeItem(STORAGE_KEY);
      }
      set({ token: null, user: null, isLoading: false });
    },
    setLoading: (isLoading: boolean) => {
      update(state => ({ ...state, isLoading }));
    }
  };
}

export const auth = createAuthStore();
export const isAuthenticated = derived(auth, $auth => !!$auth.token);
export const currentUser = derived(auth, $auth => $auth.user);