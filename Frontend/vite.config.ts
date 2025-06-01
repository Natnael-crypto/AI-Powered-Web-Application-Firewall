import path from 'path';
import react from '@vitejs/plugin-react';
import { defineConfig, loadEnv } from 'vite';
import fs from 'fs'

const key = fs.readFileSync('./certs/key.pem')
const cert = fs.readFileSync('./certs/cert.pem')

export default defineConfig(({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      https: {
        key: key,
        cert: cert,
      },
      proxy: {
        '/api': {
          target: process.env.VITE_BACKEND_URL,
          changeOrigin: true,
          secure: false,
          rewrite: path => path.replace(/^\/api/, ''),
          headers: {
            Connection: 'keep-alive',
          },
          cookieDomainRewrite: '',
        },
      }
    },
  };
});
