const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:3000';

export async function register(username, email, password) {
	const res = await fetch(`${API_BASE}/api/v1/auth/register`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username, email, password })
	});
	return res.json();
}

export async function login(email, password) {
	const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, password })
	});
	return res.json();
}

export async function getDownloads(token) {
	const res = await fetch(`${API_BASE}/api/v1/downloads`, {
		headers: { 'Authorization': `Bearer ${token}` }
	});
	return res.json();
}

export async function createDownload(url, token) {
	const res = await fetch(`${API_BASE}/api/v1/downloads`, {
		method: 'POST',
		headers: { 
			'Content-Type': 'application/json',
			'Authorization': `Bearer ${token}`
		},
		body: JSON.stringify({ url })
	});
	return res.json();
}
