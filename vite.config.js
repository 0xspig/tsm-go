import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig({
  root: 'ui/vite',
  build: {
    manifest: true,
    rollupOptions: {
      input: {
        main: 'ui/vite/main.js',
      },
    },
  },
  server: {
    origin: 'http://localhost:3000',
  },
})