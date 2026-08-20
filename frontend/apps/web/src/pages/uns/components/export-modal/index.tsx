import type { Dispatch, RefObject, SetStateAction } from 'react';
import { type FC, useImperativeHandle, useState } from 'react';
import { App, Divider, Flex, Form, Segmented } from 'antd';
import { useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';
import { UnsTree } from '@/pages/uns/components/export-modal/uns-tree.tsx';
import styles from './index.module.scss';
import { TreeStoreProvider, useTreeStore } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { CodeDom } from '@/pages/uns/components/export-modal/code-dom.tsx';
import ComButton from '../../../../components/com-button';
import OtherDom from '@/pages/uns/components/export-modal/other-dom.tsx';
import { processedCheckedKeys } from '@/pages/uns/store/utils.ts';
import { getParamsForArray, SelectAllId } from '@/utils';
import { exportExcelGlobal } from '@/apis/core-api';
import { useDownloadNotification } from '@/hooks/useDownloadNotification';

interface ExportModalRef {
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export interface ExportModalProps {
  exportRef?: RefObject<ExportModalRef>;
}

const Tab = ({ onChange }: { onChange: (value: 'uns' | 'others') => void }) => {
  const formatMessage = useTranslate();
  const { tabType } = useTreeStore((state) => ({
    tabType: state.tabType,
  }));

  return (
    <Segmented<string>
      className={styles.exportSegmented}
      defaultValue={'uns'}
      options={[
        {
          value: 'uns',
          label: formatMessage('uns.importUnsTab'),
        },
        {
          value: 'others',
          label: formatMessage('common.others'),
        },
      ]}
      value={tabType}
      onChange={(value) => {
        onChange(value as 'uns' | 'others');
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
  const { tabType, checkedKeys, allChecked, treeData, setTabType, setCheckedKeys, setAllChecked, setJsonData } =
    useTreeStore((state) => ({
      tabType: state.tabType,
      setTabType: state.setTabType,
      checkedKeys: state.checkedKeys,
      allChecked: state.allChecked,
      treeData: state.treeData,
      setCheckedKeys: state.setCheckedKeys,
      setAllChecked: state.setAllChecked,
      setJsonData: state.setJsonData,
    }));
  const [form] = Form.useForm();
  const [showDownloadNotification, contextHolder] = useDownloadNotification();
  const { message } = App.useApp();

  const handleTabChange = (value: 'uns' | 'others') => {
    if (value === 'uns') {
      form.setFieldsValue({
        sourceFlowExportParam: [],
        eventFlowExportParam: [],
      });
    } else {
      setCheckedKeys([]);
      setAllChecked(false);
      setJsonData(undefined);
    }
    setTabType(value);
  };

  const onExport = async () => {
    const values = await form.validateFields();
    const params: any = {};
    if (
      !(allChecked || checkedKeys?.length) &&
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
    if (values.sourceFlowExportParam?.length > 0) {
      params['sourceFlowExportParam'] = getParams(values.sourceFlowExportParam);
    }
    if (values.eventFlowExportParam?.length > 0) {
      params['eventFlowExportParam'] = getParams(values.eventFlowExportParam);
    }
    return exportExcelGlobal(params).then((zip) => {
      showDownloadNotification({ data: zip, name: 'global-export.zip' });
      message.success(formatMessage('common.optsuccess'));
      onClose();
    });
  };
  const isUnsTab = tabType === 'uns';
  const rootStyle = isFullscreen ? { height: '100%' } : isUnsTab ? { height: 600 } : undefined;
  const unsPanelStyle =
    isFullscreen || isUnsTab ? { flex: 1, minHeight: 0, display: isUnsTab ? ('inherit' as const) : ('none' as const), overflow: 'hidden' as const } : { display: 'none' as const };
  const othersPanelStyle =
    tabType === 'others'
      ? { flex: isFullscreen ? 1 : undefined, minHeight: isFullscreen ? 0 : undefined, overflow: isFullscreen ? ('hidden' as const) : undefined }
      : { display: 'none' as const };

  return (
    <Flex style={rootStyle} vertical>
      {contextHolder}
      <div
        style={{
          marginBottom: 24,
          flexShrink: 0,
        }}
      >
        <Tab onChange={handleTabChange} />
      </div>
      {/*uns*/}
      <Flex gap={16} style={unsPanelStyle}>
        <Flex vertical style={{ flex: 1, height: '100%', overflow: 'hidden' }}>
          <Flex className={styles['export-label']}>
            <span>{formatMessage('uns.importUnsTab')}</span>
          </Flex>
          <UnsTree open={open} />
        </Flex>
        <Flex vertical style={{ flex: 1, height: '100%', overflow: 'hidden' }}>
          <Flex className={styles['export-label']}>
            <span>{formatMessage('uns.exportJsonLabel')}</span>
          </Flex>
          <CodeDom />
        </Flex>
      </Flex>

      <div style={othersPanelStyle}>
        <OtherDom form={form} />
      </div>

      <div style={{ flexShrink: 0 }}>
        <Divider className={styles.exportFooterDivider} />
        <Flex align="center" gap={8} justify="flex-end">
          <ComButton color="default" variant="filled" onClick={onClose}>
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
            <Content open={open} isFullscreen={isFullscreen} onClose={close} />
          </TreeStoreProvider>
        );
      }}
    </ProModal>
  );
};
export default Module;
