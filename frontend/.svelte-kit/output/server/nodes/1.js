

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.B9nbHRxz.js","_app/immutable/chunks/DYjCt7Qj.js","_app/immutable/chunks/ZjpciyFr.js","_app/immutable/chunks/j7kPOhj-.js"];
export const stylesheets = [];
export const fonts = [];
