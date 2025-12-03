import fs from 'fs';
import path from 'path';

interface Config {
  src: string;
  dist: string;
}

// 解析命令行参数
const args = process.argv.slice(2);
const prefix = args.find((arg) => arg?.startsWith('--prefix='))?.split('=')[1];
const configPath = args.find((arg) => arg?.startsWith('--config='))?.split('=')[1];

// 读取配置文件
const config: Config = configPath
  ? JSON.parse(fs.readFileSync(configPath, 'utf-8'))
  : {
      src: './src/locale',
      dist: './properties-output',
    };

const inputDir = path.resolve(config.src);
const outputDir = path.resolve(config.dist);

// 确保输出目录存在
fs.mkdirSync(outputDir, { recursive: true });

// 读取目录下所有JSON文件
const files = fs.readdirSync(inputDir).filter((file) => file.endsWith('.json'));

files.forEach((file) => {
  const inputPath = path.join(inputDir, file);
  const outputPath = path.join(outputDir, file.replace('.json', '.properties'));

  // 读取JSON文件
  const jsonData = JSON.parse(fs.readFileSync(inputPath, 'utf-8'));

  // 转换为properties格式并添加前缀
  let propertiesContent = '';

  function convertJsonToProperties(obj: Record<string, any>, parentKey = ''): void {
    for (const [key, value] of Object.entries(obj)) {
      const currentKey = parentKey ? `${parentKey}.${key}` : key;
      if (typeof value === 'object' && value !== null) {
        convertJsonToProperties(value, currentKey);
      } else {
        // 添加前缀
        const prefixedKey = prefix ? `${prefix}.${currentKey}` : currentKey;
        propertiesContent += `${prefixedKey}=${value}\n`;
      }
    }
  }

  convertJsonToProperties(jsonData);

  // 写入properties文件
  fs.writeFileSync(outputPath, propertiesContent);
  console.log(`Successfully converted ${file} to ${outputPath} with prefix: ${prefix}`);
});
