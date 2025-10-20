import { defineConfig } from 'tsdown';

export default defineConfig({
  entry: ['./src/index.ts'],
  platform: 'node',
  format: 'esm',
  outDir: 'dist',
  unbundle: true,
  minify: true,
  external: ['dotenv'],
  watch: process.argv.includes('--watch') ? ['./src/**/*.ts'] : false,
});
