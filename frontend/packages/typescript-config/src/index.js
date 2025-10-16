import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';
import { readFileSync } from 'fs';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// 读取配置文件
function readConfig(name) {
  const path = resolve(__dirname, `${name}.json`);
  return JSON.parse(readFileSync(path, 'utf8'));
}

// 导出所有配置
export const base = readConfig('base');
export const react = readConfig('react');
export const node = readConfig('node');

export default {
  base,
  react,
  node,
};
