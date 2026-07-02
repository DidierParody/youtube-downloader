import { writable } from 'svelte/store';

export const user = writable(null);
export const token = writable(typeof window !== 'undefined' ? localStorage.getItem('token') : null);

user.subscribe((u) => {
	if (typeof window !== 'undefined' && u) {
		localStorage.setItem('user', JSON.stringify(u));
	}
});

token.subscribe((t) => {
	if (typeof window !== 'undefined' && t) {
		localStorage.setItem('token', t);
	}
});
