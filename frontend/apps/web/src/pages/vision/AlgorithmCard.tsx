import { Button, Dropdown, Tooltip, Typography } from 'antd';
import { OverflowMenuVertical, TrashCan } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import type { VisionAlgorithm } from '@/apis/core-api/algorithm';
import { SOURCE_META, fitAlgorithmLabels } from './algorithm-meta';
import ModelStatusTag from './ModelStatusTag';
import styles from './index.module.scss';

const { Text, Paragraph } = Typography;

type AlgorithmCardProps = {
  algorithm: VisionAlgorithm;
  canManage?: boolean;
  onDetail: (algorithm: VisionAlgorithm) => void;
  onUpdateModel: (algorithm: VisionAlgorithm) => void;
  onDelete: (algorithm: VisionAlgorithm) => void;
};

const AlgorithmCard = ({ algorithm, canManage, onDetail, onUpdateModel, onDelete }: AlgorithmCardProps) => {
  const formatMessage = useTranslate();
  const sourceMeta = SOURCE_META[algorithm.source] || SOURCE_META.custom;
  const labels = algorithm.labels || [];
  const { shown: shownLabels, rest: remainingLabels, all: allLabels } = fitAlgorithmLabels(labels, 332);

  return (
    <div className={styles.algoCard}>
      <div className={styles.algoCardHead}>
        <Text ellipsis={{ tooltip: algorithm.name }} className={styles.algoCardTitle}>
          {algorithm.name}
        </Text>
        <span className={styles.algoCardStatus}>
          <ModelStatusTag status={algorithm.modelStatus} />
        </span>
      </div>
      <Text type="secondary" className={styles.algoCardSource}>
        {formatMessage('Vision.algorithm.source')}: {formatMessage(sourceMeta.labelKey)}
      </Text>
      <Paragraph
        type="secondary"
        ellipsis={{ rows: 2, tooltip: algorithm.description }}
        className={styles.algoCardDesc}
      >
        {algorithm.description || '-'}
      </Paragraph>
      <div className={styles.algoCardLabels}>
        <Text type="secondary" className={styles.algoCardLabelsTitle}>
          {formatMessage('Vision.algorithm.labels')}:
        </Text>
        <div className={styles.algoCardLabelRow}>
          {labels.length === 0 ? (
            <Text type="secondary">--</Text>
          ) : (
<>
              {shownLabels.map((label) => (
                <span key={label} className={styles.algoCardLabelTag}>
                  {label}
                </span>
              ))}
              {remainingLabels > 0 && (
                <Tooltip title={allLabels.join(', ')}>
                  <span className={`${styles.algoCardLabelTag} ${styles.algoCardLabelMore}`}>+{remainingLabels}</span>
                </Tooltip>
              )}
            </>
          )}
        </div>
      </div>
      <div className={styles.algoCardDivider} />
      <div className={styles.algoCardFooter}>
        <Button type="primary" onClick={() => onDetail(algorithm)}>
          {formatMessage('Vision.algorithm.details')}
        </Button>
        {canManage && (
          <Button onClick={() => onUpdateModel(algorithm)}>{formatMessage('Vision.algorithm.uploadModel')}</Button>
        )}
        {canManage && (
          <Dropdown
            trigger={['click']}
            menu={{
              items: [
                {
                  key: 'delete',
                  danger: true,
                  // 内置算法不允许删除,置灰并说明原因(与表格的操作菜单一致)
                  disabled: algorithm.source !== 'custom',
                  icon: <TrashCan size={16} />,
                  label:
                    algorithm.source === 'custom' ? (
                      formatMessage('common.delete')
                    ) : (
                      <Tooltip title={formatMessage('Vision.algorithm.builtinNoDelete')}>
                        {formatMessage('common.delete')}
                      </Tooltip>
                    ),
                  onClick: () => onDelete(algorithm),
                },
              ],
            }}
          >
            <button type="button" className={styles.algoCardMore} aria-label={formatMessage('common.operation')}>
              <OverflowMenuVertical size={16} />
            </button>
          </Dropdown>
        )}
      </div>
    </div>
  );
};

export default AlgorithmCard;
