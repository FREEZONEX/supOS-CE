import CodeSnippet from '@/components/code-snippet';
import type { FC } from 'react';

const RawData: FC<{ payload?: string }> = ({ payload }) => {
  // 尝试首先解析payload
  let parsedPayload;
  try {
    if (typeof payload === 'string') {
      const firstParse = JSON.parse(payload);
      // 检查第一次解析的结果是否也是字符串，如果是，则尝试第二次解析
      parsedPayload = typeof firstParse === 'string' ? JSON.parse(firstParse) : firstParse;
    } else {
      parsedPayload = payload;
    }
  } catch (error) {
    console.error('Failed to parse payload:', error);
    return null;
  }

  if (!parsedPayload) return null;

  // 转换成美观打印的字符串
  const formattedPayload = JSON.stringify(parsedPayload, null, 2);

  return (
    <CodeSnippet className="codeViewWrap" type="multi" minCollapsedNumberOfRows={1}>
      {formattedPayload}
    </CodeSnippet>
  );
};

export default RawData;
