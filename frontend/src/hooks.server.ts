import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

const handleCORS: Handle = async ({ event, resolve }) => {
  const response = await resolve(event);
  
  if (event.url.pathname.startsWith('/api/')) {
    response.headers.set('Access-Control-Allow-Origin', '*');
  }
  
  return response;
};

export const handle = sequence(handleCORS);
// Or simply: export const handle: Handle = async ({ event, resolve }) => { ... }