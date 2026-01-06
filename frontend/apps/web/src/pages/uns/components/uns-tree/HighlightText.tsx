import React from 'react';

interface HighlightTextProps {
  needle: string;
  haystack: string;
}

// 自定义正则表达式转义函数
const escapeRegExp = (str: string) => str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

const HighlightText: React.FC<HighlightTextProps> = ({ needle, haystack }) => {
  // 如果 needle 为空或不存在于 haystack 中，则直接返回原始文本
  if (!needle || !haystack.toLowerCase().includes(needle.toLowerCase())) {
    return <>{haystack}</>;
  }

  // 使用转义后的 needle 创建正则表达式
  const escapedNeedle = escapeRegExp(needle);
  const parts = haystack.split(new RegExp(`(${escapedNeedle})`, 'gi'));

  // 使用 map 方法遍历 parts 数组，为每个匹配项添加 <span>
  const highlightedParts = parts.map((part, index) =>
    new RegExp(needle, 'i').test(part) ? (
      <span key={index} style={{ color: 'var(--supos-theme-color)' }}>
        {part}
      </span>
    ) : (
      part
    )
  );

  return <>{highlightedParts}</>;
};

export default HighlightText;
