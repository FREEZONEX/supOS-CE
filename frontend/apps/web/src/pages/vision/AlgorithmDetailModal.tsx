import type { ReactNode } from 'react';
import { Button, Space, Typography } from 'antd';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import type { VisionAlgorithm } from '@/apis/core-api/algorithm';
import ModelStatusTag from './ModelStatusTag';
import SourceTag from './SourceTag';
import { formatAlgorithmLabel } from './algorithm-meta';
import styles from './index.module.scss';

const { Text } = Typography;

type AlgorithmDetailModalProps = {
  open: boolean;
  algorithm: VisionAlgorithm | null;
  canManage?: boolean;
  onEdit: (algorithm: VisionAlgorithm) => void;
  onClose: () => void;
};

const AlgorithmDetailModal = ({ open, algorithm, canManage, onEdit, onClose }: AlgorithmDetailModalProps) => {
  const formatMessage = useTranslate();
  if (!algorithm) return null;

  const renderTags = (values: string[], tagClass: string) =>
    values.length === 0 ? (
      <Text type="secondary">--</Text>
    ) : (
      <span className={styles.labelTags}>
        {values.map((value) => (
          <span key={value} className={tagClass}>
            {formatAlgorithmLabel(value)}
          </span>
        ))}
      </span>
    );

  const items: { label: string; value: ReactNode }[] = [
    { label: formatMessage('Vision.algorithm.name'), value: algorithm.name },
    { label: formatMessage('Vision.algorithm.algorithmId'), value: algorithm.algorithmCode },
    { label: formatMessage('Vision.algorithm.source'), value: <SourceTag source={algorithm.source} /> },
    { label: formatMessage('Vision.algorithm.modelStatus'), value: <ModelStatusTag status={algorithm.modelStatus} /> },
    {
      label: formatMessage('Vision.algorithm.modelFile'),
      value: algorithm.modelFile || <Text type="secondary">{formatMessage('Vision.algorithm.modelNotUploaded')}</Text>,
    },
    { label: formatMessage('Vision.algorithm.modelFormat'), value: algorithm.modelFormat?.toUpperCase() || 'ONNX' },
    {
      label: formatMessage('Vision.algorithm.linkedModel'),
      value: algorithm.model ? (
        <Space size={8}>
          <span>
            {algorithm.model.name}
            {algorithm.model.version ? ` ${algorithm.model.version}` : ''}
          </span>
          <ModelStatusTag status={algorithm.model.status} />
        </Space>
      ) : (
        '-'
      ),
    },
  ];

  return (
    <ProModal
      open={open}
      title={formatMessage('Vision.algorithm.detailTitle')}
      width={640}
      fullScreenable={false}
      onCancel={onClose}
      footer={null}
      destroyOnHidden
    >
      <div className={styles.detailGrid}>
        {items.map((item) => (
          <div key={item.label} className={styles.detailItem}>
            <Text className={styles.detailLabel}>{item.label}</Text>
            <div className={styles.detailValue}>{item.value}</div>
          </div>
        ))}
      </div>
      <div className={styles.detailBlock}>
        <Text className={styles.detailLabel}>{formatMessage('Vision.algorithm.labels')}</Text>
        <div className={styles.detailValue}>{renderTags(algorithm.labels || [], styles.algoCardLabelTag)}</div>
      </div>
      <div className={styles.detailBlock}>
        <Text className={styles.detailLabel}>{formatMessage('Vision.algorithm.defaultOutput')}</Text>
        <div className={styles.detailValue}>{renderTags(algorithm.defaultOutput || [], styles.algoOutputTag)}</div>
      </div>
      <div className={styles.detailBlock}>
        <Text className={styles.detailLabel}>{formatMessage('common.description')}</Text>
        <div className={styles.detailValue}>{algorithm.description || <Text type="secondary">--</Text>}</div>
      </div>
      <div className={styles.detailFooter}>
        {canManage && (
          <Button
            onClick={() => {
              onEdit(algorithm);
              onClose();
            }}
          >
            {formatMessage('common.edit')}
          </Button>
        )}
        <Button type="primary" onClick={onClose}>
          {formatMessage('Vision.algorithm.done')}
        </Button>
      </div>
    </ProModal>
  );
};

export default AlgorithmDetailModal;
