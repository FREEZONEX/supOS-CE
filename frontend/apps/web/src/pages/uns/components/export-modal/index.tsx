import { type FC, useState, useImperativeHandle } from 'react';
import { Divider, Flex, Form, Segmented } from 'antd';
import { useTranslate } from '@/hooks';

import type { RefObject, Dispatch, SetStateAction } from 'react';
import ProModal from '@/components/pro-modal';
import { UnsTree } from '@/pages/uns/components/export-modal/uns-tree.tsx';
import styles from './index.module.scss';
import { TreeStoreProvider, useTreeStore } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { CodeDom } from '@/pages/uns/components/export-modal/code-dom.tsx';
import ComButton from '../../../../components/com-button';
import OtherDom from '@/pages/uns/components/export-modal/other-dom.tsx';

interface ExportModalRef {
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export interface ExportModalProps {
  exportRef?: RefObject<ExportModalRef>;
}

const Content = ({ isFullscreen, open }: { isFullscreen?: boolean; open: boolean }) => {
  const formatMessage = useTranslate();
  const { tabType, setTabType } = useTreeStore((state) => ({
    tabType: state.tabType,
    setTabType: state.setTabType,
  }));
  const [form] = Form.useForm();

  return (
    <Flex
      style={{
        height: isFullscreen ? '100%' : 600,
      }}
      vertical
    >
      <div
        style={{
          marginBottom: 24,
          flexShrink: 0,
        }}
      >
        <Segmented<string>
          defaultValue={'uns'}
          options={[
            {
              value: 'uns',
              label: 'UNS',
            },
            {
              value: 'others',
              label: 'Others',
            },
          ]}
          value={tabType}
          onChange={(value: any) => {
            setTabType(value);
          }}
        />
      </div>
      {/*uns*/}
      <Flex gap={16} style={{ flex: 1, display: tabType === 'uns' ? 'inherit' : 'none' }}>
        <Flex vertical style={{ flex: 1, height: '100%', overflow: 'hidden' }}>
          <Flex className={styles['export-label']}>
            <span>UNS</span>
          </Flex>
          <UnsTree open={open} />
        </Flex>
        <Flex vertical style={{ flex: 1, height: '100%', overflow: 'hidden' }}>
          <Flex className={styles['export-label']}>
            <span>JSON</span>
          </Flex>
          <CodeDom />
        </Flex>
      </Flex>

      <div style={{ flex: 1, display: tabType === 'others' ? 'inherit' : 'none' }}>
        <OtherDom form={form} />
      </div>

      <div style={{ flexShrink: 0 }}>
        <Divider style={{ backgroundColor: 'rgb(198, 198, 198)' }} />
        <Flex align="center" gap={8} justify="flex-end">
          <ComButton
            onClick={() => {
              return;
            }}
          >
            {formatMessage('common.cancel')}
          </ComButton>
          <ComButton
            type="primary"
            onClick={() => {
              return;
            }}
          >
            {formatMessage('common.export')}
          </ComButton>
        </Flex>
      </div>
    </Flex>
  );
};

const Module: FC<ExportModalProps> = (props) => {
  const { exportRef } = props;
  const [open, setOpen] = useState(false);
  const formatMessage = useTranslate();

  const close = () => {
    setOpen(false);
  };
  useImperativeHandle(exportRef, () => ({
    setOpen: setOpen,
  }));

  return (
    <ProModal
      className="exportModalWrap"
      open={open}
      onCancel={close}
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span>{formatMessage('uns.export')}</span>
        </div>
      }
      width={750}
      maskClosable={false}
      keyboard={false}
      destroyOnHidden
    >
      {(isFullscreen) => {
        return (
          <TreeStoreProvider>
            <Content open={open} isFullscreen={isFullscreen} />
          </TreeStoreProvider>
        );
      }}
    </ProModal>
  );
};
export default Module;
