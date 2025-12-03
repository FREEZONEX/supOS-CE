import { type FC, useState, useImperativeHandle, useMemo, useEffect } from 'react';
import { Button, Form, App, Dropdown, Input, TreeSelect, Divider, ConfigProvider, Space, Badge, Modal } from 'antd';
import { getDashboardList, getTreeData, getUnsLazyTree } from '@/apis/inter-api/uns';
import { useTranslate } from '@/hooks';

import type { RefObject, Dispatch, SetStateAction } from 'react';
import ProModal from '@/components/pro-modal';
import { Document, Download, Folder, FolderOpen } from '@carbon/icons-react';
import { ProTreeSelect } from '@/components/pro-tree-select';
import { flowPage } from '@/apis/inter-api/flow.ts';
import { flowPage as EventFlowPage } from '@/apis/inter-api/event-flow.ts';
import { downloadFn } from '@/utils/blob.ts';
import { CustomAxiosConfigEnum } from '@/utils/request.ts';
import {
  downloadGlobalFile,
  exportGlobal,
  getGlobalExportRecords,
  globalExportRecordConfirm,
} from '@/apis/inter-api/global.ts';
import { DownOutlined } from '@ant-design/icons';
import { useBaseStore } from '@/stores/base';
import { getParamsForArray } from '@/utils/uns.ts';
interface ExportModalRef {
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export interface ExportModalProps {
  exportRef?: RefObject<ExportModalRef>;
  setButtonExportRecords?: any;
}

const { SHOW_PARENT } = TreeSelect;

// function generateDatedFilename(url: string, dateStr: string): string {
//   const filenameWithExt = url.split('/').pop() || '';
//   const [filename, ext] = filenameWithExt.split('.');
//   return `${filename}_${dateStr}.${ext}`;
// }

const Module: FC<ExportModalProps> = (props) => {
  const { exportRef, setButtonExportRecords } = props;
  const [form] = Form.useForm();
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();
  const [modal, contextHolder] = Modal.useModal();
  const [exportRecords, setExportRecords] = useState([]);

  const formatMessage = useTranslate();
  const { containerList, dashboardType, lazyTree } = useBaseStore((state) => ({
    containerList: state.containerList,
    dashboardType: state.dashboardType,
    lazyTree: state.systemInfo?.lazyTree,
  }));
  const hasNodeRed = !!containerList?.aboutUs?.some((s) => s.name === 'nodered');
  const hasEventflow = !!containerList?.aboutUs?.some((s) => s.name === 'eventflow');
  const getRecords = () => {
    return getGlobalExportRecords().then((data) => {
      return data;
    });
  };

  const save = async () => {
    const values = await form.validateFields();
    const { name, unsExportParam, dashboardExportParam, eventFlowExportParam, sourceFlowExportParam } = values;
    if (
      !unsExportParam?.length &&
      !dashboardExportParam?.length &&
      !eventFlowExportParam?.length &&
      !sourceFlowExportParam?.length
    ) {
      return message.warning(formatMessage('home.mustOne'));
    }
    setLoading(true);
    return exportGlobal({
      name,
      sourceFlowExportParam: getParamsForArray(sourceFlowExportParam),
      eventFlowExportParam: getParamsForArray(eventFlowExportParam),
      dashboardExportParam: getParamsForArray(dashboardExportParam),
      unsExportParam: {
        fileType: 'json',
        ...getParamsForArray(unsExportParam, 'type', {
          groups: {
            0: 'models',
            2: 'instances',
          },
          extract: 'value',
        }),
      },
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
        // close();
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
  const formDom = useMemo(() => {
    if (!open) return null;
    return (
      <ConfigProvider
        theme={{
          components: {
            Form: {
              itemMarginBottom: 12,
            },
          },
        }}
      >
        <Form
          layout="vertical"
          name="exportForm"
          form={form}
          colon={false}
          style={{ color: 'var(--supos-text-color)', maxHeight: '500px', overflowY: 'auto' }}
          initialValues={{
            name: '',
            unsExportParam: [],
            sourceFlowExportParam: [],
            eventFlowExportParam: [],
            dashboardExportParam: [],
          }}
          disabled={loading}
        >
          <Form.Item
            label={formatMessage('home.exportName')}
            name="name"
            rules={[
              {
                required: true,
                message: formatMessage('rule.required'),
              },
            ]}
          >
            <Input showCount maxLength={10} placeholder={formatMessage('uns.pleaseInputName')} />
          </Form.Item>
          <Divider style={{ borderColor: '#c6c6c6', margin: '8px 0' }} />
          <Form.Item label={formatMessage('home.uns')} name="unsExportParam">
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
                if (dataNode.type === 0) {
                  return _treeExpandedKeys.includes(dataNode.key!) ? (
                    <FolderOpen style={{ flexShrink: 0, marginRight: '5px' }} />
                  ) : (
                    <Folder style={{ flexShrink: 0, marginRight: '5px' }} />
                  );
                } else if (dataNode.type === 2) {
                  return <Document style={{ flexShrink: 0, marginRight: '5px' }} />;
                }
                return null;
              }}
            />
          </Form.Item>
          {hasNodeRed && (
            <Form.Item label={formatMessage('home.sourceFlow')} name="sourceFlowExportParam">
              <ProTreeSelect
                lazy
                maxTagCount="responsive"
                treeCheckable
                popupMatchSelectWidth={500}
                api={(params) =>
                  flowPage({
                    k: params?.searchValue,
                    pageNo: params?.pageNo,
                    pageSize: params?.pageSize,
                  }).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        value: item.id,
                        title: item.flowName,
                        key: item.id,
                      })),
                    };
                  })
                }
              />
            </Form.Item>
          )}
          {hasEventflow && (
            <Form.Item label={formatMessage('home.eventFlow')} name="eventFlowExportParam">
              <ProTreeSelect
                allowClear
                maxTagCount="responsive"
                treeCheckable
                showSwitcherIcon={false}
                popupMatchSelectWidth={500}
                lazy={true}
                api={(params) =>
                  EventFlowPage({
                    k: params?.searchValue,
                    pageNo: params?.pageNo,
                    pageSize: params?.pageSize,
                  }).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        value: item.id,
                        title: item.flowName,
                        key: item.id,
                      })),
                    };
                  })
                }
              />
            </Form.Item>
          )}
          {dashboardType?.length && (
            <Form.Item label={formatMessage('home.dashboard')} name="dashboardExportParam">
              <ProTreeSelect
                allowClear
                maxTagCount="responsive"
                treeCheckable
                showSwitcherIcon={false}
                popupMatchSelectWidth={500}
                lazy={true}
                api={(params) =>
                  getDashboardList({
                    k: params?.searchValue,
                    pageNo: params?.pageNo,
                    pageSize: params?.pageSize,
                    type: dashboardType?.length >= 2 ? undefined : dashboardType?.includes('fuxa') ? 2 : 1,
                  }).then((data) => {
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: data?.data?.map((item: any) => ({
                        ...item,
                        value: item.id,
                        title: item.name,
                        key: item.id,
                      })),
                    };
                  })
                }
              />
            </Form.Item>
          )}
        </Form>
      </ConfigProvider>
    );
  }, [open]);

  useEffect(() => {
    if (open) {
      getRecords().then((data: any) => {
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
          <span>{formatMessage('common.export')}</span>
          <Dropdown
            onOpenChange={(open) => {
              if (open) {
                getRecords().then((data: any) => {
                  setExportRecords(data);
                  const ids = data?.filter((f: any) => !f.confirm)?.map((d: any) => d.id);
                  if (ids?.length > 0) {
                    globalExportRecordConfirm({
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
                        // const fileName = generateDatedFilename(
                        //   m.filePath,
                        //   formatTimestamp(m.exportTime, 'YYMMDD', true)
                        // );
                        return {
                          label: m?.fileName,
                          // <Flex justify="flex-start" align="center" gap={4}>
                          //   <div
                          //     style={{
                          //       width: 3,
                          //       height: 3,
                          //       background: !m.confirm ? '#FF832B' : 'transparent',
                          //       borderRadius: '50%',
                          //     }}
                          //   ></div>
                          //   {m?.fileName}
                          // </Flex>
                          key: m.id,
                          extra: <Download style={{ verticalAlign: 'middle' }} />,
                          onClick: () => {
                            downloadGlobalFile({ path: m.filePath }).then((data) => {
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
      destroyOnHidden
      forceRender={true}
    >
      {contextHolder}
      {formDom}
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
