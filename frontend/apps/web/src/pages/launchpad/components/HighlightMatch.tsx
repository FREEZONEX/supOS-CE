import { type FC } from 'react';
import styles from './HighlightMatch.module.scss';

interface HighlightMatchProps {
  /** 展示文本 */
  text: string;
  /** 搜索关键词，命中子串保持原色、其余置灰 */
  keyword?: string;
}

/**
 * 搜索输入态：命中的子串保持原色，其余字符置灰（对齐设计稿）。
 * 空关键词或未命中时原样返回。
 */
const HighlightMatch: FC<HighlightMatchProps> = ({ text, keyword }) => {
  const trimmed = keyword?.trim();
  if (!trimmed) {
    return <>{text}</>;
  }

  const index = text.toLowerCase().indexOf(trimmed.toLowerCase());
  if (index === -1) {
    return <>{text}</>;
  }

  return (
    <>
      <span className={styles.dim}>{text.slice(0, index)}</span>
      {text.slice(index, index + trimmed.length)}
      <span className={styles.dim}>{text.slice(index + trimmed.length)}</span>
    </>
  );
};

export default HighlightMatch;
