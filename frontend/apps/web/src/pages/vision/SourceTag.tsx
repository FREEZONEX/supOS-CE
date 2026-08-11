import { useTranslate } from '@/hooks';
import type { AlgorithmSource } from '@/apis/core-api/algorithm';
import { SOURCE_META } from './algorithm-meta';
import styles from './index.module.scss';

// 来源标签：内置算法用蓝色胶囊，自定义走中性灰（对齐 Figma 13632-148096）。
const SourceTag = ({ source }: { source: AlgorithmSource }) => {
  const formatMessage = useTranslate();
  const meta = SOURCE_META[source] || SOURCE_META.custom;
  const className = source === 'builtin' ? `${styles.sourceTag} ${styles.sourceTagBuiltin}` : styles.sourceTag;
  return <span className={className}>{formatMessage(meta.labelKey)}</span>;
};

export default SourceTag;
