import { useState } from 'react';
import { App, Button, Space, Upload } from 'antd';
import { Document, Upload as UploadIcon, Close } from '@/components/lucide-icon/carbon';
import { importAlgorithmPackage } from '@/apis/core-api/algorithm';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

type AlgorithmImportModalProps = {
  open: boolean;
  onCancel: () => void;
  onSaved: () => void;
};

const formatSize = (bytes: number) => {
  if (bytes >= 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${bytes} B`;
};

const AlgorithmImportModal = ({ open, onCancel, onSaved }: AlgorithmImportModalProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [file, setFile] = useState<File | null>(null);
  const [saving, setSaving] = useState(false);

  const submit = async () => {
    if (!file) return;
    setSaving(true);
    try {
      await importAlgorithmPackage(file);
      message.success(formatMessage('common.optsuccess'));
      onSaved();
    } finally {
      setSaving(false);
    }
  };

  return (
    <ProModal
      open={open}
      title={formatMessage('Vision.algorithm.addTitle')}
      width={520}
      fullScreenable={false}
      onCancel={onCancel}
      footer={null}
      destroyOnHidden
    >
      {file ? (
        <div className={styles.uploadPicked}>
          <Document size={20} />
          <span className={styles.uploadPickedName}>{file.name}</span>
          <span className={styles.uploadPickedSize}>{formatSize(file.size)}</span>
          <button
            type="button"
            className={styles.uploadPickedRemove}
            aria-label={formatMessage('common.delete')}
            onClick={() => setFile(null)}
          >
            <Close size={16} />
          </button>
          <Upload accept=".zip" showUploadList={false} beforeUpload={(f) => (setFile(f), false)}>
            <Button size="small" icon={<UploadIcon size={14} />}>
              {formatMessage('Vision.algorithm.uploadAgain')}
            </Button>
          </Upload>
        </div>
      ) : (
        <Upload.Dragger accept=".zip" showUploadList={false} beforeUpload={(f) => (setFile(f), false)}>
          <p className={styles.uploadDragIcon}>
            <UploadIcon size={32} />
          </p>
          <p className={styles.uploadDragText}>{formatMessage('Vision.algorithm.uploadModelHint')}</p>
          <p className={styles.uploadDragTip}>{formatMessage('Vision.algorithm.uploadZipTip')}</p>
        </Upload.Dragger>
      )}
      <div className={styles.detailFooter}>
        <Space>
          <Button onClick={onCancel} disabled={saving}>
            {formatMessage('common.cancel')}
          </Button>
          <Button type="primary" loading={saving} disabled={!file} onClick={() => void submit()}>
            {formatMessage('common.add', {}, 'Add')}
          </Button>
        </Space>
      </div>
    </ProModal>
  );
};

export default AlgorithmImportModal;
