import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			// Tauri expects a static build output
			fallback: 'index.html'
		}),
		// Disable SSR – Tauri loads the app as a static file bundle
		prerender: {
			handleHttpError: 'warn'
		}
	}
};

export default config;
