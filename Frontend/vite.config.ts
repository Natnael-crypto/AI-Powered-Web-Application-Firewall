import path from 'path'
import react from '@vitejs/plugin-react'
import {defineConfig} from 'vite'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      // Proxy all requests starting with `/api` to the backend
      '/api': {
        target: 'http://localhost:8080', // Your backend server URL
        changeOrigin: true, // Needed for virtual hosted sites
        rewrite: path => path.replace(/^\/api/, ''), // Remove `/api` prefix when forwarding
      },
    },
  },
})
