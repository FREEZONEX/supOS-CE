import { execSync } from 'child_process';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import fs from 'fs';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

export function runScript(scriptName: string): void {
  const cwd = process.cwd();
  switch (scriptName) {
    // 翻译一次
    case 'intl:once':
      execSync('node ' + join(__dirname, '../dist/i18n.js'), { stdio: 'inherit' });
      break;
    // 自动化翻译
    case 'intl:watch':
      {
        // 使用 Node.js 内置的文件监视功能，避免依赖 nodemon
        const watchFile = join(cwd, 'src/locale/index.js');

        if (!fs.existsSync(watchFile)) {
          console.error(`监视文件不存在: ${watchFile}`);
          process.exit(1);
        }

        console.log(`开始监视文件: ${watchFile}`);
        console.log('按 Ctrl+C 停止监视');

        // 立即执行一次
        execSync('node ' + join(__dirname, '../dist/i18n.js'), { stdio: 'inherit' });

        // 设置文件监视
        fs.watchFile(watchFile, { interval: 1000 }, (curr, prev) => {
          if (curr.mtime !== prev.mtime) {
            console.log('检测到文件变更，重新执行国际化处理...');
            try {
              execSync('node ' + join(__dirname, '../dist/i18n.js'), { stdio: 'inherit' });
            } catch (error) {
              console.error('执行失败:', error);
            }
          }
        });

        // 保持进程运行
        process.on('SIGINT', () => {
          console.log('\n停止监视');
          process.exit(0);
        });
      }

      break;
    // 打包
    case 'build':
      execSync('npm run clean && tsc -b && vite build', {
        stdio: 'inherit',
        cwd,
      });
      execSync('node ' + join(__dirname, '../dist/afterBuild.js'), {
        stdio: 'inherit',
        cwd,
      });
      break;
    case 'json-to-properties':
      {
        const args = process.argv.slice(3).join(' ');
        execSync(`node ${join(__dirname, '../dist/jsonToProperties.js')} ${args}`, {
          stdio: 'inherit',
          cwd,
        });
      }
      break;
    default:
      throw new Error(`Unknown script: ${scriptName}`);
  }
}
