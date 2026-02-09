import { type FC, useRef, useState } from 'react';
import { App, Breadcrumb, Button, Empty, Flex, Form, Pagination, Segmented, Tag } from 'antd';
import { useNavigate } from 'react-router';
import {
  addFlow,
  copyFlow,
  deleteFlow,
  editFlow,
  getEventFlowAndGroupList,
  markFlow,
  unmarkFlow,
} from '@/apis/inter-api/event-flow';
import { usePagination, useTranslate } from '@/hooks';
import type { PageProps } from '@/common-types';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import {
  CopyFile,
  Edit,
  FlowData,
  Folder,
  FolderAdd,
  FolderMoveTo,
  Grid,
  List,
  Search,
  TrashCan,
  Undo,
} from '@carbon/icons-react';
import ComDrawer from '@/components/com-drawer';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComSearch from '@/components/com-search';
import OperationForm from '@/components/operation-form';
import { validInputPattern } from '@/utils/pattern';
import { getSearchParamsString } from '@/utils/url-util';
import { AuthButton } from '@/components/auth';
import { useLocalStorageState } from 'ahooks';
import ProCardContainer from '@/components/pro-card/ProCardContainer.tsx';
import ProTable from '@/components/pro-table';
import ProCard from '@/components/pro-card/ProCard.tsx';
import SecondaryList from '@/components/pro-card/SecondaryList.tsx';
import { hasPermission } from '@/utils/auth.ts';
import { formatTimestamp } from '@/utils/format.ts';
import AddGroupModal from '@/pages/dashboards/components/AddGroupModal.tsx';
import { deleteGroup, markGroup } from '@/apis/inter-api/group.ts';
import MoveGroupModal from '@/pages/dashboards/components/MoveGroupModal.tsx';

