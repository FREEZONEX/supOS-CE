import {
  addFlow,
  copyFlow,
  deleteFlow,
  editFlow,
  getEventFlowAndGroupList,
  markFlow,
  unmarkFlow,
  updateFlowStatus,
} from '@/apis/core-api/event-flow';
import { createFlowRowSorter } from '@/apis/core-api/flow';
import { deleteGroup, markGroup } from '@/apis/core-api/group.ts';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { AuthButton, ComEmptyState } from '@/components';
import ComButton from '@/components/com-button';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComSearch from '@/components/com-search';
import ViewModeSegmented from '@/components/lucide-icon/ViewModeSegmented';
import { PageTitleRow } from '@/components/lucide-icon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import OperationForm from '@/components/operation-form';
import ProCard from '@/components/pro-card/ProCard.tsx';
import ProCardContainer from '@/components/pro-card/ProCardContainer.tsx';
import ProModal from '@/components/pro-modal';
import SecondaryList from '@/components/pro-card/SecondaryList.tsx';
import ProTable from '@/components/pro-table';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { usePagination, useTranslate, useViewModeStorage, VIEW_MODE_STORAGE_KEYS } from '@/hooks';
import FlowStatusTag from '@/pages/flow/components/FlowStatusTag';
import FlowItemIcon from '@/pages/flow/components/FlowItemIcon';
import AddGroupModal from '@/components/group-modal/AddGroupModal.tsx';
import MoveGroupModal from '@/components/group-modal/MoveGroupModal.tsx';
import { getSearchParamsString, validInputPattern } from '@/utils';
import { hasPermission } from '@/utils/auth.ts';
import { formatTimestamp } from '@/utils/format.ts';
import {
  Add,
  CopyFile,
  Edit,
  FolderAdd,
  FolderMoveTo,
  PauseOutline,
  PlayOutline,
  Search,
  TrashCan,
  Undo,
} from '@/components/lucide-icon/carbon';
import { App, Breadcrumb, Button, Flex, Form, Pagination, Spin } from 'antd';
import { forwardRef, useEffect, useImperativeHandle, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { useLocation, useNavigate } from 'react-router';
import type { FlowListPanelProps } from '@/pages/collection-flow';
import type { FlowListPanelHandle } from '@/pages/flow/types';
import './index.scss';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';

type EventFlowPanelProps = Omit<FlowListPanelProps, 'flowScope'>;

const EventFlow = forwardRef<FlowListPanelHandle, EventFlowPanelProps>(function EventFlow(
  { title, embedded, viewMode: viewModeProp, onViewModeChange, hideViewMode, toolbarPortalHost, hideNewFlow },
  ref
) {
  const { modal, message } = App.useApp();
  const addGroupModalRef = useRef<any>(null);
  const moveGroupModalRef = useRef<any>(null);
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const location = useLocation();
  const [isEdit, setIsEdit] = useState('create');
  const [apiLoading, setApiLoading] = useState(false);
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [show, setShow] = useState(false);
  const [breadcrumbItem, setBreadcrumbItem] = useState<any>([{ title: formatMessage('common.all') }]);

  const renderFlowEmpty = (variant: 'page' | 'inline' = 'page') => (
    <ComEmptyState variant={variant} description={formatMessage('uns.noData')} />
  );

  const [storedMode, setStoredMode] = useViewModeStorage(VIEW_MODE_STORAGE_KEYS.eventFlow);
  const mode = viewModeProp ?? storedMode;
  const setMode = onViewModeChange ?? setStoredMode;
  const listFromPath = embedded ? '/flow' : location.pathname;
  const tableHostRef = useRef<HTMLDivElement>(null);
  const [tableScrollY, setTableScrollY] = useState<number>();

  const getEditorPath = (item: any) =>
    `/event-flow/flow-editor?${getSearchParamsString({
      id: item.id,
      name: item.flowName,
      status: item.flowStatus,
      flowId: item.runtimeFlowId || item.flowId,
      from: listFromPath,
    })}`;

  const { loading, pagination, data, reload, refreshRequest, setSearchParams, onChange } = usePagination({
    fetchApi: getEventFlowAndGroupList,
    initPageSize: 18,
  });

  useEffect(() => {
    if (!embedded || mode !== 'list') return;
    const host = tableHostRef.current;
    if (!host) return;

    const updateScrollY = () => {
      const next = Math.max(160, host.clientHeight - 110);
      setTableScrollY(next);
    };

    updateScrollY();
    const observer = new ResizeObserver(updateScrollY);
    observer.observe(host);
    return () => observer.disconnect();
  }, [embedded, mode, data?.length, pagination?.page, pagination?.pageSize]);

  const runStatusOptions = [
    {
      value: 'RUNNING',
      text: 'common.running',
    },
    {
      value: 'PENDING',
      text: 'common.pending',
    },
    {
      value: 'STOPPED',
      text: 'common.stopped',
    },
    {
      value: 'DISABLED',
      text: 'resourceAction.disabled',
    },
    {
      value: 'DRAFT',
      text: 'common.draft',
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
        {
          max: 128,
          message: formatMessage('uns.labelMaxLength', {
            label: formatMessage('common.name'),
            length: 128,
          }),
        },
        { pattern: validInputPattern, message: formatMessage('rule.flowNameIllegal') },
      ],
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

  useImperativeHandle(
    ref,
    () => ({
      refreshRequest: () => refreshRequest(),
      getCreateContext: () => ({
        groupId: breadcrumbItem?.length === 2 ? breadcrumbItem[1]?.groupId : undefined,
      }),
    }),
    [breadcrumbItem, refreshRequest]
  );

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
  const onStatusHandle = (item: any) => {
    const isDisabled = item.flowStatus === 'DISABLED';
    return updateFlowStatus(item.id, isDisabled ? 'deployed' : 'disabled').then(() => {
      message.success(formatMessage('common.optsuccess'));
      reload();
    });
  };
  const flowDeleteConfirmOptions = ({ content, name }: { content?: string; name?: string }) =>
    createDeleteConfirmOptions({
      title: formatMessage('common.deleteConfirm'),
      content,
      name,
      formatMessage,
      okText: formatMessage('common.confirm'),
      cancelText: formatMessage('common.cancel'),
    });
  const actions: any = (record: any) => {
    if (record?.category === 'group') {
      return [
        {
          key: 'edit',
          label: formatMessage('common.edit'),
          auth: ButtonPermission['EventFlow.edit'],
          icon: <Edit size={16} />,
          onClick: () => {
            addGroupModalRef?.current?.onOpen?.(2, {
              id: record.id,
              name: record.name,
              description: record.description,
            });
          },
        },
        { type: 'divider' },
        {
          key: 'delete',
          label: formatMessage('common.delete'),
          auth: ButtonPermission['EventFlow.delete'],
          icon: <TrashCan size={16} />,
          onClick: () => {
            modal.confirm({
              ...flowDeleteConfirmOptions({
                content: formatMessage('uns.deleteGroupInfo', { module: formatMessage('home.eventFlow') }),
              }),
              onOk: () => {
                return deleteGroup(record.id).then(() => {
                  message.success(formatMessage('common.optsuccess'));
                  reload();
                });
              },
            });
          },
        },
      ];
    }
    return [
      {
        key: 'copy',
        label: formatMessage('common.copy'),
        auth: ButtonPermission['EventFlow.copy'],
        icon: <CopyFile size={16} />,
        onClick: () => {
          setIsEdit('copy');
          setShow(true);
          form.setFieldsValue({
            id: record.id,
            groupId: record.groupId,
          });
        },
      },
      {
        key: 'edit',
        label: formatMessage('common.edit'),
        auth: ButtonPermission['EventFlow.edit'],
        icon: <Edit size={16} />,
        onClick: () => onEditHandle(record),
      },
      {
        key: 'moveToGroup',
        label: formatMessage('uns.moveToGroup'),
        auth: ButtonPermission['EventFlow.moveToGroup'],
        icon: <FolderMoveTo size={16} />,
        onClick: () => {
          moveGroupModalRef.current?.onOpen(2, { bizId: record.id, id: record.groupId });
        },
      },
      { type: 'divider' },
      {
        key: record.flowStatus === 'DISABLED' ? 'enable' : 'disable',
        label: formatMessage(record.flowStatus === 'DISABLED' ? 'common.enable' : 'common.disable'),
        auth: ButtonPermission['EventFlow.edit'],
        icon: record.flowStatus === 'DISABLED' ? <PlayOutline size={16} /> : <PauseOutline size={16} />,
        onClick: () => onStatusHandle(record),
      },
      {
        key: 'delete',
        label: formatMessage('common.delete'),
        auth: ButtonPermission['EventFlow.delete'],
        icon: <TrashCan size={16} />,
        onClick: () =>
          modal.confirm({
            ...flowDeleteConfirmOptions({ name: record.name || record.flowName || formatMessage('common.delete') }),
            onOk: () => onDeleteHandle(record),
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

  const toolbarExtra = (
    <>
      <ComSearch
        form={searchForm}
        formItemOptions={[
          {
            hidden: true,
            name: 'groupId',
          },
          {
            name: 'k',
            properties: {
              prefix: <Search {...toolbarIconProps} />,
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
          <FolderAdd {...toolbarIconProps} />
          {formatMessage('common.newFolder')}
        </AuthButton>
      )}
      {!hideNewFlow && (
        <AuthButton
          auth={ButtonPermission['EventFlow.add']}
          type="primary"
          icon={<Add {...toolbarIconProps} />}
          onClick={onAddHandle}
        >
          {formatMessage('common.newFlow')}
        </AuthButton>
      )}
    </>
  );

  const flowListColumns = [
    {
      titleIntlId: 'common.name',
      dataIndex: 'name',
      key: 'flowName',
      width: 280,
      maxWidth: 360,
      ellipsis: true,
      sorter: createFlowRowSorter('name'),
      sortKey: 'name',
      render: (text: any, item: any) => {
        const hasDesign = hasPermission(ButtonPermission['EventFlow.design']);
        return (
          <Flex gap={8} align="center" className="flow-list-name-cell">
            <FlowItemIcon category={item.category} flowKind={item.flowKind ?? 'event'} size="sm" />
            {item.category === 'group' ? (
              <Button
                type="link"
                className="table-link-button table-link-button-neutral"
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
                className="table-link-button table-link-button-neutral"
                type="link"
                onClick={() => {
                  navigate(getEditorPath(item));
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
      width: 120,
      maxWidth: 132,
      render: (flowStatus: any, item: any) => {
        if (flowStatus) {
          return <FlowStatusTag status={flowStatus}>{titleStatehandle(item)}</FlowStatusTag>;
        }
        return null;
      },
    },
    {
      titleIntlId: 'common.description',
      dataIndex: 'description',
      key: 'description',
      width: 240,
      maxWidth: 320,
      ellipsis: true,
    },
    {
      title: () => formatMessage('common.creationTime'),
      dataIndex: 'createAt',
      width: 180,
      maxWidth: 200,
      ellipsis: true,
      sorter: createFlowRowSorter('createAt'),
      sortKey: 'createAt',
      render: (item: any) => {
        const text = formatTimestamp(item);
        return (
          <span
            style={{
              display: 'block',
              maxWidth: '100%',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            }}
            title={text || undefined}
          >
            {text || '-'}
          </span>
        );
      },
    },
    {
      title: () => formatMessage('common.creator'),
      dataIndex: 'creator',
      width: 120,
      maxWidth: 160,
      ellipsis: true,
      render: (value: string) => (
        <span
          style={{
            display: 'block',
            maxWidth: '100%',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
          title={value || undefined}
        >
          {value || '-'}
        </span>
      ),
    },
  ] as any;

  const panelBody = (
    <>
      <Flex
        justify="space-between"
        align="center"
        style={{
          marginBottom: embedded ? 8 : 16,
          marginTop: embedded ? 0 : 16,
          padding: embedded ? 0 : '0 16px',
          flexShrink: 0,
        }}
      >
        <Flex align="center" gap={16}>
          {breadcrumbItem?.length > 1 ? (
            <>
              <Button
                size="small"
                style={{ background: 'var(--ui-switchwrap-bg-color)' }}
                onClick={() => {
                  searchForm.setFieldsValue({
                    groupId: undefined,
                  });
                  setBreadcrumbItem((pre: any) => pre.slice(0, -1));
                  onSearch?.();
                }}
              >
                <Undo {...toolbarIconProps} />
                {formatMessage('common.back')}
              </Button>
              <Breadcrumb items={breadcrumbItem} separator=">" />
            </>
          ) : (
            <span></span>
          )}
        </Flex>
        {!hideViewMode ? (
          <ViewModeSegmented
            value={mode}
            onChange={setMode}
            cardTitle={formatMessage('common.cardMode')}
            listTitle={formatMessage('common.listMode')}
          />
        ) : null}
      </Flex>
      <div
        className={embedded ? 'flow-list-panel-body' : undefined}
        style={{
          flex: 1,
          minHeight: 0,
          display: 'flex',
          flexDirection: 'column',
          padding: embedded ? 0 : '0 16px 16px',
          overflow: embedded && mode !== 'card' ? 'hidden' : 'auto',
        }}
      >
        {mode === 'card' ? (
          data?.length > 0 ? (
            <ProCardContainer>
              {data?.map((d: any) => {
                return (
                  <ProCard
                    key={d?.id}
                    iconBg={false}
                    classNames={{ headerTitle: 'flow-card-title' }}
                    header={{
                      customIcon: <FlowItemIcon category={d.category} flowKind={d.flowKind ?? 'event'} size="md" />,
                      title: d.name,
                      titleDescription: formatTimestamp(d?.createAt),
                      showChevron: false,
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
                            ? () => navigate(getEditorPath(d))
                            : undefined,
                    }}
                    statusHeader={{
                      statusTag:
                        d.category !== 'group' ? (
                          <FlowStatusTag status={d.flowStatus} title={titleStatehandle(d)} ellipsis>
                            {titleStatehandle(d)}
                          </FlowStatusTag>
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
                        ]}
                      />
                    }
                    item={d}
                  />
                );
              })}
            </ProCardContainer>
          ) : (
            renderFlowEmpty()
          )
        ) : (
          <div
            ref={tableHostRef}
            className={embedded ? 'flow-list-table-host' : undefined}
            style={{ flex: 1, minHeight: 0 }}
          >
            <ProTable
              locale={{ emptyText: renderFlowEmpty('inline') }}
              resizeable
              fixedPosition={embedded}
              style={{ height: embedded ? '100%' : undefined }}
              onChange={onChange}
              scroll={
                embedded
                  ? tableScrollY
                    ? { y: tableScrollY, x: 'max-content' }
                    : undefined
                  : { y: 'calc(100vh - 285px)', x: 'max-content' }
              }
              dataSource={data as any}
              columns={flowListColumns}
              pagination={{
                total: pagination?.total,
                style: { display: 'flex', justifyContent: 'flex-end', padding: embedded ? '10px 0 0' : '10px 0' },
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
          </div>
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
    </>
  );

  const panelModals = (
    <>
      <ProModal
        open={show}
        onCancel={onClose}
        title={formatMessage(`eventFlow.${isEdit}Flow`)}
        width={500}
        styles={{
          body: {
            paddingBlockStart: 0,
          },
        }}
      >
        {() => (
          <OperationForm
            loading={apiLoading}
            form={form}
            onCancel={onClose}
            onSave={onSave}
            formConfig={{
              layout: 'vertical',
              labelCol: { span: 24 },
              wrapperCol: { span: 24 },
            }}
            formItemOptions={formItemOptions(isEdit).filter((item: any) => item?.name)}
            style={{ padding: 0 }}
            footer={
              <Flex gap="10px" justify="end">
                <ComButton
                  color="default"
                  variant="filled"
                  onClick={onClose}
                  title={formatMessage('common.cancel')}
                >
                  {formatMessage('common.cancel')}
                </ComButton>
                <ComButton
                  type="primary"
                  variant="solid"
                  onClick={onSave}
                  title={formatMessage('common.save')}
                  loading={apiLoading}
                >
                  {formatMessage('common.save')}
                </ComButton>
              </Flex>
            }
          />
        )}
      </ProModal>
      <AddGroupModal ref={addGroupModalRef} refreshRequest={refreshRequest} />
      <MoveGroupModal ref={moveGroupModalRef} refreshRequest={refreshRequest} />
    </>
  );

  if (embedded) {
    return (
      <>
        {toolbarPortalHost ? createPortal(toolbarExtra, toolbarPortalHost) : null}
        <Spin spinning={loading}>
          <div
            className="flow-embedded-panel"
            style={{ display: 'flex', flexDirection: 'column', height: '100%', overflow: 'hidden' }}
          >
            {!toolbarPortalHost ? (
              <div
                style={{ display: 'flex', justifyContent: 'flex-end', gap: 10, padding: '0 16px 12px', flexShrink: 0 }}
              >
                {toolbarExtra}
              </div>
            ) : null}
            <div style={{ flex: 1, minHeight: 0, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              {panelBody}
            </div>
            {panelModals}
          </div>
        </Spin>
      </>
    );
  }

  return (
    <ComLayout loading={loading}>
      <ComContent
        title={
          <PageTitleRow resourceKey="flow.event.page">
            <span>{title}</span>
          </PageTitleRow>
        }
        mustHasBack={false}
        style={{
          overflow: 'hidden',
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
        }}
        extra={toolbarExtra}
      >
        {panelBody}
      </ComContent>
      {panelModals}
    </ComLayout>
  );
});

export default EventFlow;
