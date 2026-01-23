import type { Dispatch, RefObject, SetStateAction } from 'react';
import { type FC, useImperativeHandle, useMemo, useState } from 'react';
import { Divider, Flex, Form, type FormInstance, Segmented } from 'antd';
import { useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';
import { UnsTree } from '@/pages/uns/components/export-modal/uns-tree.tsx';
import styles from './index.module.scss';
import { TreeStoreProvider, useTreeStore } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { CodeDom } from '@/pages/uns/components/export-modal/code-dom.tsx';
import ComButton from '../../../../components/com-button';
import OtherDom from '@/pages/uns/components/export-modal/other-dom.tsx';
import ComStatusDot from '@/components/com-status-dot/ComStatusDot.tsx';

interface ExportModalRef {
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export interface ExportModalProps {
  exportRef?: RefObject<ExportModalRef>;
}

const Tab = ({ form }: { form: FormInstance }) => {
  const { tabType, setTabType, checkedKeys, allChecked } = useTreeStore((state) => ({
    tabType: state.tabType,
    setTabType: state.setTabType,
    checkedKeys: state.checkedKeys,
    allChecked: state.allChecked,
  }));

  const isSelect = Form.useWatch((values) => {
    return (
      values?.dashboardExportParam?.length ||
      values?.eventFlowExportParam?.length ||
      values?.sourceFlowExportParam?.length
    );
  }, form);

  const unsStatus = useMemo(() => {
    const isSelect = allChecked || checkedKeys?.length > 0;
    if (tabType === 'uns') {
      if (isSelect) {
        return 'breathing';
      } else {
        return 'active';
      }
    } else {
      if (isSelect) {
        return 'unactive';
      } else {
        return 'stop';
      }
    }
  }, [tabType, checkedKeys, allChecked]);

  const othersStatus = useMemo(() => {
    if (tabType === 'others') {
      if (isSelect) {
        return 'breathing';
      } else {
        return 'active';
      }
    } else {
      if (isSelect) {
        return 'unactive';
      } else {
        return 'stop';
      }
    }
  }, [tabType, isSelect]);
  return (
    <Segmented<string>
      defaultValue={'uns'}
      options={[
        {
          value: 'uns',
          label: (
            <Flex gap={4} align="center">
              <ComStatusDot status={unsStatus} />
              <span>UNS</span>
            </Flex>
          ),
        },
        {
          value: 'others',
          label: (
            <Flex gap={4} align="center">
              <ComStatusDot status={othersStatus} />
              <span>Others</span>
            </Flex>
          ),
        },
      ]}
      value={tabType}
      onChange={(value: any) => {
        setTabType(value);
      }}
    />
  );
};

const Content = ({ isFullscreen, open }: { isFullscreen?: boolean; open: boolean }) => {
  const formatMessage = useTranslate();
  const { tabType, checkedKeys, allChecked } = useTreeStore((state) => ({
    tabType: state.tabType,
    setTabType: state.setTabType,
    checkedKeys: state.checkedKeys,
    allChecked: state.allChecked,
  }));
  const [form] = Form.useForm();

  const onExport = async () => {
    const values = await form.validateFields();
    console.log(values, checkedKeys, allChecked);
  };
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
        <Tab form={form} />
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
          <ComButton type="primary" onClick={onExport}>
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
