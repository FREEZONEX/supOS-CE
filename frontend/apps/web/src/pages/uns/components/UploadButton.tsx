import { Button, Flex, message, Upload as AntUpload } from 'antd';
import { useTranslate } from '@/hooks';
import { FolderAdd, Upload } from '@carbon/icons-react';
import { useState } from 'react';
import { uploadAttachment } from '@/apis/inter-api/attachments.ts';
import { AuthButton } from '@/components/auth';
import ProModal from '@/components/pro-modal';
const { Dragger } = AntUpload;

const UploadButton = ({
  alias,
  documentListRef,
  auth,
  setActiveList,
}: {
  auth?: string;
  alias: string;
  documentListRef: any;
  setActiveList?: any;
}) => {
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const [fileList, setFileList] = useState<any[]>([]);
  const [show, setShow] = useState(false);
  const onClose = () => {
    setFileList([]);
    setShow(false);
  };
  const beforeUpload = (file: any) => {
    if (file.size <= 1024 * 1024 * 10) {
      setFileList((pre) => {
        return [...pre, file];
      });
    } else {
      message.warning(formatMessage('uns.importDocumentMax'));
    }
    return false;
  };

  const onSave = () => {
    setLoading(true);
    uploadAttachment(
      fileList?.map((item: any) => ({ value: item, name: 'files', fileName: item.name })),
      { alias }
    )
      .then(() => {
        documentListRef?.current?.refresh?.();
        message.success(formatMessage('common.optsuccess'));
        onClose();
        setActiveList?.((pre: string[]) => {
          return [...new Set([...(pre || []), 'document'])];
        });
      })
      .finally(() => {
        setLoading(false);
      });
  };
  return (
    <>
      <AuthButton
        auth={auth}
        onClick={() => setShow(true)}
        style={{ border: '1px solid #C6C6C6', background: 'var(--supos-uns-button-color)' }}
        icon={<Upload />}
      >
        {formatMessage('common.upload')}
      </AuthButton>
      <ProModal
        aria-label=""
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>{formatMessage('uns.importDocument')}</span>
          </div>
        }
        onCancel={onClose}
        open={show}
        className="importModalWrap"
        size="xxs"
      >
        <Dragger
          className="uploadWrap"
          action=""
          multiple
          fileList={fileList}
          beforeUpload={beforeUpload}
          onRemove={(file) => {
            setFileList(fileList?.filter((item) => item.uid !== file.uid));
          }}
        >
          <Flex vertical align="center" gap={10}>
            <FolderAdd size={100} style={{ color: '#E0E0E0' }} />
            <span style={{ fontSize: 12 }}>{formatMessage('uns.importDocumentMax')}</span>
          </Flex>
        </Dragger>
        <Button loading={loading} color="primary" variant="solid" block onClick={onSave} style={{ marginTop: 20 }}>
          {formatMessage('common.save')}
        </Button>
      </ProModal>
    </>
  );
};

export default UploadButton;
