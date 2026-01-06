import fs from 'fs';
import translateText from './tmtClient.js';

const prefixPath = process.cwd();

// 语言配置项
const LANGUAGES = ['zh-CN', 'en-US'] as const;

interface LanguageMessages {
  [key: string]: string;
}

const loadLanguageFile = (key: string): LanguageMessages => {
  const filePath = `${prefixPath}/src/locale/${key}.json`.replace(/\\/g, '/');
  return fs.existsSync(filePath) ? JSON.parse(fs.readFileSync(filePath, 'utf-8')) : {};
};

const oldMessages: Record<string, LanguageMessages> = {}; // 旧的语言文件
const newMessages: Record<string, LanguageMessages> = {}; // 新的语言文件

// 根据配置加载语言文件
LANGUAGES.forEach((langKey) => {
  oldMessages[langKey] = loadLanguageFile(langKey);
});

// 遍历messages，将其写入语言文件
const intl = async (): Promise<void> => {
  const filePath = `${prefixPath}/src/locale/index.js`.replace(/\\/g, '/');
  const { default: messages } = await import(`file://${filePath}`);

  for (const key of Object.keys(messages)) {
    for (const langKey of LANGUAGES) {
      if (!newMessages[langKey]) {
        newMessages[langKey] = {};
      }

      // 如果是中文，直接写入
      if (langKey === 'zh-CN') {
        newMessages[langKey][key] = messages[key];
        continue;
      }

      // 如果不存在，翻译并写入
      if (!oldMessages[langKey]?.[key]) {
        newMessages[langKey][key] = await translateText(messages[key], langKey);
        continue;
      }

      // 已存在，直接复用
      newMessages[langKey][key] = oldMessages[langKey][key];
    }
  }

  // 生成新的语言文件
  LANGUAGES.forEach((langKey) => {
    const outputPath = `${prefixPath}/src/locale/${langKey}.json`.replace(/\\/g, '/');
    const newContent = JSON.stringify(newMessages[langKey], null, 2) + '\n';

    // 检查文件是否存在且内容是否相同
    if (fs.existsSync(outputPath)) {
      const oldContent = fs.readFileSync(outputPath, 'utf-8');
      // 规范化行分隔符进行比较（处理 Windows 的 \r\n 和 Unix 的 \n 差异）
      const normalizedOldContent = oldContent.replace(/\r\n/g, '\n');
      const normalizedNewContent = newContent.replace(/\r\n/g, '\n');

      if (normalizedOldContent === normalizedNewContent) {
        console.log(`语言文件 ${langKey}.json 无变更，跳过写入`);
        return;
      }
    }

    fs.writeFileSync(outputPath, newContent);
    console.log(`语言文件 ${langKey}.json 已更新`);
  });
};

// 执行异步处理
intl().catch(console.error);
