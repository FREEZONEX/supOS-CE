import { defineConfig } from 'tsdown';

export default defineConfig({
  entry: ['./src/index.ts'],
  platform: 'node',
  format: 'esm', // 输出为 ES 模块，Node.js 现代版本推荐
  outDir: 'dist', // 输出目录
  unbundle: true, // 输出目录和源文件木有一一对应 https://tsdown.dev/zh-CN/options/unbundle
  minify: true,
  external: ['dotenv', '@modelcontextprotocol/sdk'],
  watch: process.argv.includes('--watch') ? ['./src/**/*.ts'] : false,
});
