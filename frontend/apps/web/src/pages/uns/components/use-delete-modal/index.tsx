import { useCallback, useState } from 'react';
import { Button, Form, App, Divider, Flex } from 'antd';
import { WarningFilled } from '@carbon/icons-react';
import { deleteTreeNode, detectIfRemoveApi } from '@/apis/inter-api/uns';
import { useTranslate } from '@/hooks';

import type { UnsTreeNode, InitTreeDataFnType } from '@/pages/uns/types';
import ComCheckbox from '@/components/com-checkbox';
import ProModal from '@/components/pro-modal';
import { useBaseStore } from '@/stores/base';
import { ROOT_NODE_ID } from '@/pages/uns/store/treeStore.tsx';

export interface DeleteModalProps {
  successCallBack: InitTreeDataFnType;
  currentNode?: UnsTreeNode;
  setSelectedNode?: any;
  lazyTree?: boolean;
}

const Module = ({ successCallBack, currentNode, lazyTree }: DeleteModalProps) => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const [open, setOpen] = useState(false);
  const [deleteDetail, setDeleteDetail] = useState<UnsTreeNode | null>();
  const [loading, setLoading] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const { message } = App.useApp();
  const dashboardType = useBaseStore((state) => state.dashboardType);
  const setModalOpen = useCallback(
    (detail: UnsTreeNode) => {
      setDeleteDetail(detail);
      setOpen(true);
      setLoading(true);
      // 先调校验接口
      detectIfRemoveApi({ id: detail?.id })
        .then((data) => {
          if (data?.refs > 0) {
            setShowConfirm(true);
          } else {
            setShowConfirm(false);
          }
        })
        .finally(() => {
          setLoading(false);
        });
    },
    [setOpen, setDeleteDetail]
  );

  const close = () => {
    setOpen(false);
    form.resetFields();
    setLoading(false);
    setShowConfirm(false);
  };
  const cancel = () => {
    close();
  };

  const deleteRequest = async (params: any) => {
    setLoading(true);
    deleteTreeNode(params)
      .then(() => {
        message.success(formatMessage('common.deleteSuccessfully'));
        const clearSelect =
          currentNode?.path?.startsWith(deleteDetail?.path || '') ||
          currentNode?.id === deleteDetail?.id ||
          showConfirm;
        const sourceId = deleteDetail?.parentId;
        const config = lazyTree
          ? {
              queryType: deleteDetail?.pathType === 0 ? 'deleteFolder' : 'deleteFile',
              clearSelect,
              key: sourceId ? sourceId : ROOT_NODE_ID,
              newNodeKey: deleteDetail?.preId || deleteDetail?.nextId,
              nodeDetail: deleteDetail,
              reset: !sourceId,
            }
          : { clearSelect };
        successCallBack(config);
        close();
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const confirm = async (cascade?: boolean) => {
    const { id, pathType } = deleteDetail || {};
    if (pathType === 0) {
      deleteRequest({ cascade, id });
    } else {
      const formData = await form.validateFields();
      const params = { id, ...formData };
      if (cascade) params.cascade = cascade;
      deleteRequest(params);
    }
  };

  const Dom = (
    <ProModal
      aria-label=""
      className="importModalWrap"
      open={open}
      onCancel={close}
      maskClosable={false}
      title={formatMessage(deleteDetail?.pathType === 2 ? 'uns.deleteFile' : 'uns.deleteFolder')}
      width={460}
    >
      <Form name="deleteForm" form={form} colon={false}>
        {deleteDetail?.pathType === 2 ? (
          <>
            {showConfirm && (
              <>
                <div style={{ display: 'flex', gap: '5px', fontSize: '16px' }}>
                  <WarningFilled style={{ color: '#faad14', flexShrink: 0, height: '25px' }} />
                  <span>{formatMessage('uns.thisNodeHasAssociatedComputingNodes')}</span>
                </div>
                <Divider variant="dashed" style={{ margin: '8px 0', borderColor: '#c6c6c6' }} />
              </>
            )}
            <Form.Item
              name="withFlow"
              label=""
              valuePropName="checked"
              initialValue={false}
              style={{ marginBottom: 0 }}
            >
              <ComCheckbox label={formatMessage('uns.deleteAutoFlow')} />
            </Form.Item>
            {dashboardType?.includes('grafana') && (
              <Form.Item
                name="withDashboard"
                label=""
                valuePropName="checked"
                initialValue={true}
                style={{ marginBottom: 0 }}
              >
                <ComCheckbox label={formatMessage('uns.deleteAutoDashboard')} />
              </Form.Item>
            )}
          </>
        ) : (
          <div style={{ display: 'flex', gap: '5px', fontSize: '16px' }}>
            <WarningFilled style={{ color: '#faad14', flexShrink: 0, height: '25px' }} />
            <Flex vertical justify="flex-start">
              <span>
                {showConfirm ? '1、' : ''}
                {formatMessage('uns.deleteFolderTip')}
              </span>
              {showConfirm && <span>2、{formatMessage('uns.thisFolderHasAssociatedComputingNodes')}</span>}
            </Flex>
          </div>
        )}
      </Form>
      <div style={{ marginTop: '20px' }}>
        <Button
          color="default"
          variant="filled"
          onClick={cancel}
          style={{ width: '48%', marginRight: '4%' }}
          size="large"
          disabled={loading}
          title={formatMessage('common.cancel')}
        >
          {formatMessage('common.cancel')}
        </Button>

        <Button
          color="primary"
          variant="solid"
          onClick={() => {
            confirm(true);
          }}
          style={{ width: '48%' }}
          size="large"
          loading={loading}
          title={formatMessage('common.confirm')}
        >
          {formatMessage('common.confirm')}
        </Button>
      </div>
    </ProModal>
  );
  return {
    DeleteModal: Dom,
    setDeleteOpen: setModalOpen,
  };
};
export default Module;
