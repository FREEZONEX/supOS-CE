import { type FC, useState, useImperativeHandle, useEffect } from 'react';
import { Button, Form, App, Select, TreeSelect, Modal, Dropdown, Badge, Space } from 'antd';
import {
  downloadUnsFile,
  exportExcel,
  getTreeData,
  getUnsExportRecordsApi,
  getUnsLazyTree,
  unsExportRecordConfirmApi,
} from '@/apis/inter-api/uns';
import { useTranslate } from '@/hooks';

import type { RefObject, Dispatch, SetStateAction } from 'react';
import ProModal from '@/components/pro-modal';
import { ProTreeSelect } from '@/components';
import { CustomAxiosConfigEnum, downloadFn } from '@/utils';
import { Document, Download, Folder, FolderOpen } from '@carbon/icons-react';
import { useBaseStore } from '@/stores/base';
import { getParamsForArray } from '@/utils/uns.ts';
import { DownOutlined } from '@ant-design/icons';

interface ExportModalRef {
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export interface ExportModalProps {
  exportRef?: RefObject<ExportModalRef>;
  setButtonExportRecords?: any;
}
const { SHOW_PARENT } = TreeSelect;

const Module: FC<ExportModalProps> = (props) => {
  const { exportRef, setButtonExportRecords } = props;
  const [form] = Form.useForm();
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();
  const { lazyTree, currentUserInfo } = useBaseStore((state) => ({
    lazyTree: state.systemInfo?.lazyTree,
    currentUserInfo: state?.currentUserInfo,
  }));
  const [exportRecords, setExportRecords] = useState([]);
  const [modal, contextHolder] = Modal.useModal();

  const formatMessage = useTranslate();

  const save = async () => {
    const values = await form.validateFields();
    const { uns, fileType } = values;
    const params = getParamsForArray(uns, 'pathType', {
      groups: {
        0: 'models',
        2: 'instances',
      },
      extract: 'value',
    });
    if (!params?.models?.length && !params?.instances?.length && !params?.exportType) {
      message.warning(formatMessage('uns.pleaseSelectTheInstanceToExport'));
      return;
    }

    setLoading(true);
    return exportExcel({
      fileType,
      ...params,
    })
      .then(() => {
        let secondsToGo = 5;
        const instance = modal.success({
          title: formatMessage('home.exportSuccess'),
          okText: `${formatMessage('common.ok')}(${secondsToGo})`,
        });
        const timer = setInterval(() => {
          secondsToGo -= 1;
          instance.update({ okText: `${formatMessage('common.ok')}(${secondsToGo})` });
        }, 1000);
        setTimeout(() => {
          clearInterval(timer);
          instance.destroy();
        }, 5 * 1000);
        form.resetFields();
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const close = () => {
    setOpen(false);
    if (!loading) {
      form.resetFields();
    }
  };
  useImperativeHandle(exportRef, () => ({
    setOpen: setOpen,
  }));

  useEffect(() => {
    if (open) {
      getUnsExportRecordsApi({
        userId: currentUserInfo?.sub,
      }).then((data: any) => {
        setExportRecords(data);
      });
    }
  }, [open]);

  return (
    <ProModal
      className="exportModalWrap"
      open={open}
      onCancel={close}
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span>{formatMessage('uns.export')}</span>
          <Dropdown
            onOpenChange={(open) => {
              if (open) {
                getUnsExportRecordsApi().then((data: any) => {
                  setExportRecords(data);
                  const ids = data?.filter((f: any) => !f.confirm)?.map((d: any) => d.id);
                  if (ids?.length > 0) {
                    unsExportRecordConfirmApi({
                      ids,
                    }).then(() => {
                      setButtonExportRecords((pre: any) => {
                        return pre.map((i: any) => ({
                          ...i,
                          confirm: true,
                        }));
                      });
                      setExportRecords((pre: any) => {
                        return pre.map((i: any) => ({
                          ...i,
                          confirm: true,
                        }));
                      });
                    });
                  }
                });
              }
            }}
            menu={{
              items:
                exportRecords?.length > 0
                  ? [
                      ...(exportRecords?.map((m: any) => {
                        return {
                          label: m?.fileName,
                          key: m.id,
                          extra: <Download style={{ verticalAlign: 'middle' }} />,
                          onClick: () => {
                            downloadUnsFile({ path: m.filePath }).then((data) => {
                              downloadFn({ data, name: m.fileName });
                            });
                          },
                        };
                      }) || []),
                      {
                        type: 'divider',
                      },
                      {
                        key: '-2',
                        label: formatMessage('home.fiveRecord'),
                        disabled: true,
                      },
                    ]
                  : [
                      {
                        disabled: true,
                        label: formatMessage('home.noExport'),
                        key: '-1',
                      },
                    ],
            }}
          >
            <Badge dot={exportRecords?.some((s: any) => !s.confirm)}>
              <Button color="default" variant="filled" iconPosition="end" style={{ padding: '4px 10px' }}>
                <Space>
                  {formatMessage('common.exported')}
                  <DownOutlined />
                </Space>
              </Button>
            </Badge>
          </Dropdown>
        </div>
      }
      width={500}
      maskClosable={false}
      keyboard={false}
    >
      {contextHolder}
      <Form
        name="exportForm"
        form={form}
        colon={false}
        style={{ color: 'var(--supos-text-color)', maxHeight: '300px', overflowY: 'auto' }}
        labelCol={{ span: 6 }}
        wrapperCol={{ span: 18 }}
        labelAlign="left"
        labelWrap
        initialValues={{ uns: [], all: false }}
        disabled={loading}
      >
        <Form.Item name="fileType" label={formatMessage('uns.fileType')} initialValue={'excel'}>
          <Select
            options={[
              { label: 'EXCEL', value: 'excel' },
              { label: 'JSON', value: 'json' },
            ]}
          />
        </Form.Item>
        <Form.Item label={formatMessage('home.uns')} name="uns">
          <ProTreeSelect
            showSearch={false}
            loadDataEnable
            lazy={lazyTree}
            listHeight={350}
            maxTagCount="responsive"
            showSwitcherIcon
            treeCheckable
            popupMatchSelectWidth={700}
            fieldNames={{ label: 'pathName', value: 'id' }}
            showCheckedStrategy={SHOW_PARENT}
            api={(params, config) =>
              lazyTree
                ? getUnsLazyTree(
                    {
                      parentId: params?.key,
                      keyword: params?.searchValue,
                      pageSize: params!.pageSize!,
                      pageNo: params!.pageNo!,
                      searchType: 1,
                    },
                    {
                      [CustomAxiosConfigEnum.BusinessResponse]: true,
                      ...config,
                    }
                  ).then((info: any) => {
                    return {
                      ...info,
                      data: info.data?.map((item: any) => ({
                        ...item,
                        isLeaf: !item.hasChildren,
                        key: item.id,
                      })),
                    };
                  })
                : getTreeData({ key: params?.searchValue }).then((data: any) => {
                    return {
                      data,
                    };
                  })
            }
            treeNodeIcon={(dataNode: any, _treeExpandedKeys = []) => {
              if (dataNode.pathType === 0) {
                return _treeExpandedKeys.includes(dataNode.key!) ? (
                  <FolderOpen style={{ flexShrink: 0, marginRight: '5px' }} />
                ) : (
                  <Folder style={{ flexShrink: 0, marginRight: '5px' }} />
                );
              } else if (dataNode.pathType === 2) {
                return <Document style={{ flexShrink: 0, marginRight: '5px' }} />;
              }
              return null;
            }}
          />
        </Form.Item>
      </Form>
      <Button
        className="exportConfirm"
        color="primary"
        variant="solid"
        onClick={save}
        block
        style={{ marginTop: '10px' }}
        loading={loading}
        disabled={loading}
      >
        {formatMessage('common.confirm')}
      </Button>
    </ProModal>
  );
};
export default Module;
