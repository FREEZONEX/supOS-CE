import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import legacy from '@vitejs/plugin-legacy';
// import { federation } from '@module-federation/vite';
import path from 'path';
import packageJson from './package.json';
import { getDevInfo, getProxy, logBuildTime, logDevInfo } from './config/supos.dev';
// import { mfConfig } from './config/mfConfig.ts';

const devInfo = getDevInfo();
const proxy = getProxy(devInfo.API_PROXY_URL, devInfo.SINGLE_API_PROXY_LIST, devInfo.SINGLE_API_PROXY_URL);
logDevInfo(devInfo);
const buildTime = logBuildTime();

// https://vite.dev/config/
export default defineConfig({
  base: devInfo.VITE_ASSET_PREFIX || '/',
  esbuild: {
    drop: ['debugger'],
    pure: ['console.log'],
    supported: {
      'top-level-await': true,
    },
  },
  plugins: [
    react(),
    legacy({
      targets: ['chrome>=89', 'safari>=15', 'firefox>=89', 'edge>=89'],
      modernPolyfills: true,
    }),
    // federation(mfConfig),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
    //导入文件时省略的扩展名
    extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx'],
  },
  define: {
    'process.env': { ...devInfo },
    'import.meta.env.VITE_APP_VERSION': JSON.stringify(packageJson.version),
    'import.meta.env.VITE_APP_BUILD_TIMESTAMP': JSON.stringify(buildTime),
  },
  envPrefix: ['REACT_APP_', 'VITE_', 'OPENAI_'],
  server: {
    origin: devInfo.VITE_ASSET_PREFIX,
    proxy: {
      ...proxy,
      '/copilotkit': 'http://localhost:4000',
      '/open-api': 'http://localhost:4000',
      ...(devInfo.VITE_ASSET_PREFIX !== '1'
        ? {
            '/plugin/': {
              target: devInfo.API_PROXY_URL,
              changeOrigin: true,
            },
          }
        : {
            '/mf-manifest.json': devInfo.VITE_ASSET_PREFIX,
          }),
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', 'react-router'],
          antd: ['antd', '@ant-design/icons'],
          charts: ['@antv/x6'],
          utils: ['ahooks', 'lodash-es', 'dayjs'],
        },
      },
    },
    target: ['chrome89', 'edge89', 'firefox89', 'safari15'],
  },
});
