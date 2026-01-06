import { type FC, useState, useImperativeHandle } from 'react';
import { Flex } from 'antd';
import { useTranslate } from '@/hooks';

import type { RefObject, Dispatch, SetStateAction } from 'react';
import ProModal from '@/components/pro-modal';
import { UnsTree } from '@/pages/uns/components/export-modal/uns-tree.tsx';
import styles from './index.module.scss';
import { TreeStoreProvider } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { CodeDom } from '@/pages/uns/components/export-modal/code-dom.tsx';

interface ExportModalRef {
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export interface ExportModalProps {
  exportRef?: RefObject<ExportModalRef>;
}

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
            <Flex gap={16} style={{ height: isFullscreen ? '100%' : 500 }}>
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
          </TreeStoreProvider>
        );
      }}
    </ProModal>
  );
};
export default Module;
