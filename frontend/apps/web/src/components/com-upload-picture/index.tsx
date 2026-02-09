import { type FC, useRef } from 'react';
import { Upload, App } from 'antd';
import { AddLarge } from '@carbon/icons-react';
import useTranslate from '@/hooks/useTranslate.ts';

import type { UploadProps, UploadFile } from 'antd';
import usePropsValue from '@/hooks/usePropsValue.ts';

export interface UploadPictureProps extends Omit<UploadProps, 'onChange'> {
  value?: UploadFile;
  onChange?: (file: UploadFile) => void;
  acceptList?: string[];
  className?: string;
  onActionChange?: UploadProps['onChange'];
}

const ComUploadPicture: FC<UploadPictureProps> = ({
  value,
  onChange,
  maxCount = 1,
  acceptList = ['jpg', 'jpeg', 'png', 'svg'],
  className,
  action = '',
  onActionChange,
  ...restProps
}) => {
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const [fileList, setFileList] = usePropsValue<UploadFile>({
    value,
    onChange,
  });
  const accept = acceptList.map((item) => `.${item}`).join(',');
  const acceptMsg = acceptList.map((item) => `.${item}`).join('、');

  const uploadRef = useRef<any>(null);

  const beforeUpload = (file: any) => {
    const fileType = file.name.split('.').pop();
    if (acceptList?.length === 0 || acceptList.includes(fileType.toLowerCase())) {
      if (!action) {
        const previewUrl = URL.createObjectURL(file);
        const newFile = {
          ...file,
          file,
          url: previewUrl,
          thumbUrl: previewUrl,
          status: 'done',
        };
        setFileList([newFile]);
      }
    } else {
      message.warning(formatMessage('common.imgFormatSupport', { format: acceptMsg }));
      return Upload.LIST_IGNORE; //阻止无效文件挂载到组件本身
    }
    if (!action) {
      return false;
    } //阻止调用Upload上传
  };

  const onRemove = () => {
    setFileList([]);
  };

  return (
    <div className={className}>
      <Upload
        action={action}
        listType="picture-card"
        {...restProps}
        fileList={fileList}
        accept={accept}
        beforeUpload={beforeUpload}
        onRemove={onRemove}
        ref={uploadRef}
        onChange={
          onActionChange
            ? (info) => {
                const { file } = info;
                const fileType = file.name.split('.').pop() || '';
                if (acceptList?.length === 0 || acceptList.includes(fileType.toLowerCase())) {
                  setFileList(info.fileList);
                  onActionChange?.(info);
                }
              }
            : undefined
        }
      >
        {fileList?.length >= maxCount ? null : (
          <button style={{ color: 'inherit', cursor: 'inherit', border: 0, background: 'none' }} type="button">
            <AddLarge />
          </button>
        )}
      </Upload>
      <span style={{ color: '#6F6F6F', marginTop: 4 }}>{formatMessage('common.imageSize', { size: '28*28' })}</span>
    </div>
  );
};
export default ComUploadPicture;