const CollectionFlow: FC<PageProps> = ({ title }) => {
  const { modal, message } = App.useApp();
  const addGroupModalRef = useRef<any>(null);
  const moveGroupModalRef = useRef<any>(null);
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const [isEdit, setIsEdit] = useState('create');
  const [apiLoading, setApiLoading] = useState(false);
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [show, setShow] = useState(false);
  const [breadcrumbItem, setBreadcrumbItem] = useState<any>([{ title: formatMessage('common.all') }]);
  const { loading, pagination, data, reload, refreshRequest, setSearchParams, onChange } = usePagination({
    fetchApi: getEventFlowAndGroupList,
    initPageSize: 18,
  });
  const [mode, setMode] = useLocalStorageState<string>('SUPOS_EVENTFLOW_MODE', {
    defaultValue: 'card',
  });
  const runStatusOptions = [
    {
      value: 'RUNNING',
      text: 'common.running',
      bgType: 'green',
    },
    {
      value: 'PENDING',
      text: 'common.pending',
      bgType: 'purple',
    },
    {
      value: 'STOPPED',
      text: 'common.stopped',
      bgType: 'red',
    },
    {
      value: 'DRAFT',
      text: 'common.draft',
      bgType: 'blue',
    },
  ];

  const titleStatehandle = (item: any) => {
    const key = runStatusOptions?.find((f: any) => f.value === item.flowStatus)?.text;
    return key ? formatMessage(key) : item.flowStatus;
  };

  const formItemOptions = (isEdit: string) => [
    {
      label: `${formatMessage(`eventFlow.${isEdit}Flow`)}`,
    },
    {
      name: 'groupId',
      hidden: true,
    },
    {
      label: formatMessage('common.name'),
      name: 'flowName',
      rules: [
        { required: true, message: formatMessage('rule.required') },
        { pattern: validInputPattern, message: formatMessage('rule.flowNameIllegal') },
      ],
    },
    {
      label: formatMessage('collectionFlow.flowTemplate'),
      name: 'template',
      type: 'Select',
      properties: {
        options: [
          {
            label: 'node-red',
            value: 'node-red',
          },
        ],
        disabled: isEdit !== 'create',
      },
      initialValue: 'node-red',
      rules: [{ required: true, message: '' }],
    },
    {
      label: formatMessage('uns.description'),
      name: 'description',
    },
    {
      label: 'id',
      name: 'id',
      hidden: true,
    },
    {
      type: 'divider',
    },
  ];
  useActivate(() => {
    refreshRequest?.();
  });
  const onClose = () => {
    setShow(false);
    form.resetFields();
  };
  const onAddHandle = () => {
    setIsEdit('create');
    form.resetFields();
    if (breadcrumbItem?.length === 2) {
      form.setFieldsValue({
        groupId: breadcrumbItem?.[1]?.groupId,
      });
    }
    if (show) return;
    setShow(true);
  };
  const onSave = async () => {
    const values = await form.validateFields();
    setApiLoading(true);
    const apiObj: any = {
      copy: copyFlow,
      edit: editFlow,
      create: addFlow,
    };
    const api = apiObj[isEdit || 'create'];
    return api({
      ...values,
      template: isEdit === 'edit' ? undefined : values.template,
      id: isEdit === 'edit' ? values.id : undefined,
      sourceId: isEdit === 'copy' ? values.id : undefined,
    })
      .then(() => {
        refreshRequest();
        message.success(formatMessage('common.optsuccess'));
        onClose();
      })
      .finally(() => {
        setApiLoading(false);
      });
  };
  const onDeleteHandle = (item: any) => {
    return deleteFlow(item.id).then(() => {
      message.success(formatMessage('common.deleteSuccessfully'));
      reload();
    });
  };
  const onEditHandle = (item: any) => {
    setIsEdit('edit');
    setShow(true);
    form.setFieldsValue({
      ...item,
    });
  };
  const actions: any = (record: any) => {
    if (record?.category === 'group') {
      return [
        {
          key: 'edit',
          label: formatMessage('common.edit'),
          auth: ButtonPermission['EventFlow.edit'],
          onClick: () => {
            addGroupModalRef?.current?.onOpen?.(2, {
              id: record.id,
              name: record.name,
              description: record.description,
            });
          },
          extra: (
            <Flex justify="center" align="center">
              <Edit />
            </Flex>
          ),
        },
        {
          key: 'delete',
          label: formatMessage('common.delete'),
          auth: ButtonPermission['EventFlow.delete'],
          onClick: () => {
            modal.confirm({
              title: formatMessage('common.deleteConfirm'),
              onOk: () => {
                return deleteGroup(record.id).then(() => {
                  message.success(formatMessage('common.optsuccess'));
                  reload();
                });
              },
            });
          },
          extra: (
            <Flex justify="center" align="center">
              <TrashCan />
            </Flex>
          ),
        },
      ];
    }
    return [
      {
        key: 'copy',
        label: formatMessage('common.copy'),
        auth: ButtonPermission['EventFlow.copy'],
        button: {
          type: 'primary',
        },
        extra: (
          <Flex justify="center" align="center">
            <CopyFile />
          </Flex>
        ),
        onClick: () => {
          setIsEdit('copy');
          setShow(true);
          form.setFieldsValue({
            id: record.id,
          });
        },
      },
      {
        key: 'edit',
        label: formatMessage('common.edit'),
        auth: ButtonPermission['EventFlow.edit'],
        extra: (
          <Flex justify="center" align="center">
            <Edit />
          </Flex>
        ),
        onClick: () => onEditHandle(record),
      },
      {
        key: 'moveToGroup',
        label: formatMessage('uns.moveToGroup'),
        auth: ButtonPermission['EventFlow.moveToGroup'],
        onClick: () => {
          moveGroupModalRef.current?.onOpen(2, { bizId: record.id, id: record.groupId });
        },
        extra: (
          <Flex justify="center" align="center">
            <FolderMoveTo />
          </Flex>
        ),
      },
      {
        key: 'delete',
        label: formatMessage('common.delete'),
        auth: ButtonPermission['EventFlow.delete'],
        extra: (
          <Flex justify="center" align="center">
            <TrashCan />
          </Flex>
        ),
        onClick: () =>
          modal.confirm({
            title: formatMessage('common.deleteConfirm'),
            onOk: () => onDeleteHandle(record),
            okButtonProps: {
              title: formatMessage('common.confirm'),
            },
            cancelButtonProps: {
              title: formatMessage('common.cancel'),
            },
          }),
      },
    ];
  };
  const pinOptions = {
    onClick: (record: any) => {
      const isMark = record?.sort === 1;
      if (record?.category === 'group') {
        return markGroup(record.id, !isMark).then(() => {
          message.success(formatMessage('common.optsuccess'));
          reload();
        });
      } else {
        const api = isMark ? unmarkFlow : markFlow;
        return api?.(record.id).then(() => {
          message.success(formatMessage('common.optsuccess'));
          reload();
        });
      }
    },
    renderPinIcon: (record: any) => {
      return record?.sort !== 1;
    },
  };

  const onSearch = () => {
    const params = searchForm.getFieldsValue();
    if (params?.category === 'all') {
      params.category = undefined;
    }
    setSearchParams(params);
  };

  return (
    <ComLayout loading={loading}>
      <ComContent
        title={
          <Flex align="center" gap={8}>
            <FlowData size={20} />
            <span>{title}</span>
          </Flex>
        }
        mustHasBack={false}
        style={{
          overflow: 'hidden',
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
        }}
        extra={
          <>
            <ComSearch
              form={searchForm}
              formItemOptions={[
                {
                  hidden: true,
                  name: 'groupId',
                },
                {
                  type: 'Filter',
                  name: 'category',
                  properties: {
                    placeholder: formatMessage('common.searchPlaceholder'),
                    defaultValue: 'all',
                  },
                },
                {
                  name: 'k',
                  properties: {
                    prefix: <Search />,
                    placeholder: formatMessage('common.searchPlaceholder'),
                    style: { width: 300 },
                    allowClear: true,
                  },
                },
              ]}
              formConfig={{
                onFinish: onSearch,
              }}
              onSearch={onSearch}
            />
            {breadcrumbItem?.length === 1 && (
              <AuthButton
                auth={ButtonPermission['EventFlow.add']}
                type="primary"
                onClick={() => {
                  addGroupModalRef?.current?.onOpen?.(2);
                }}
              >
                <FolderAdd />
                {formatMessage('common.newGroup')}
              </AuthButton>
            )}
            <AuthButton auth={ButtonPermission['EventFlow.add']} type="primary" onClick={onAddHandle}>
              + {formatMessage('eventFlow.newFlow')}
            </AuthButton>
          </>
        }
      >
        <Flex justify="space-between" align="center" style={{ marginBottom: 16, marginTop: 16, padding: '0 16px' }}>
          <Flex align="center" gap={16}>
            {breadcrumbItem?.length > 1 && (
              <Button
                size="small"
                style={{ background: 'var(--supos-switchwrap-bg-color)' }}
                onClick={() => {
                  searchForm.setFieldsValue({
                    groupId: undefined,
                  });
                  setBreadcrumbItem((pre: any) => pre.slice(0, -1));
                  onSearch?.();
                }}
              >
                <Undo />
                {formatMessage('common.back')}
              </Button>
            )}
            <Breadcrumb items={breadcrumbItem} separator=">" />
          </Flex>
          <Segmented
            size="small"
            value={mode}
            onChange={(v) => setMode(v)}
            options={[
              {
                value: 'card',
                icon: (
                  <span title={formatMessage('common.cardMode')}>
                    <Grid />
                  </span>
                ),
              },
              {
                value: 'list',
                icon: (
                  <span title={formatMessage('common.listMode')}>
                    <List />
                  </span>
                ),
              },
            ]}
          />
        </Flex>
        <div style={{ flex: 1, padding: '0 16px 16px', overflow: 'auto', alignItems: 'center' }}>
          {mode === 'card' ? (
            data?.length > 0 ? (
              <ProCardContainer>
                {data?.map((d: any) => {
                  return (
                    <ProCard
                      key={d?.id}
                      header={{
                        customIconBg: d.category === 'group' ? '#A8A8A8' : undefined,
                        customIcon: (
                          <Flex align="center" justify="center">
                            {d.category === 'group' ? (
                              <Folder size="28" style={{ color: 'white' }} />
                            ) : (
                              <FlowData size="28" />
                            )}
                          </Flex>
                        ),
                        title: d.name,
                        titleDescription: formatTimestamp(d?.createAt),
                        onClick:
                          d.category === 'group'
                            ? () => {
                                setBreadcrumbItem((pre: any) => {
                                  return [...pre, { title: d.name, groupId: d.id }];
                                });
                                searchForm.setFieldsValue({
                                  groupId: d.id,
                                });
                                onSearch?.();
                              }
                            : hasPermission(ButtonPermission['EventFlow.design'])
                              ? () =>
                                  navigate(
                                    `/EventFlow/Editor?${getSearchParamsString({ id: d.id, name: d.flowName, status: d.flowStatus, flowId: d.flowId })}`
                                  )
                              : undefined,
                      }}
                      statusHeader={{
                        statusTag:
                          d.category !== 'group' ? (
                            <Tag
                              style={{ borderRadius: 15, lineHeight: '16px', margin: 0 }}
                              bordered={false}
                              color={
                                (runStatusOptions?.find((f: any) => f.value === d.flowStatus)?.bgType || 'red') as any
                              }
                            >
                              {titleStatehandle(d)}
                            </Tag>
                          ) : (
                            <div></div>
                          ),
                        pinOptions,
                        actions,
                      }}
                      description={d?.description}
                      secondaryDescription={
                        <SecondaryList
                          options={[
                            {
                              label: formatMessage('common.creator'),
                              content: d?.creator,
                              span: 24,
                              key: 'creator',
                            },
                            {
                              label: formatMessage('collectionFlow.flowTemplate'),
                              content: d?.template,
                              span: 24,
                              key: 'flowTemplate',
                            },
                          ]}
                        />
                      }
                      item={d}
                    />
                  );
                })}
              </ProCardContainer>
            ) : (
              <Empty />
            )
          ) : (
            <ProTable
              resizeable
              style={{ height: '100%' }}
              onChange={onChange}
              scroll={{ y: 'calc(100vh  - 285px)', x: 'max-content' }}
              dataSource={data as any}
              columns={
                [
                  {
                    titleIntlId: 'common.name',
                    dataIndex: 'name',
                    key: 'flowName',
                    width: 200,
                    sorter: true,
                    render: (text: any, item: any) => {
                      const hasDesign = hasPermission(ButtonPermission['EventFlow.design']);
                      return (
                        <Flex gap={8} align="center">
                          <Flex
                            style={{
                              borderRadius: 3,
                              backgroundColor: item.category === 'group' ? '#A8A8A8' : '#F4F4F4',
                              padding: 6,
                              height: 26,
                            }}
                          >
                            <Flex align="center" justify="center">
                              {item.category === 'group' ? (
                                <Folder size="16" style={{ color: 'white' }} />
                              ) : (
                                <FlowData size="16" style={{ color: 'black' }} />
                              )}
                            </Flex>
                          </Flex>
                          {item.category === 'group' ? (
                            <Button
                              type="link"
                              className="table-link-button"
                              onClick={() => {
                                searchForm.setFieldsValue({
                                  groupId: item.id,
                                });
                                setBreadcrumbItem((pre: any) => {
                                  return [...pre, { title: item.name, groupId: item.id }];
                                });
                                onSearch?.();
                              }}
                              title={text}
                            >
                              {text}
                            </Button>
                          ) : hasDesign ? (
                            <Button
                              className="table-link-button"
                              type="link"
                              onClick={() => {
                                navigate(
                                  `/EventFlow/Editor?${getSearchParamsString({ id: item.id, name: item.flowName, status: item.flowStatus, flowId: item.flowId })}`
                                );
                              }}
                              title={text}
                            >
                              {text}
                            </Button>
                          ) : (
                            text
                          )}
                        </Flex>
                      );
                    },
                  },
                  {
                    title: () => formatMessage('common.status'),
                    dataIndex: 'flowStatus',
                    width: 150,
                    render: (flowStatus: any, item: any) => {
                      if (flowStatus) {
                        return (
                          <Tag
                            style={{ borderRadius: 15, lineHeight: '16px', margin: 0 }}
                            bordered={false}
                            color={(runStatusOptions?.find((f: any) => f.value === flowStatus)?.bgType || 'red') as any}
                          >
                            {titleStatehandle(item)}
                          </Tag>
                        );
                      } else {
                        return null;
                      }
                    },
                  },
                  {
                    titleIntlId: 'collectionFlow.flowTemplate',
                    dataIndex: 'template',
                    key: 'template',
                    width: 150,
                  },
                  {
                    titleIntlId: 'common.description',
                    dataIndex: 'description',
                    key: 'description',
                    width: 400,
                    ellipsis: true,
                  },
                  {
                    title: () => formatMessage('common.creationTime'),
                    dataIndex: 'createAt',
                    width: 200,
                    sorter: true,
                    render: (item: any) => formatTimestamp(item),
                  },
                  {
                    title: () => formatMessage('common.creator'),
                    dataIndex: 'creator',
                    width: 150,
                  },
                ] as any
              }
              pagination={{
                total: pagination?.total,
                style: { display: 'flex', justifyContent: 'flex-end', padding: '10px 0' },
                pageSize: pagination?.pageSize || 18,
                current: pagination?.page,
                showQuickJumper: true,
                pageSizeOptions: pagination?.pageSizes,
                showSizeChanger: true,
                onChange: pagination.onChange,
                onShowSizeChange: (current, size) => {
                  pagination.onChange({ page: current, pageSize: size });
                },
              }}
              pinOptions={pinOptions}
              operationOptions={{
                render: actions,
              }}
            />
          )}
        </div>
        {mode === 'card' && (
          <Pagination
            size="small"
            className="custom-pagination"
            align="center"
            style={{ margin: '20px 0' }}
            total={pagination?.total}
            showSizeChanger={false}
            onChange={pagination.onChange}
            pageSize={pagination?.pageSize || 18}
            current={pagination?.page}
          />
        )}
      </ComContent>
      <ComDrawer title=" " open={show} onClose={onClose}>
        <OperationForm
          loading={apiLoading}
          form={form}
          onCancel={onClose}
          onSave={onSave}
          formItemOptions={formItemOptions(isEdit)}
        />
      </ComDrawer>
      <AddGroupModal ref={addGroupModalRef} refreshRequest={refreshRequest} />
      <MoveGroupModal ref={moveGroupModalRef} refreshRequest={refreshRequest} />
    </ComLayout>
  );
};

export default CollectionFlow;
