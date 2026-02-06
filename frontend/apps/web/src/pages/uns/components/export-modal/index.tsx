import type { Dispatch, RefObject, SetStateAction } from 'react';
import { type FC, useImperativeHandle, useMemo, useState } from 'react';
import { App, Divider, Flex, Form, type FormInstance, Segmented } from 'antd';
import { useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';
import { UnsTree } from '@/pages/uns/components/export-modal/uns-tree.tsx';
import styles from './index.module.scss';
import { TreeStoreProvider, useTreeStore } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { CodeDom } from '@/pages/uns/components/export-modal/code-dom.tsx';
import ComButton from '../../../../components/com-button';
import OtherDom from '@/pages/uns/components/export-modal/other-dom.tsx';
import ComStatusDot from '@/components/com-status-dot/ComStatusDot.tsx';
import { processedCheckedKeys } from '@/pages/uns/store/utils.ts';
import { getParamsForArray, SelectAllId } from '@/utils';
import { exportExcelGlobal } from '@/apis/inter-api';
import { useDownloadNotification } from '@/hooks/useDownloadNotification';

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

const getParams = (list: any) => {
  if (list.some((s: any) => s.value === SelectAllId)) {
    return {};
  }
  return list.reduce(
    (result: any, item: any) => {
      if (item.category === 'group') {
        result.gids.push(item.id);
      } else if (item.category === 'file') {
        result.ids.push(item.id);
      }
      return result;
    },
    { gids: [] as number[], ids: [] as number[] }
  );
};

const Content = ({ isFullscreen, open, onClose }: { isFullscreen?: boolean; open: boolean; onClose: () => void }) => {
  const formatMessage = useTranslate();
  const { tabType, checkedKeys, allChecked, treeData } = useTreeStore((state) => ({
    tabType: state.tabType,
    setTabType: state.setTabType,
    checkedKeys: state.checkedKeys,
    allChecked: state.allChecked,
    treeData: state.treeData,
  }));
  const [form] = Form.useForm();
  const [showDownloadNotification, contextHolder] = useDownloadNotification();
  const { message } = App.useApp();

  const onExport = async () => {
    const values = await form.validateFields();
    const params: any = {};
    if (
      !(allChecked || checkedKeys?.length) &&
      !values?.dashboardExportParam?.length &&
      !values?.eventFlowExportParam?.length &&
      !values?.sourceFlowExportParam?.length
    ) {
      message.warning(formatMessage('home.mustOne'));
      return;
    }
    if (allChecked) {
      params['unsExportParam'] = {
        exportType: 'ALL',
      };
    } else {
      // 根据checkedKeys匹配节点信息
      const matchedNodes = processedCheckedKeys({
        checkedKeys,
        strategy: 'SHOW_PARENT',
        treeData, // 添加treeData参数
      });
      params['unsExportParam'] = {
        ...getParamsForArray(matchedNodes as any[], 'pathType', {
          groups: {
            0: 'folders',
            2: 'files',
          },
          extract: 'id',
        }),
        checkSmallFile: false,
      };
    }
    if (values.dashboardExportParam?.length > 0) {
      params['dashboardExportParam'] = getParams(values.dashboardExportParam);
    }
    if (values.sourceFlowExportParam?.length > 0) {
      params['sourceFlowExportParam'] = getParams(values.sourceFlowExportParam);
    }
    if (values.eventFlowExportParam?.length > 0) {
      params['eventFlowExportParam'] = getParams(values.eventFlowExportParam);
    }
    return exportExcelGlobal(params).then((zip) => {
      showDownloadNotification({ data: zip, name: 'global-export.zip' });
    });
  };
  return (
    <Flex
      style={{
        height: isFullscreen ? '100%' : 600,
      }}
      vertical
    >
      {contextHolder}
      <div
        style={{
          marginBottom: 24,
          flexShrink: 0,
        }}
      >
        <Tab form={form} />
      </div>
      {/*uns*/}
      <Flex gap={16} style={{ flex: 1, display: tabType === 'uns' ? 'inherit' : 'none', overflow: 'hidden' }}>
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
          <ComButton onClick={onClose}>{formatMessage('common.cancel')}</ComButton>
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
            <Content open={open} isFullscreen={isFullscreen} onClose={close} />
          </TreeStoreProvider>
        );
      }}
    </ProModal>
  );
};
export default Module;
