import { type FC, useRef, useState } from 'react';
import { App, Breadcrumb, Button, Empty, Flex, Form, Pagination, Segmented } from 'antd';
import { useNavigate } from 'react-router';
import {
  markDashboard,
  unmarkDashboard,
  addDashboard,
  editDashboard,
  deleteDashboard,
  getDashboardAndGroupList,
} from '@/apis/inter-api/dashboard';
import { usePagination, useTranslate } from '@/hooks';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import type { PageProps } from '@/common-types';
import {
  Dashboard,
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
  View,
} from '@carbon/icons-react';
import ComDrawer from '@/components/com-drawer';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComSearch from '@/components/com-search';
import OperationForm from '@/components/operation-form';
import { validInputPattern } from '@/utils/pattern';
import { getSearchParamsString } from '@/utils/url-util';
import { useBaseStore } from '@/stores/base';
import { AuthButton } from '@/components/auth';
import './index.scss';
import { useLocalStorageState } from 'ahooks';
import ProCardContainer from '@/components/pro-card/ProCardContainer.tsx';
import ProTable from '@/components/pro-table';
import ProCard from '@/components/pro-card/ProCard.tsx';
import SecondaryList from '../../components/pro-card/SecondaryList.tsx';
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
  const dashboardType = useBaseStore((state) => state.dashboardType);

  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [show, setShow] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [clickItem, setClickItem] = useState<any>({});
  const [breadcrumbItem, setBreadcrumbItem] = useState<any>([{ title: formatMessage('common.all') }]);
  const {
    loading,
    pagination,
    data: _data,
    reload,
    refreshRequest,
    setSearchParams,
    onChange,
  } = usePagination({
    fetchApi: getDashboardAndGroupList,
    initPageSize: 18,
  });
  const [mode, setMode] = useLocalStorageState<string>('SUPOS_DASHBOARD_MODE', {
    defaultValue: 'list',
  });

  const typeOptions = [
    {
      label: 'Grafana',
      value: 1,
      key: 'grafana',
    },
    {
      label: 'Fuxa',
      value: 2,
      key: 'fuxa',
    },
  ]?.filter((f) => dashboardType?.includes(f.key));
  const data = (_data || [])?.map((e: any) => {
    return {
      ...e,
      typeName: e.type
        ? typeOptions?.find((o) => o.value === e.type)?.label
        : typeOptions?.find((o) => o.value === 1)?.label,
    };
  });
  const formItemOptions = (isEdit: boolean) => {
    return [
      {
        label: `${isEdit ? formatMessage('common.edit') : formatMessage('common.create')} ${formatMessage('dashboards.dashboard')}`,
      },
      {
        name: 'groupId',
        hidden: true,
      },
      {
        label: formatMessage('common.name'),
        name: 'name',
        rules: [
          { required: true, message: '' },
          { pattern: validInputPattern, message: '' },
        ],
      },
      {
        label: formatMessage('dashboards.dashboardsTemplate'),
        name: 'type',
        type: 'Select',
        properties: {
          options: typeOptions,
          disabled: isEdit,
        },
        rules: [{ required: true, message: '' }],
      },
      {
        label: formatMessage('uns.description'),
        name: 'description',
      },
      {
        type: 'divider',
      },
    ];
  };
  useActivate(() => {
    refreshRequest?.();
  });

  const onClose = () => {
    setShow(false);
    form.resetFields();
  };
  const onAddHandle = () => {
    form.resetFields();
    setIsEdit(false);
    if (breadcrumbItem?.length === 2) {
      form.setFieldsValue({
        groupId: breadcrumbItem?.[1]?.groupId,
      });
    }
    if (show) return;
    setShow(true);
  };
  const onDeleteHandle = (item: any) => {
    return deleteDashboard(item.id)
      .then(() => {
        message.success(formatMessage('common.deleteSuccessfully'));
        reload();
      })
      .catch(() => {});
  };
  const edit = (item: any) => {
    form.setFieldsValue({ name: item.name, type: item.type || 1, description: item.description });
    setShow(true);
    setClickItem(item);
  };

  const onSave = () => {
    return form
      .validateFields()
      .then((info) => {
        const params = info;
        if (isEdit) {
          params.id = clickItem.id;
        }
        const request = isEdit ? editDashboard : addDashboard;
        return request(params)
          .then(() => {
            message.success(formatMessage('common.optsuccess'));
            onClose();
            refreshRequest();
          })
          .catch(() => {});
      })
      .catch(() => {});
  };

  const onEditHandle = (item: any) => {
    setIsEdit(true);
    edit(item);
  };

  const actions: any = (record: any) => {
    if (record?.category === 'group') {
      return [
        {
          key: 'edit',
          label: formatMessage('common.edit'),
          auth: ButtonPermission['Dashboards.edit'],
          onClick: () => {
            addGroupModalRef?.current?.onOpen?.(3, {
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
          auth: ButtonPermission['Dashboards.delete'],
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
        key: 'preview',
        label: formatMessage('dashboards.preview'),
        auth: ButtonPermission['Dashboards.preview'],
        button: {
          type: 'primary',
        },
        onClick: () => {
          setClickItem(record);
          navigate(
            `/dashboards/preview?${getSearchParamsString({ id: record.id, type: record.type, status: 'preview', name: record.name })}`
          );
        },
        extra: (
          <Flex justify="center" align="center">
            <View />
          </Flex>
        ),
      },
      {
        key: 'edit',
        label: formatMessage('common.edit'),
        auth: ButtonPermission['Dashboards.edit'],
        onClick: () => {
          onEditHandle(record);
        },
        extra: (
          <Flex justify="center" align="center">
            <Edit />
          </Flex>
        ),
      },
      {
        key: 'moveToGroup',
        label: formatMessage('uns.moveToGroup'),
        auth: ButtonPermission['Dashboards.moveToGroup'],
        onClick: () => {
          moveGroupModalRef.current?.onOpen(3, { bizId: record.id, id: record.groupId });
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
        auth: ButtonPermission['Dashboards.delete'],
        onClick: () => {
          modal.confirm({
            title: formatMessage('common.deleteConfirm'),
            onOk: () => onDeleteHandle(record),
          });
        },
        extra: (
          <Flex justify="center" align="center">
            <TrashCan />
          </Flex>
        ),
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
        const api = isMark ? unmarkDashboard : markDashboard;
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
        mustHasBack={false}
        title={
          <Flex align="center" gap={8}>
            <Dashboard size={20} />
            <span>{title}</span>
          </Flex>
        }
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
                auth={ButtonPermission['SourceFlow.add']}
                type="primary"
                onClick={() => {
                  addGroupModalRef?.current?.onOpen?.(3);
                }}
              >
                <FolderAdd />
                {formatMessage('common.newGroup')}
              </AuthButton>
            )}
            <AuthButton auth={ButtonPermission['Dashboards.add']} type="primary" onClick={onAddHandle}>
              + {formatMessage('dashboards.newDashboard')}
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
                            : hasPermission(ButtonPermission['Dashboards.preview'])
                              ? () => {
                                  setClickItem(d);
                                  navigate(
                                    `/dashboards/preview?${getSearchParamsString({ id: d.id, type: d.type, status: 'preview', name: d.name })}`
                                  );
                                }
                              : undefined,
                      }}
                      statusHeader={{
                        statusTag: <div></div>,
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
                              label: formatMessage('dashboards.dashboardsTemplate'),
                              content: d?.typeName,
                              span: 24,
                              key: 'dashboardsTemplate',
                            },
                          ]}
                        />
                      }
                      // actions={actions}
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
              onChange={onChange}
              style={{ height: '100%' }}
              scroll={{ y: 'calc(100vh  - 285px)', x: 'max-content' }}
              dataSource={data as any}
              pinOptions={pinOptions}
              columns={
                [
                  {
                    titleIntlId: 'common.name',
                    dataIndex: 'name',
                    width: '14%',
                    sorter: true,
                    render: (text: any, item: any) => {
                      const hasDesign = hasPermission(ButtonPermission['Dashboards.design']);
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
                              type="link"
                              className="table-link-button"
                              onClick={() => {
                                setClickItem(item);
                                navigate(
                                  `/dashboards/preview?${getSearchParamsString({ id: item.id, type: item.type, status: 'design', name: item.name })}`
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
                    titleIntlId: 'dashboards.dashboardsTemplate',
                    dataIndex: 'typeName',
                    width: '10%',
                  },
                  {
                    titleIntlId: 'common.description',
                    dataIndex: 'description',
                    width: '20%',
                    ellipsis: true,
                  },
                  {
                    title: () => formatMessage('common.creationTime'),
                    dataIndex: 'createAt',
                    width: '10%',
                    sorter: true,
                    render: (item: any) => formatTimestamp(item),
                  },
                  {
                    title: () => formatMessage('common.creator'),
                    dataIndex: 'creator',
                    width: '10%',
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
        <OperationForm form={form} onCancel={onClose} onSave={onSave} formItemOptions={formItemOptions(isEdit)} />
      </ComDrawer>
      <AddGroupModal ref={addGroupModalRef} refreshRequest={refreshRequest} />
      <MoveGroupModal ref={moveGroupModalRef} refreshRequest={refreshRequest} />
    </ComLayout>
  );
};

export default CollectionFlow;
