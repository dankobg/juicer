// import path from 'node:path';
// import fs from 'node:fs';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { purgeCss } from 'vite-plugin-tailwind-purgecss';

export default defineConfig({
	plugins: [sveltekit(), purgeCss()],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
	},
	server: {
		port: 3974,
		// https: {
		// 	key: fs.readFileSync(path.join('./certs', 'local-key.pem')),
		// 	cert: fs.readFileSync(path.join('./certs', 'local-cert.pem')),
		// },
	},
});
