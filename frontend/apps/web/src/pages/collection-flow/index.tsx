import { type FC, useState } from 'react';
import { App, Button, Empty, Flex, Form, Pagination, Segmented, Tag } from 'antd';
import { useNavigate } from 'react-router';
import { addFlow, copyFlow, deleteFlow, editFlow, flowPage, markFlow, unmarkFlow } from '@/apis/inter-api/flow';
import { usePagination, useTranslate } from '@/hooks';
import type { PageProps } from '@/common-types';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { ConnectSource, CopyFile, Edit, Grid, List, Search, TrashCan } from '@carbon/icons-react';
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
import ProCard from '../../components/pro-card/ProCard.tsx';
import SecondaryList from '../../components/pro-card/SecondaryList.tsx';
import { hasPermission } from '@/utils/auth.ts';
import { formatTimestamp } from '@/utils/format.ts';

const CollectionFlow: FC<PageProps> = ({ title }) => {
  const { modal, message } = App.useApp();
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const [isEdit, setIsEdit] = useState('create');
  const [apiLoading, setApiLoading] = useState(false);
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [show, setShow] = useState(false);
  const [mode, setMode] = useLocalStorageState<string>('SUPOS_COLLECTION_MODE', {
    defaultValue: 'list',
  });

  const { loading, pagination, data, reload, refreshRequest, setSearchParams, onChange } = usePagination({
    fetchApi: flowPage,
    initPageSize: 18,
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
      label: `${formatMessage(`collectionFlow.${isEdit}Flow`)}`,
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
    api({
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
    deleteFlow(item.id).then(() => {
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
    return [
      {
        key: 'copy',
        label: formatMessage('common.copy'),
        auth: ButtonPermission['SourceFlow.copy'],
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
        auth: ButtonPermission['SourceFlow.edit'],
        extra: (
          <Flex justify="center" align="center">
            <Edit />
          </Flex>
        ),
        onClick: () => onEditHandle(record),
      },
      {
        key: 'delete',
        label: formatMessage('common.delete'),
        auth: ButtonPermission['SourceFlow.delete'],
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
      const isMark = record?.mark === 1;
      const api = isMark ? unmarkFlow : markFlow;
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
        title={
          <Flex align="center" gap={8}>
            <ConnectSource size={20} />
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
              onSearch={() => {
                setSearchParams(searchForm.getFieldsValue());
              }}
            />
            <AuthButton auth={ButtonPermission['SourceFlow.add']} type="primary" onClick={onAddHandle}>
              + {formatMessage('collectionFlow.newFlow')}
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
                        title: d.flowName,
                        titleDescription: formatTimestamp(d?.createTime),
                        onClick: hasPermission(ButtonPermission['SourceFlow.design'])
                          ? () =>
                              navigate(
                                `/collection-flow/flow-editor?${getSearchParamsString({ id: d.id, name: d.flowName, status: d.flowStatus, flowId: d.flowId, from: location.pathname })}`
                              )
                          : undefined,
                      }}
                      statusHeader={{
                        statusTag: (
                          <Tag
                            style={{
                              borderRadius: 9,
                              height: 16,
                              lineHeight: '16px',
                              maxWidth: 120,
                              overflow: 'hidden',
                              whiteSpace: 'nowrap',
                              textOverflow: 'ellipsis',
                            }}
                            bordered={false}
                            title={titleStatehandle(d)}
                            color={
                              (runStatusOptions?.find((f: any) => f.value === d.flowStatus)?.bgType || 'red') as any
                            }
                          >
                            {titleStatehandle(d)}
                          </Tag>
                        ),
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
                              label: formatMessage('collectionFlow.flowTemplate'),
                              content: d?.template,
                              span: 24,
                              key: 'flowTemplate',
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
              columns={
                [
                  {
                    titleIntlId: 'common.name',
                    dataIndex: 'flowName',
                    width: '20%',
                    sorter: true,
                    render: (text: any, item: any) => {
                      const hasDesign = hasPermission(ButtonPermission['SourceFlow.design']);
                      return (
                        <>
                          {hasDesign ? (
                            <Button
                              className="table-link-button"
                              type="link"
                              onClick={() => {
                                navigate(
                                  `/collection-flow/flow-editor?${getSearchParamsString({ id: item.id, name: item.flowName, status: item.flowStatus, flowId: item.flowId, from: location.pathname })}`
                                );
                              }}
                              title={text}
                            >
                              {text}
                            </Button>
                          ) : (
                            text
                          )}
                          {item.flowStatus && (
                            <Tag
                              style={{ borderRadius: 15, lineHeight: '16px', margin: 0 }}
                              bordered={false}
                              color={
                                (runStatusOptions?.find((f: any) => f.value === item.flowStatus)?.bgType ||
                                  'red') as any
                              }
                            >
                              {titleStatehandle(item)}
                            </Tag>
                          )}
                        </>
                      );
                    },
                  },
                  {
                    titleIntlId: 'collectionFlow.flowTemplate',
                    dataIndex: 'template',
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
    </ComLayout>
  );
};

export default CollectionFlow;
