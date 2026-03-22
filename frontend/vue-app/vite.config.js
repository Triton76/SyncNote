import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/auth': {
        target: 'http://127.0.0.1:8000',
        changeOrigin: true,
      },
      '/api/v1': {
        target: 'http://127.0.0.1:8001',
        changeOrigin: true,
      },
      '/api/user': {
        target: 'http://127.0.0.1:8003',
        changeOrigin: true,
      },
    },
  },
})
