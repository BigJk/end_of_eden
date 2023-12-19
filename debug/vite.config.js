import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';
import path from 'path';

export default defineConfig({
	plugins: [tsconfigPaths()],
	root: 'src/',
	resolve: {
		alias: {
			js: path.resolve(__dirname, 'src/js'),
			src: path.resolve(__dirname, 'src'),
		},
	},
	build: {
		outDir: '../dist',
	},
	server: {
		host: '127.0.0.1',
		port: 3000,
	},
});