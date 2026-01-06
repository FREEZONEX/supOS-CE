import { forwardRef, useEffect, useImperativeHandle, useRef, useState } from 'react';
import ProModal from '@/components/pro-modal';
import { Tabs } from 'antd';
import { useTranslate } from '@/hooks';
import ComCodeSnippet from '@/components/com-code-snippet';
import { getFileSchema, getFolderSchema, getLabelSchema, getTemplateSchema } from '@/apis/inter-api/uns';

export interface SchemaModalRef {
  onOpen: (props: any) => void;
  onClose: () => void;
}

export interface SchemaModalProps {
  [key: string]: any;
}

const SchemaModal = forwardRef<SchemaModalRef, SchemaModalProps>((_props, ref) => {
  const [visible, setVisible] = useState(false);
  const formatMessage = useTranslate();
  const modalPropsRef = useRef({});
  const [fileJson, setFileJson] = useState<any>('');
  const [folderJson, setFolderJson] = useState<any>('');
  const [labelJson, setLabelJson] = useState<any>('');
  const [templateJson, setTemplateJson] = useState<any>('');
  const onOpen = (props: any) => {
    modalPropsRef.current = props;
    setVisible(true);
  };
  const onClose = () => {
    setVisible(false);
  };
  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  useEffect(() => {
    if (visible) {
      getFileSchema().then((data) => {
        setFileJson(data);
      });
      getFolderSchema().then((data) => {
        setFolderJson(data);
      });
      getTemplateSchema().then((data) => {
        setTemplateJson(data);
      });
      getLabelSchema().then((data) => {
        setLabelJson(data);
      });
    }
  }, [visible]);

  return (
    <ProModal
      open={visible}
      onCancel={onClose}
      title={'Schema'}
      width={500}
      styles={{
        body: {
          paddingBlockStart: 0,
        },
      }}
    >
      {() => {
        return (
          <Tabs
            items={[
              {
                key: 'folder',
                label: formatMessage('uns.model'),
                children: (
                  <ComCodeSnippet
                    style={{ borderRadius: 3, border: '1px solid #c6c6c6' }}
                    minCollapsedNumberOfRows={24}
                    maxCollapsedNumberOfRows={24}
                    maxExpandedNumberOfRows={36}
                    minExpandedNumberOfRows={36}
                  >
                    <>{JSON.stringify(folderJson, null, 2)}</>
                  </ComCodeSnippet>
                ),
              },
              {
                key: 'file',
                label: formatMessage('uns.instance'),
                children: (
                  <ComCodeSnippet
                    style={{ borderRadius: 3, border: '1px solid #c6c6c6' }}
                    minCollapsedNumberOfRows={24}
                    maxCollapsedNumberOfRows={24}
                    maxExpandedNumberOfRows={36}
                    minExpandedNumberOfRows={36}
                  >
                    <>{JSON.stringify(fileJson, null, 2)}</>
                  </ComCodeSnippet>
                ),
              },
              {
                key: 'template',
                label: formatMessage('common.template'),
                children: (
                  <ComCodeSnippet
                    style={{ borderRadius: 3, border: '1px solid #c6c6c6' }}
                    minCollapsedNumberOfRows={24}
                    maxCollapsedNumberOfRows={24}
                    maxExpandedNumberOfRows={36}
                    minExpandedNumberOfRows={36}
                  >
                    <>{JSON.stringify(templateJson, null, 2)}</>
                  </ComCodeSnippet>
                ),
              },
              {
                key: 'label',
                label: formatMessage('common.label'),
                children: (
                  <ComCodeSnippet
                    style={{ borderRadius: 3, border: '1px solid #c6c6c6' }}
                    minCollapsedNumberOfRows={24}
                    maxCollapsedNumberOfRows={24}
                    maxExpandedNumberOfRows={36}
                    minExpandedNumberOfRows={36}
                  >
                    <>{JSON.stringify(labelJson, null, 2)}</>
                  </ComCodeSnippet>
                ),
              },
            ]}
          />
        );
      }}
    </ProModal>
  );
});

export default SchemaModal;
