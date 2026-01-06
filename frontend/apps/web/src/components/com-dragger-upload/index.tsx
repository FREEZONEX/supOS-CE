import { FolderAdd } from '@carbon/icons-react';
import { App, type DrawerProps, type GetRef, Upload, type UploadFile } from 'antd';
import { type CSSProperties, forwardRef } from 'react';
import usePropsValue from '@/hooks/usePropsValue.ts';
import { useTranslate } from '@/hooks';
import './index.scss';

const { Dragger } = Upload;
type DraggerRef = GetRef<typeof Dragger>;

interface ComDraggerUploadProps extends Omit<DrawerProps, 'action' | 'fileList' | 'height'> {
  value?: UploadFile[];
  defaultValue?: UploadFile[];
  onChange?: (fileList: UploadFile[]) => void;
  acceptList?: string[];
  style?: CSSProperties;
}

const ComDraggerUpload = forwardRef<DraggerRef, ComDraggerUploadProps>(
  ({ value, onChange, defaultValue, acceptList = [], children, style, ...restProps }, ref) => {
    const [fileList, setFileList] = usePropsValue<UploadFile[]>({
      value,
      onChange,
      defaultValue,
    });
    const formatMessage = useTranslate();
    const { message } = App.useApp();
    const accept = acceptList.map((item) => `.${item}`).join(',');

    const beforeUpload = (file: any) => {
      const fileType = file.name.split('.').pop();
      if (acceptList?.length === 0 || acceptList.includes(fileType.toLowerCase())) {
        setFileList([file]);
      } else {
        message.warning(formatMessage('common.theFileFormatType', { fileType: accept }));
      }
      return false;
    };

    return (
      <div style={style}>
        <Dragger
          className="com-dragger-upload"
          action=""
          accept={accept}
          maxCount={1}
          beforeUpload={beforeUpload}
          fileList={fileList}
          onRemove={() => {
            setFileList([]);
          }}
          {...restProps}
          ref={ref}
        >
          {children ? children : <FolderAdd size={100} style={{ color: '#E0E0E0' }} />}
        </Dragger>
      </div>
    );
  }
);

export default ComDraggerUpload;
