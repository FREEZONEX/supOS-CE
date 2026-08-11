import { Image as ImageIcon } from '@carbon/icons-react';
import { App, Spin, Upload, type UploadProps } from 'antd';
import type { RcFile } from 'antd/es/upload';
import classNames from 'classnames';

import { useTranslate } from '@/hooks';

import styles from './AppMediaFields.module.scss';

const { Dragger } = Upload;
const MAX_IMAGE_SIZE = 20 * 1024 * 1024;
const SUPPORTED_IMAGE_TYPES = new Set(['image/jpeg', 'image/png']);

export type AppMediaKind = 'cover' | 'icon';

export interface AppMediaFieldsProps {
  coverUrl?: string;
  iconUrl?: string;
  disabled?: boolean;
  uploadingKind?: AppMediaKind;
  onSelect: (kind: AppMediaKind, file: RcFile) => Promise<void> | void;
}

const AppMediaFields = ({ coverUrl, iconUrl, disabled = false, uploadingKind, onSelect }: AppMediaFieldsProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();

  const createBeforeUpload =
    (kind: AppMediaKind): UploadProps['beforeUpload'] =>
    async (file) => {
      if (!SUPPORTED_IMAGE_TYPES.has(file.type)) {
        message.warning(formatMessage('common.fileFormatType', { fileType: '.jpg, .png' }));
        return Upload.LIST_IGNORE;
      }
      if (file.size > MAX_IMAGE_SIZE) {
        message.warning(formatMessage('common.fileSizeMax', { size: '20MB' }));
        return Upload.LIST_IGNORE;
      }
      await onSelect(kind, file);
      return false;
    };

  const renderPreview = (kind: AppMediaKind, url?: string) => {
    const isCover = kind === 'cover';
    const loading = uploadingKind === kind;

    return (
      <div className={styles.preview}>
        {url ? <img src={url} alt="" className={isCover ? styles['cover-image'] : styles['icon-image']} /> : null}
        {!url && !loading ? (
          <div className={isCover ? styles['cover-empty'] : styles['icon-empty']}>
            <ImageIcon size={isCover ? 48 : 32} />
            <span className={styles['primary-hint']}>
              {isCover
                ? formatMessage('apps.media.dropHint', {}, 'Click or drag image to this area')
                : formatMessage('apps.media.selectIcon', {}, 'Select Icon')}
            </span>
            {isCover ? (
              <span className={styles['secondary-hint']}>
                {formatMessage('apps.media.formatHint', {}, 'Supported formats: .jpg, .png (Max size: 20MB)')}
              </span>
            ) : null}
          </div>
        ) : null}
        {loading ? (
          <div className={styles['loading-overlay']}>
            <Spin size="small" />
          </div>
        ) : null}
        {url && !loading ? (
          <div className={styles['change-overlay']}>
            <ImageIcon size={32} />
          </div>
        ) : null}
      </div>
    );
  };

  return (
    <div className={styles['media-fields']}>
      <div className={styles['cover-field']}>
        <label className={styles.label}>{formatMessage('apps.coverImage', {}, 'Cover Image')}</label>
        <Dragger
          className={classNames(styles['cover-upload'], {
            [styles['cover-upload-filled']]: Boolean(coverUrl),
          })}
          accept=".jpg,.jpeg,.png"
          showUploadList={false}
          disabled={disabled}
          beforeUpload={createBeforeUpload('cover')}
        >
          {renderPreview('cover', coverUrl)}
        </Dragger>
      </div>
      <div className={styles['icon-field']}>
        <label className={styles.label}>{formatMessage('apps.icon')}</label>
        <Upload
          className={styles['icon-upload']}
          accept=".jpg,.jpeg,.png"
          showUploadList={false}
          disabled={disabled}
          beforeUpload={createBeforeUpload('icon')}
        >
          {renderPreview('icon', iconUrl)}
        </Upload>
      </div>
    </div>
  );
};

export default AppMediaFields;
