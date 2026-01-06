#!/usr/bin/env node

import { existsSync } from 'fs';
import { join } from 'path';
import { fileURLToPath } from 'url';
import { execSync } from 'child_process';

const __filename = fileURLToPath(import.meta.url);
const __dirname = join(__filename, '..');

const scriptsPackagePath = join(__dirname, '..', 'packages', 'scripts');
const distPath = join(scriptsPackagePath, 'dist');

// 检查dist目录是否存在
if (!existsSync(distPath)) {
  console.log('scripts包没有dist目录，开始构建...');

  try {
    // 切换到scripts包目录并执行构建
    execSync('pnpm build:scripts', {
      stdio: 'inherit',
      shell: true,
    });
    console.log('scripts包构建完成！');
  } catch (error) {
    console.error('scripts包构建失败:', error.message);
    process.exit(1);
  }
} else {
  console.log('scripts包已有dist目录，跳过构建');
}
