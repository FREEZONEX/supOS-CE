import { type FC, useState } from 'react';
import { App, Button, Empty, Flex, Form, Pagination, Segmented } from 'antd';
import { useNavigate } from 'react-router';
import {
  getDashboardList,
  addDashboard,
  editDashboard,
  deleteDashboard,
  markDashboard,
  unmarkDashboard,
} from '@/apis/inter-api/uns';
import { usePagination, useTranslate } from '@/hooks';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import type { PageProps } from '@/common-types';
import { Dashboard, Edit, Grid, List, Search, TrashCan, View } from '@carbon/icons-react';
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

const CollectionFlow: FC<PageProps> = ({ title }) => {
  const { modal, message } = App.useApp();
  const formatMessage = useTranslate();
  const dashboardType = useBaseStore((state) => state.dashboardType);

  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [show, setShow] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [clickItem, setClickItem] = useState<any>({});
  const {
    loading,
    pagination,
    data: _data,
    reload,
    refreshRequest,
    setSearchParams,
    onChange,
  } = usePagination({
    fetchApi: getDashboardList,
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
    if (show) return;
    setShow(true);
  };
  const onDeleteHandle = (item: any) => {
    deleteDashboard(item.id)
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
    form
      .validateFields()
      .then((info) => {
        const params = info;
        if (isEdit) {
          params.id = clickItem.id;
        }
        const request = isEdit ? editDashboard : addDashboard;
        request(params)
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
        key: 'delete',
        label: formatMessage('common.delete'),
        auth: ButtonPermission['Dashboards.delete'],
        onClick: () => {
          modal.confirm({
            title: formatMessage('common.deleteConfirm'),
            onOk: () => onDeleteHandle(record),
            okButtonProps: {
              title: formatMessage('common.confirm'),
            },
            cancelButtonProps: {
              title: formatMessage('common.cancel'),
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
  };

  const pinOptions = {
    onClick: (record: any) => {
      const isMark = record?.mark === 1;
      const api = isMark ? unmarkDashboard : markDashboard;
      return api?.(record.id).then(() => {
        message.success(formatMessage('common.optsuccess'));
        reload();
      });
    },
    renderPinIcon: (record: any) => {
      return record?.mark !== 1;
    },
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
                onFinish: () => {
                  setSearchParams(searchForm.getFieldsValue());
                },
              }}
              onSearch={() => setSearchParams(searchForm.getFieldsValue())}
            />
            <AuthButton auth={ButtonPermission['Dashboards.add']} type="primary" onClick={onAddHandle}>
              + {formatMessage('dashboards.newDashboard')}
            </AuthButton>
          </>
        }
      >
        <Flex justify="flex-end" align="center" style={{ marginBottom: 16, marginTop: 16, paddingRight: 16 }}>
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
            data.length > 0 ? (
              <ProCardContainer>
                {data?.map((d: any) => {
                  return (
                    <ProCard
                      key={d?.id}
                      header={{
                        title: d.name,
                        titleDescription: formatTimestamp(d?.createTime),
                        onClick: hasPermission(ButtonPermission['Dashboards.preview'])
                          ? () =>
                              navigate(
                                `/dashboards/preview?${getSearchParamsString({ id: d.id, type: d.type, status: 'preview', name: d.name })}`
                              )
                          : undefined,
                      }}
                      statusHeader={{
                        statusTag: <div></div>,
                        pinOptions,
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
                      actions={actions}
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
                    width: '20%',
                    sorter: true,
                    render: (text: any, item: any) => {
                      const hasDesign = hasPermission(ButtonPermission['Dashboards.design']);
                      if (hasDesign)
                        return (
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
                        );
                      return text;
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
                    width: '30%',
                    ellipsis: true,
                  },
                  {
                    title: () => formatMessage('common.creationTime'),
                    dataIndex: 'createTime',
                    width: '15%',
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
    </ComLayout>
  );
};

export default CollectionFlow;
