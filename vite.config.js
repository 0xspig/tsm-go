import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig({
  root: 'ui/vite',
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'ui/vite/index.html'),
      },
    },
  },
})