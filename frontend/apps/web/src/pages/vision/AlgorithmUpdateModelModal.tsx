import { useState } from 'react';
import { Alert, App, Button, Space, Upload } from 'antd';
import { Document, Upload as UploadIcon, Close } from '@/components/lucide-icon/carbon';
import { uploadAlgorithmModel, type VisionAlgorithm } from '@/apis/core-api/algorithm';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

type AlgorithmUpdateModelModalProps = {
  open: boolean;
  algorithm: VisionAlgorithm | null;
  onClose: () => void;
  onSaved: () => void;
};

const formatSize = (bytes: number) => {
  if (bytes >= 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${bytes} B`;
};

const AlgorithmUpdateModelModal = ({ open, algorithm, onClose, onSaved }: AlgorithmUpdateModelModalProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [failed, setFailed] = useState(false);

  const submit = async () => {
    if (!file || !algorithm) return;
    setUploading(true);
    setFailed(false);
    try {
      await uploadAlgorithmModel(algorithm.id, file);
      message.success(formatMessage('common.optsuccess'));
      onSaved();
    } catch {
      setFailed(true);
    } finally {
      setUploading(false);
    }
  };

  return (
    <ProModal
      open={open}
      title={formatMessage('Vision.algorithm.updateModelTitle')}
      width={520}
      fullScreenable={false}
      onCancel={onClose}
      footer={null}
      destroyOnHidden
    >
      {failed && (
        <Alert
          className={styles.uploadAlert}
          type="error"
          showIcon
          message={formatMessage('Vision.algorithm.modelUploadFailed')}
        />
      )}
      {file && !failed && (
        <Alert
          className={styles.uploadAlert}
          type="info"
          showIcon
          message={formatMessage('Vision.algorithm.modelRestartHint')}
        />
      )}
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
          <Upload accept=".onnx" showUploadList={false} beforeUpload={(f) => (setFile(f), setFailed(false), false)}>
            <Button size="small" icon={<UploadIcon size={14} />}>
              {formatMessage('Vision.algorithm.uploadAgain')}
            </Button>
          </Upload>
        </div>
      ) : (
        <Upload.Dragger
          accept=".onnx"
          showUploadList={false}
          beforeUpload={(f) => (setFile(f), setFailed(false), false)}
        >
          <p className={styles.uploadDragIcon}>
            <UploadIcon size={32} />
          </p>
          <p className={styles.uploadDragText}>{formatMessage('Vision.algorithm.uploadModelHint')}</p>
          <p className={styles.uploadDragTip}>{formatMessage('Vision.algorithm.uploadOnnxTip')}</p>
        </Upload.Dragger>
      )}
      <div className={styles.detailFooter}>
        <Space>
          <Button onClick={onClose} disabled={uploading}>
            {formatMessage('common.cancel')}
          </Button>
          <Button type="primary" loading={uploading} disabled={!file} onClick={() => void submit()}>
            {formatMessage('Vision.algorithm.updateModel')}
          </Button>
        </Space>
      </div>
    </ProModal>
  );
};

export default AlgorithmUpdateModelModal;
