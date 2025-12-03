import dayjs from 'dayjs';

/**
 * 日期格式化说明：
 *
 * - 固定格式：使用明确占位符（如 HH/mm/dd），不依赖语言环境
 * * 示例：'HH:mm:ss DD/MM/YYYY' → "14:05:00 20/02/2025"
 *
 * - 动态格式：包含本地化标记（如 MMMM/dddd），随 locale 变化
 * * 示例：'dddd, MMMM D, YYYY' → 中文环境显示"星期四, 二月 20, 2025"
 * * 英文环境显示"Thursday, February 20, 2025"
 */
const DAY_FORMATS: { [key: string]: any } = {
  all: {
    en: 'HH:mm:ss DD/MM/YYYY',
    'zh-cn': 'YYYY/MM/DD HH:mm:ss',
  },
};

/**
 * @param timestamp 时间戳
 * @param format 格式化类型
 * @param customize 是否自定义格式化
 * @returns 格式化后的时间字符串
 */
export const formatTimestamp = (timestamp: any, format = 'all', customize = false) => {
  const locale = dayjs.locale();
  if (!timestamp) return '';
  if (customize) return dayjs(timestamp).format(format === 'all' ? 'HH:mm:ss DD/MM/YYYY' : format);
  return dayjs(timestamp).format(DAY_FORMATS[format]?.[locale ?? 'en'] ?? 'HH:mm:ss DD/MM/YYYY');
};

export const simpleFormat = (data: any) => {
  if (data == null) return ''; // 处理 null 和 undefined
  if (typeof data === 'boolean') return String(data); // 处理布尔值，转换为字符串
  if (typeof data === 'object') return JSON.stringify(data) || ''; // 处理数组和对象
  return data; // 处理数字、字符串和其他基本类型
};

/**
 * @param item 字符串
 * @returns 字符串里面时间戳格式化替换后最终字符串
 */
export const processTimestamp = (item: string) => {
  const match = item.match(/(\d{1,3}(?:,\d{3})*) /);
  if (match && match[1]) {
    const timestamp = match[1].replace(/,/g, '');
    const formattedTime: any = formatTimestamp(Number(timestamp));
    const newStr = item.replace(match[1], formattedTime);
    return newStr;
  }
  return item;
};

/**
 * @param item 字符串
 * @returns 转化国际化方法
 */
export const formatShowName = ({
  code,
  showName,
  formatMessage,
  finallyShowName,
}: {
  code?: string;
  showName?: string;
  formatMessage: any;
  finallyShowName?: string;
}) => {
  // 如果code和showName一样，表明后端翻译未成功
  if (showName && code !== showName) return showName;
  if (code) {
    try {
      return formatMessage(code);
    } catch (e) {
      console.log(e);
      return code;
    }
  }
  return finallyShowName || '未配置国际化';
};
