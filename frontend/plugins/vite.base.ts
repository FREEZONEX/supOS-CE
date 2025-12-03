import { defineConfig } from 'vite';
import type { PluginOption } from 'vite';
import react from '@vitejs/plugin-react';
import { federation } from '@module-federation/vite';
import { shared } from './pluginConfig';
import path from 'path';

interface Config {
  assetPrefix?: string; // 资源前缀
  name: string; // 插件名称
  hostOrigin?: string; // 主机地址
  exposes?: { [key: string]: string }; // 暴露的模块
}

const config = ({ assetPrefix, name, hostOrigin, exposes = {} }: Config) => {
  console.log(`当前remote：supos-ce/${name}`, '\n', { assetPrefix, hostOrigin });

  return defineConfig({
    base: assetPrefix || `/plugin/${name}`,
    plugins: [
      react(),
      federation({
        name: `supos-ce/${name}`,
        manifest: true,
        exposes: {
          './index': './src/App.tsx',
          './enUS': './src/locale/en-US.json',
          './zhCN': './src/locale/zh-CN.json',
          ...exposes,
        },
        remotes: {
          '@supos_host': `supos-ce/host@${hostOrigin || ''}/mf-manifest.json`,
        },
        shared: shared,
      }) as PluginOption,
    ],
    resolve: {
      alias: {
        '@': path.resolve(process.cwd(), 'src'),
      },
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx'],
    },
    build: {
      minify: false,
      target: 'chrome89',
      outDir: `./dist/${name}`,
      rollupOptions: {
        input: 'src/main.tsx', // 指定自定义入口文件
      },
    },
    esbuild: {
      target: 'esnext',
    },
    server: {
      origin: assetPrefix,
    },
  });
};

export default config;
