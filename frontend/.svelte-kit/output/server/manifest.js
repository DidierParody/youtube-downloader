export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set([]),
	mimeTypes: {},
	_: {
		client: {start:"_app/immutable/entry/start.CckIzgOn.js",app:"_app/immutable/entry/app.C3mwfFqp.js",imports:["_app/immutable/entry/start.CckIzgOn.js","_app/immutable/chunks/j7kPOhj-.js","_app/immutable/chunks/DYjCt7Qj.js","_app/immutable/entry/app.C3mwfFqp.js","_app/immutable/chunks/DYjCt7Qj.js","_app/immutable/chunks/ZjpciyFr.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js'))
		],
		remotes: {
			
		},
		routes: [
			
		],
		prerendered_routes: new Set(["/","/login","/register"]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
