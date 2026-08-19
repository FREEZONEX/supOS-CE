import { Tag } from 'antd';
import { useTranslate } from '@/hooks';
import type { AlgorithmModelStatus } from '@/apis/core-api/algorithm';
import { MODEL_STATUS_META } from './algorithm-meta';
import styles from './index.module.scss';

// 模型状态标签：missing 用白底中性样式，available/error 用彩色 Tag。
const ModelStatusTag = ({ status }: { status: AlgorithmModelStatus }) => {
  const formatMessage = useTranslate();
  const meta = MODEL_STATUS_META[status] || MODEL_STATUS_META.missing;
  if (meta.color === 'default') {
    return <span className={`${styles.algoStatusNeutral} ${styles.pillTag}`}>{formatMessage(meta.labelKey)}</span>;
  }
  return (
    <Tag color={meta.color} bordered={false} className={styles.pillTag}>
      {formatMessage(meta.labelKey)}
    </Tag>
  );
};

export default ModelStatusTag;
