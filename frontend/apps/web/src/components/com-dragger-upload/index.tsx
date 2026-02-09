import { FolderAdd } from '@carbon/icons-react';
import { App, type UploadProps, type GetRef, Upload, type UploadFile } from 'antd';
import { type CSSProperties, forwardRef } from 'react';
import usePropsValue from '@/hooks/usePropsValue.ts';
import { useTranslate } from '@/hooks';
import './index.scss';

const { Dragger } = Upload;
type DraggerRef = GetRef<typeof Dragger>;

interface ComDraggerUploadProps extends Omit<UploadProps, 'onChange'> {
  acceptList?: string[];
  style?: CSSProperties;
  size?: number;
  className?: string;
  value?: UploadFile[];
  defaultValue?: UploadFile[];
  onChange?: (fileList: UploadFile[]) => void;
  onActionChange?: UploadProps['onChange'];
}

const ProDraggerUpload = forwardRef<DraggerRef, ComDraggerUploadProps>(
  (
    {
      value,
      onChange,
      defaultValue,
      acceptList = [],
      children,
      style,
      className,
      onActionChange,
      action = '',
      size,
      ...restProps
    },
    ref
  ) => {
    const [fileList, setFileList] = usePropsValue<UploadFile[]>({
      value,
      onChange,
      defaultValue,
    });
    const formatMessage = useTranslate();
    const { message } = App.useApp();
    const accept = acceptList.map((item) => `.${item}`).join(',');

    const beforeUpload = (file: any) => {
      if (size && file?.size > size) {
        message.warning(formatMessage('common.theFileSizeMax', { size: '2GB' }));
        return false;
      }
      const fileType = file.name.split('.').pop();
      if (acceptList?.length === 0 || acceptList.includes(fileType.toLowerCase())) {
        if (!action) {
          setFileList([file]);
        }
      } else {
        message.warning(formatMessage('common.theFileFormatType', { fileType: accept }));
        return false;
      }
      if (!action) {
        return false;
      }
    };

    return (
      <div style={style} className={className}>
        <Dragger
          className="com-dragger-upload"
          action={action}
          accept={accept}
          maxCount={1}
          beforeUpload={beforeUpload}
          fileList={fileList}
          onRemove={() => {
            setFileList([]);
          }}
          onChange={
            onActionChange
              ? (info) => {
                  const { file } = info;
                  const fileType = file.name.split('.').pop() || '';
                  if (
                    size &&
                    (file?.size || 0) < size &&
                    (acceptList?.length === 0 || acceptList.includes(fileType.toLowerCase()))
                  ) {
                    setFileList(info.fileList);
                    onActionChange?.(info);
                  }
                }
              : undefined
          }
          {...restProps}
          ref={ref}
        >
          {children ? children : <FolderAdd size={100} style={{ color: '#E0E0E0' }} />}
        </Dragger>
      </div>
    );
  }
);

export default ProDraggerUpload;
