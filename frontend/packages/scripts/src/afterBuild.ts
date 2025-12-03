import fs from 'fs';
import path from 'path';

const prefixPath = process.cwd();

// 源目录和目标目录路径
const sourceDir = `${prefixPath}/dist`;
const targetDir = `${prefixPath}/../../dist`;

// 确保目标目录存在
if (!fs.existsSync(targetDir)) {
  fs.mkdirSync(targetDir, { recursive: true });
}

// 复制目录函数
function copyDirSync(source: string, target: string): void {
  const files = fs.readdirSync(source);
  files.forEach((file) => {
    const sourcePath = path.join(source, file);
    const targetPath = path.join(target, file);

    if (fs.lstatSync(sourcePath).isDirectory()) {
      if (!fs.existsSync(targetPath)) {
        fs.mkdirSync(targetPath);
      }
      copyDirSync(sourcePath, targetPath);
    } else {
      fs.copyFileSync(sourcePath, targetPath);
    }
  });
}

// 执行复制
copyDirSync(sourceDir, targetDir);
console.log('构建成功');
