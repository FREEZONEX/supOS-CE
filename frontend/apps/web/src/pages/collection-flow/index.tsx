import { forwardRef, useCallback, useEffect, useImperativeHandle, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { App, Breadcrumb, Button, Flex, Form, Pagination, Spin } from 'antd';
import { useLocation, useNavigate } from 'react-router';
import {
  addFlow,
  copyFlow,
  deleteFlow,
  editFlow,
  createFlowRowSorter,
  getFlowAndGroupList,
  markFlow,
  unmarkFlow,
  updateFlowStatus,
} from '@/apis/core-api/flow';
import {
  addFlow as addEventFlow,
  copyFlow as copyEventFlow,
  deleteFlow as deleteEventFlow,
  editFlow as editEventFlow,
} from '@/apis/core-api/event-flow';
import { usePagination, useTranslate, useViewModeStorage, VIEW_MODE_STORAGE_KEYS } from '@/hooks';
import type { PageProps } from '@/common-types';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { ButtonPermission } from '@/common-types/button-permission.ts';
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
import FlowStatusTag from '@/pages/flow/components/FlowStatusTag';
import FlowItemIcon from '@/pages/flow/components/FlowItemIcon';
import ComButton from '@/components/com-button';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComSearch from '@/components/com-search';
import ViewModeSegmented from '@/components/lucide-icon/ViewModeSegmented';
import { PageTitleRow } from '@/components/lucide-icon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import OperationForm from '@/components/operation-form';
import { validInputPattern } from '@/utils/pattern';
import { getSearchParamsString } from '@/utils/url-util';
import { AuthButton, ComEmptyState } from '@/components';
import ProCardContainer from '@/components/pro-card/ProCardContainer.tsx';
import ProModal from '@/components/pro-modal';
import ProTable from '@/components/pro-table';
import ProCard from '../../components/pro-card/ProCard.tsx';
import SecondaryList from '../../components/pro-card/SecondaryList.tsx';
import { hasPermission } from '@/utils/auth.ts';
import { formatTimestamp } from '@/utils/format.ts';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import { deleteGroup, markGroup } from '@/apis/core-api/group.ts';
import AddGroupModal from '@/components/group-modal/AddGroupModal';
import MoveGroupModal from '@/components/group-modal/MoveGroupModal';
import { fetchMergedAllFlows, flowListRowKey } from '@/pages/flow/merge-all-flows';
import type { FlowListPanelHandle } from '@/pages/flow/types';
import './index.scss';

export interface FlowListPanelProps extends PageProps {
  embedded?: boolean;
  flowScope?: 'source' | 'all';
  viewMode?: string;
  onViewModeChange?: (mode: string) => void;
  hideViewMode?: boolean;
  toolbarPortalHost?: HTMLElement | null;
  hideNewFlow?: boolean;
}

const CollectionFlow = forwardRef<FlowListPanelHandle, FlowListPanelProps>(function CollectionFlow(
  {
    title,
    embedded,
    flowScope = 'source',
    viewMode: viewModeProp,
    onViewModeChange,
    hideViewMode,
    toolbarPortalHost,
    hideNewFlow,
  },
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
  const [storedMode, setStoredMode] = useViewModeStorage(VIEW_MODE_STORAGE_KEYS.collection);
  const mode = viewModeProp ?? storedMode;
  const setMode = onViewModeChange ?? setStoredMode;
  const [breadcrumbItem, setBreadcrumbItem] = useState<any>([{ title: formatMessage('common.all') }]);
  const tableHostRef = useRef<HTMLDivElement>(null);
  const [tableScrollY, setTableScrollY] = useState<number>();
  const readOnlyList = flowScope === 'all';
  const listFromPath = embedded ? '/flow' : location.pathname;

  const renderFlowEmpty = (variant: 'page' | 'inline' = 'page') => (
    <ComEmptyState variant={variant} description={formatMessage('uns.noData')} />
  );

  const fetchAllFlows = useCallback(async (params: any) => fetchMergedAllFlows(params), []);

  const getEditorPath = (item: any) => {
    const editorBase = item?.flowKind === 'event' ? '/event-flow/flow-editor' : '/collection-flow/flow-editor';
    return `${editorBase}?${getSearchParamsString({
      id: item.id,
      name: item.flowName,
      status: item.flowStatus,
      flowId: item.runtimeFlowId || item.flowId,
      from: listFromPath,
    })}`;
  };

  const { loading, pagination, data, reload, refreshRequest, setSearchParams, onChange } = usePagination({
    fetchApi: flowScope === 'all' ? fetchAllFlows : getFlowAndGroupList,
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
      label: `${formatMessage(`collectionFlow.${isEdit}Flow`)}`,
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
      name: 'id',
      hidden: true,
    },
    {
      name: 'flowKind',
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
    const isEvent = values.flowKind === 'event';
    const apiObj: any = isEvent
      ? { copy: copyEventFlow, edit: editEventFlow, create: addEventFlow }
      : { copy: copyFlow, edit: editFlow, create: addFlow };
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
    const api = item?.flowKind === 'event' ? deleteEventFlow : deleteFlow;
    return api(item.id).then(() => {
      message.success(formatMessage('common.deleteSuccessfully'));
      reload();
    });
  };
  const onEditHandle = (item: any) => {
    setIsEdit('edit');
    setShow(true);
    form.setFieldsValue({
      ...item,
      flowKind: item?.flowKind || 'source',
    });
  };
  const isEventFlow = (record: any) => record?.flowKind === 'event';
  const flowGroupType = (record: any) => (isEventFlow(record) ? 2 : 1);
  const flowPerm = (record: any) => ({
    edit: isEventFlow(record) ? ButtonPermission['EventFlow.edit'] : ButtonPermission['SourceFlow.edit'],
    copy: isEventFlow(record) ? ButtonPermission['EventFlow.copy'] : ButtonPermission['SourceFlow.copy'],
    delete: isEventFlow(record) ? ButtonPermission['EventFlow.delete'] : ButtonPermission['SourceFlow.delete'],
    moveToGroup: isEventFlow(record)
      ? ButtonPermission['EventFlow.moveToGroup']
      : ButtonPermission['SourceFlow.moveToGroup'],
  });
  const flowModuleName = (record: any) =>
    isEventFlow(record) ? formatMessage('home.eventFlow') : formatMessage('home.sourceFlow');
  const flowDeleteConfirmOptions = ({ content, name }: { content?: string; name?: string }) =>
    createDeleteConfirmOptions({
      title: formatMessage('common.deleteConfirm'),
      content,
      name,
      formatMessage,
      okText: formatMessage('common.confirm'),
      cancelText: formatMessage('common.cancel'),
    });
  const onStatusHandle = (item: any) => {
    const isDisabled = item.flowStatus === 'DISABLED';
    return updateFlowStatus(item.id, isDisabled ? 'deployed' : 'disabled').then(() => {
      message.success(formatMessage('common.optsuccess'));
      reload();
    });
  };
  const actions: any = (record: any) => {
    if (record?.category === 'group') {
      return [
        {
          key: 'edit',
          label: formatMessage('common.edit'),
          auth: flowPerm(record).edit,
          icon: <Edit size={16} />,
          onClick: () => {
            addGroupModalRef?.current?.onOpen?.(flowGroupType(record), {
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
          auth: flowPerm(record).delete,
          icon: <TrashCan size={16} />,
          onClick: () => {
            modal.confirm({
              ...flowDeleteConfirmOptions({
                content: formatMessage('uns.deleteGroupInfo', { module: flowModuleName(record) }),
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
        auth: flowPerm(record).copy,
        icon: <CopyFile size={16} />,
        onClick: () => {
          setIsEdit('copy');
          setShow(true);
          form.setFieldsValue({
            id: record.id,
            groupId: record.groupId,
            flowKind: record?.flowKind || 'source',
          });
        },
      },
      {
        key: 'edit',
        label: formatMessage('common.edit'),
        auth: flowPerm(record).edit,
        icon: <Edit size={16} />,
        onClick: () => onEditHandle(record),
      },
      {
        key: 'moveToGroup',
        label: formatMessage('uns.moveToGroup'),
        auth: flowPerm(record).moveToGroup,
        icon: <FolderMoveTo size={16} />,
        onClick: () => {
          moveGroupModalRef.current?.onOpen(flowGroupType(record), { bizId: record.id, id: record.groupId });
        },
      },
      { type: 'divider' },
      {
        key: record.flowStatus === 'DISABLED' ? 'enable' : 'disable',
        label: formatMessage(record.flowStatus === 'DISABLED' ? 'common.enable' : 'common.disable'),
        auth: ButtonPermission['SourceFlow.edit'],
        icon: record.flowStatus === 'DISABLED' ? <PlayOutline size={16} /> : <PauseOutline size={16} />,
        onClick: () => onStatusHandle(record),
      },
      {
        key: 'delete',
        label: formatMessage('common.delete'),
        auth: flowPerm(record).delete,
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
      }
      const api = isMark ? unmarkFlow : markFlow;
      return api?.(record.id).then(() => {
        message.success(formatMessage('common.optsuccess'));
        reload();
      });
    },
    renderPinIcon: (record: any) => record?.sort !== 1,
  };

  const listRowKey = readOnlyList ? flowListRowKey : 'id';

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
            name: 'groupId',
            hidden: true,
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
          auth={
            flowScope === 'all'
              ? [ButtonPermission['SourceFlow.add'], ButtonPermission['EventFlow.add']]
              : ButtonPermission['SourceFlow.add']
          }
          type="primary"
          onClick={() => {
            addGroupModalRef?.current?.onOpen?.(1);
          }}
        >
          <FolderAdd {...toolbarIconProps} />
          {formatMessage('common.newFolder')}
        </AuthButton>
      )}
      {!readOnlyList && !hideNewFlow && (
        <AuthButton
          auth={ButtonPermission['SourceFlow.add']}
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
      width: 280,
      maxWidth: 360,
      ellipsis: true,
      sorter: createFlowRowSorter('name'),
      sortKey: 'name',
      render: (text: any, item: any) => {
        const hasDesign = hasPermission(
          item.flowKind === 'event' ? ButtonPermission['EventFlow.design'] : ButtonPermission['SourceFlow.design']
        );
        return (
          <Flex gap={8} align="center" className="flow-list-name-cell">
            <FlowItemIcon category={item.category} flowKind={item.flowKind} size="sm" />
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
                    key={readOnlyList ? flowListRowKey(d) : d?.id}
                    iconBg={false}
                    classNames={{ headerTitle: 'flow-card-title' }}
                    header={{
                      customIcon: <FlowItemIcon category={d.category} flowKind={d.flowKind} size="md" />,
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
                          : hasPermission(
                                d.flowKind === 'event'
                                  ? ButtonPermission['EventFlow.design']
                                  : ButtonPermission['SourceFlow.design']
                              )
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
                    // actions={actions}
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
              onChange={onChange}
              style={{ height: embedded ? '100%' : undefined }}
              scroll={
                embedded
                  ? tableScrollY
                    ? { y: tableScrollY, x: 'max-content' }
                    : undefined
                  : { y: 'calc(100vh - 285px)', x: 'max-content' }
              }
              dataSource={data as any}
              rowKey={listRowKey}
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
        title={formatMessage(`collectionFlow.${isEdit}Flow`)}
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
          <PageTitleRow resourceKey="flow.collection.page">
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

export default CollectionFlow;
