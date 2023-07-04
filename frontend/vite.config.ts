import fs from 'fs';
import path from 'path';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [sveltekit()],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	},
	server: {
		port: 1420,
		https: {
			key: fs.readFileSync(path.join('./certs', 'local-key.pem')),
			cert: fs.readFileSync(path.join('./certs', 'local-cert.pem'))
		}
	}
});
