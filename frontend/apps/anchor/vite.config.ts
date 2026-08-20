import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// 子应用统一由 backend 静态承载在同源 /anchor/ 路径（iframe 嵌入），base 必须与之一致
export default defineConfig({
  base: '/anchor/',
  plugins: [react(), tailwindcss()],
  server: {
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: process.env.API_PROXY_URL || 'http://127.0.0.1:8088',
        changeOrigin: true,
      },
    },
  },
});
