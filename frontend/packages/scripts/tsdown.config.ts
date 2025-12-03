import { defineConfig } from 'tsdown';

export default defineConfig({
  entry: ['./src/*.ts'],
  platform: 'node',
  format: 'esm', // 输出为 ES 模块，Node.js 现代版本推荐
  outDir: 'dist', // 输出目录
});
