import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';

// https://vitejs.dev/config/
export default defineConfig(async () => ({
	plugins: [sveltekit()],

	// Vite options tailored for Tauri development and production
	clearScreen: false,
	// Tauri expects a fixed port; fail if that port is not available
	server: {
		port: 1420,
		strictPort: true,
		watch: {
			// Tell vite to ignore watching `src-tauri`
			ignored: ['**/src-tauri/**']
		}
	}
}));
