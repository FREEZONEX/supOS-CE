export function formatStreamSQL(sql: string): string {
  // 移除多余的空格和换行
  let formatted = sql.trim().replace(/\s+/g, ' ');

  // 主要关键字
  const keywords = [
    'CREATE STREAM',
    'TRIGGER',
    'WATERMARK',
    'FILL HISTORY',
    'IGNORE EXPIRED',
    'IGNORE UPDATE',
    'INTO',
    'SELECT',
    'FROM',
    'EVENT_WINDOW',
    'START WITH',
    'END WITH',
    'WHERE',
    'GROUP BY',
  ];

  // 在关键字前添加换行和缩进
  keywords.forEach((keyword) => {
    formatted = formatted.replace(new RegExp(`\\s*${keyword}\\s+`, 'gi'), `\n  ${keyword} `);
  });

  // 处理括号
  formatted = formatted.replace(/\(/g, '(').replace(/\)/g, ')').replace(/,\s*/g, ',\n  ');

  // 处理分号
  formatted = formatted.replace(/;/g, ';\n');

  // 移除多余的空行
  formatted = formatted
    .split('\n')
    .map((line) => line.trim())
    .filter((line) => line.length > 0)
    .join('\n');

  return formatted;
}
