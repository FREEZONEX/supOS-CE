import { type FC, useImperativeHandle, useState } from 'react';
import { Button, Flex, Segmented, Tooltip } from 'antd';
import { useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';
import './index.scss';
import ImportDom from '@/pages/uns/components/import-modal/ImportDom.tsx';
import JsonDom from '@/pages/uns/components/import-modal/JsonDom.tsx';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { Download } from '@carbon/icons-react';

export interface ImportModalProps {
  initTreeData: any;
  importRef: any;
}

const Module: FC<ImportModalProps> = (props) => {
  const { importRef, initTreeData } = props;
  const [open, setOpen] = useState(false);
  const formatMessage = useTranslate();
  const [type, setType] = useState('document');

  useImperativeHandle(importRef, () => ({
    setOpen: setOpen,
  }));

  const onCloseModal = () => {
    setOpen(false);
    setType('document');
  };

  if (!open) return null;
  return (
    <ProModal
      className="importModalWrap"
      open={open}
      onCancel={onCloseModal}
      title={
        <Flex
          style={{
            cursor: 'pointer',
            width: 'fit-content',
          }}
          gap={8}
        >
          <span>{formatMessage('common.import')}</span>
          <Tooltip
            title={
              <div>
                {formatMessage('uns.supportScope')}
                <div>
                  •{' '}
                  {formatMessage('uns.overallProject', {
                    text: `UNS、${formatMessage('home.sourceFlow')}、${formatMessage('home.eventFlow')}、${formatMessage('home.dashboard')}`,
                  })}
                </div>
                <div>• {formatMessage('uns.singleFile', { text: `UNS JSON${formatMessage('common.file')}` })}</div>
              </div>
            }
          >
            <QuestionCircleOutlined />
          </Tooltip>
        </Flex>
      }
      width={460}
      maskClosable={false}
      keyboard={false}
      destroyOnHidden
    >
      {(isFullscreen) => {
        return (
          <Flex vertical style={{ height: isFullscreen ? '100%' : 500 }}>
            <Flex style={{ flexShrink: 0, paddingBottom: 16 }} align="center" justify="space-between">
              <Segmented<string>
                defaultValue={'document'}
                options={[
                  {
                    value: 'document',
                    label: formatMessage('common.all'),
                  },
                  {
                    value: 'json',
                    label: 'UNS',
                  },
                ]}
                value={type}
                onChange={(value) => {
                  setType(value);
                }}
              />
              {type === 'json' && (
                <Button
                  size="small"
                  onClick={(e) => {
                    e.stopPropagation();
                    window.open(`/inter-api/supos/uns/importExport/template/download?fileType=json`, '_self');
                  }}
                >
                  <Download />
                  {formatMessage('common.downloadTemplate')}
                </Button>
              )}
            </Flex>
            <div style={{ flex: 1, overflow: 'hidden' }}>
              <div style={{ display: type === 'document' ? 'inherit' : 'none', height: '100%' }}>
                <ImportDom initTreeData={initTreeData} onCloseModal={onCloseModal} />
              </div>
              <div style={{ display: type === 'json' ? 'inherit' : 'none', height: '100%' }}>
                <JsonDom initTreeData={initTreeData} onCloseModal={onCloseModal} />
              </div>
            </div>
          </Flex>
        );
      }}
    </ProModal>
  );
};
export default Module;
