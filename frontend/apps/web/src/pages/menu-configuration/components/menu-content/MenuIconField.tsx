import { type FC, useMemo } from 'react';
import { App, Upload } from 'antd';
import useTranslate from '@/hooks/useTranslate.ts';
import { MenuLucideIcon } from '@/components/lucide-icon';
import { Image as ImageIcon } from '@/components/lucide-icon/carbon';
import type { MenuProps } from '../../store/types.ts';
import styles from './BasicInfo.module.scss';

const uploadedAsset = (response: unknown) =>
  (response as any)?.data?.list?.[0] || (response as any)?.list?.[0] || (response as any)?.data || response;

const isUploadedIconPath = (icon?: string) => {
  const value = String(icon || '');
  return (
    value.startsWith('/api/core/assets/') ||
    value.startsWith('/inter-api/supos/attachment/download') ||
    value.startsWith('http://') ||
    value.startsWith('https://')
  );
};

export interface MenuIconFieldProps {
  value?: string;
  onChange?: (icon: string) => void;
  menuItem?: MenuProps | null;
  disabled?: boolean;
}

const MenuIconField: FC<MenuIconFieldProps> = ({ value, onChange, menuItem, disabled }) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();

  const iconItem = useMemo(
    () => ({
      ...(menuItem || {}),
      icon: (value && isUploadedIconPath(value) ? value : value || menuItem?.icon) as string | undefined,
      code: menuItem?.code,
      resourceKey: menuItem?.resourceKey,
      url: menuItem?.url,
    }),
    [menuItem, value]
  );

  const preview = useMemo(() => {
    if (value && isUploadedIconPath(value)) {
      return <img src={value} alt="" className={styles.iconImage} />;
    }
    if (iconItem.icon || iconItem.code || iconItem.resourceKey) {
      return <MenuLucideIcon item={iconItem as any} size={18} className={styles.iconSvg} />;
    }
    return <ImageIcon size={18} className={styles.iconSvg} />;
  }, [iconItem, value]);

  const iconBox = <div className={styles.iconBox}>{preview}</div>;

  return (
    <div className={styles.menuIconField}>
      {disabled ? (
        <div className={`${styles.iconBox} ${styles.iconBoxReadonly}`}>{preview}</div>
      ) : (
        <Upload
          action="/api/core/assets"
          withCredentials
          accept=".jpg,.jpeg,.png,.svg"
          showUploadList={false}
          beforeUpload={(file) => {
            const fileType = file.name.split('.').pop()?.toLowerCase() || '';
            if (!['jpg', 'jpeg', 'png', 'svg'].includes(fileType)) {
              message.warning(formatMessage('common.imgFormatSupport', { format: '.jpg、.jpeg、.png、.svg' }));
              return Upload.LIST_IGNORE;
            }
            return true;
          }}
          onChange={(info) => {
            if (info.file.status === 'done') {
              const asset = uploadedAsset(info.file.response);
              const fileId = asset?.fileId || asset?.id;
              onChange?.(asset?.downloadUrl || (fileId ? `/api/core/assets/${fileId}/download` : ''));
              return;
            }
            if (info.file.status === 'error') {
              message.error(formatMessage('common.serverBusy'));
            }
          }}
        >
          {iconBox}
        </Upload>
      )}
      <span className={styles.iconHint}>{formatMessage('common.imagePixels', { size: '32×32' }, '32×32 pixels')}</span>
    </div>
  );
};

export default MenuIconField;
