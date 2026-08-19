import { type FC, type MouseEvent, useRef } from 'react';
import { Upload, App } from 'antd';
import { AddLarge } from '@carbon/icons-react';
import useTranslate from '@/hooks/useTranslate.ts';

import type { UploadProps, UploadFile } from 'antd';
import usePropsValue from '@/hooks/usePropsValue.ts';

export interface UploadPictureProps extends Omit<UploadProps, 'onChange'> {
  value?: UploadFile;
  onChange?: (file: UploadFile) => void;
}

const Module: FC<UploadPictureProps> = ({ value, onChange, maxCount = 1, ...restProps }) => {
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const [fileList, setFileList] = usePropsValue<UploadFile>({
    value,
    onChange,
  });

  const uploadRef = useRef<any>(null);
  const rootRef = useRef<HTMLDivElement>(null);

  const openFileDialog = (event?: MouseEvent<HTMLElement>) => {
    event?.stopPropagation();
    const target = event?.target as HTMLElement | null;
    if (target?.tagName === 'INPUT') return;
    rootRef.current?.querySelector<HTMLInputElement>('input[type="file"]')?.click();
  };

  const beforeUpload = (file: any) => {
    const fileType = file.name.split('.').pop();
    if (['jpg', 'jpeg', 'png', 'svg'].includes(fileType.toLowerCase())) {
      const previewUrl = URL.createObjectURL(file);
      const newFile = {
        ...file,
        file,
        url: previewUrl,
        thumbUrl: previewUrl,
        status: 'done',
      };
      setFileList([newFile]);
    } else {
      message.warning(formatMessage('common.imgFormatSupport', { format: 'jpg、jpeg、png、svg' }));
      return Upload.LIST_IGNORE; //阻止无效文件挂载到组件本身
    }
    return false; //阻止调用Upload上传
  };

  const onRemove = () => {
    setFileList([]);
  };

  return (
    <div ref={rootRef} onClick={openFileDialog}>
      <Upload
        action=""
        listType="picture-card"
        {...restProps}
        fileList={fileList}
        accept=".jpg,.jpeg,.png,.svg"
        beforeUpload={beforeUpload}
        onRemove={onRemove}
        ref={uploadRef}
        openFileDialogOnClick={false}
      >
        {fileList?.length >= maxCount ? null : (
          <button
            style={{ color: 'inherit', cursor: 'pointer', border: 0, background: 'none' }}
            type="button"
            onClick={openFileDialog}
          >
            <AddLarge />
          </button>
        )}
      </Upload>
      <span style={{ color: '#6F6F6F', marginTop: 4 }}>{formatMessage('common.imageSize', { size: '28*28' })}</span>
    </div>
  );
};
export default Module;
