import { useCallback, useEffect, useRef, useState } from 'react';
import { App, Button, Dropdown, Form, Input, Modal, Select, Space, Table, Tag, Tooltip } from 'antd';
import type { MenuProps } from 'antd';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import classNames from 'classnames';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router';
import {
  createFleetNode,
  deleteFleetNode,
  getFleetNodes,
  restoreFleetNode,
  updateFleetNode,
  type CreateFleetNodeResp,
  type FleetNode,
  type FleetNodeListResp,
  type FleetNodeStatus,
} from '@/apis/core-api/fleet';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton } from '@/components/auth';
import ComEmpty from '@/components/com-empty';
import { openConfirmModal } from '@/components/confirm-modal';
import HelpTooltip from '@/components/help-tooltip';
import { Add, Copy, Edit, OverflowMenuHorizontal, Renew, TrashCan } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import { hasPermission } from '@/utils/auth';
import { copyToClipboard } from '@/utils/common';
import { MAX_LENGTHS } from '@/utils/limits';
import styles from './index.module.scss';

type NodeFormValues = {
  name: string;
  description?: string;
};

type StatusFilter = 'all' | FleetNodeStatus;

const formatTimestamp = (value?: number) => (value ? dayjs(value).format('YYYY-MM-DD HH:mm:ss') : '-');

const statusColor: Record<FleetNodeStatus, string> = {
  online: 'success',
  offline: 'error',
  unjoined: 'warning',
  deleted: 'error',
};

const NodeIdCell = ({ value }: { value: string }) => {
  const ref = useRef<HTMLSpanElement>(null);
  const [tooltipTitle, setTooltipTitle] = useState<string>();

  return (
    <Tooltip title={tooltipTitle} placement="topLeft">
      <span
        ref={ref}
        className={styles.nodeIdCell}
        onMouseEnter={() => {
          const el = ref.current;
          setTooltipTitle(el && el.scrollWidth > el.clientWidth + 1 ? value : undefined);
        }}
      >
        {value}
      </span>
    </Tooltip>
  );
};

const FleetManagementPanel = () => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
  const [form] = Form.useForm<NodeFormValues>();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<FleetNode>();
  const [quotaExceeded, setQuotaExceeded] = useState(false);
  const [created, setCreated] = useState<CreateFleetNodeResp>();
  const [keyword, setKeyword] = useState('');
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('all');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [data, setData] = useState<FleetNodeListResp>();

  const loadNodes = useCallback(async () => {
    setLoading(true);
    try {
      setData(
        await getFleetNodes({
          keyword: keyword.trim() || undefined,
          status: statusFilter === 'all' ? undefined : statusFilter,
          page,
          size: pageSize,
        })
      );
    } catch {
      message.error(formatMessage('fleet.management.loadFailed'));
    } finally {
      setLoading(false);
    }
  }, [formatMessage, keyword, message, page, pageSize, statusFilter]);

  useEffect(() => {
    const timer = window.setTimeout(() => void loadNodes(), 0);
    return () => window.clearTimeout(timer);
  }, [loadNodes]);

  const summary = data?.summary;
  const activeNodes = summary?.active ?? 0;
  const nodeLimit = summary?.limit;
  const nodeLimitLabel = nodeLimit === -1 ? '∞' : (nodeLimit ?? '-');
  const quotaReached = nodeLimit != null && nodeLimit !== -1 && activeNodes >= nodeLimit;
  const activeQuota = `${activeNodes}/${nodeLimitLabel}`;

  const openCreate = () => {
    setEditing(undefined);
    setQuotaExceeded(quotaReached);
    form.resetFields();
    setModalOpen(true);
  };

  const openEdit = (node: FleetNode) => {
    setEditing(node);
    setQuotaExceeded(false);
    form.setFieldsValue({ name: node.name, description: node.description });
    setModalOpen(true);
  };

  const saveNode = async () => {
    if (!editing && quotaReached) {
      setQuotaExceeded(true);
      return;
    }
    const values = await form.validateFields();
    setSaving(true);
    try {
      if (editing) {
        await updateFleetNode(editing.fleetNodeID, values);
        message.success(formatMessage('fleet.management.updateSucceeded'));
      } else {
        setCreated(await createFleetNode(values));
        message.success(formatMessage('fleet.management.createSucceeded'));
      }
      setModalOpen(false);
      await loadNodes();
    } catch (error) {
      const responseMessage = String((error as { msg?: string })?.msg || '');
      if (!editing && responseMessage === 'FLEET_NODE_QUOTA_EXCEEDED') {
        setQuotaExceeded(true);
        await loadNodes();
      } else {
        message.error(formatMessage(editing ? 'fleet.management.updateFailed' : 'fleet.management.createFailed'));
      }
    } finally {
      setSaving(false);
    }
  };

  const confirmDelete = (node: FleetNode) => {
    openConfirmModal(modal, {
      title: formatMessage('fleet.management.deleteConfirmTitle'),
      content: formatMessage(
        node.status === 'online'
          ? 'fleet.management.deleteOnlineConfirmContent'
          : 'fleet.management.deleteConfirmContent',
        { name: node.name }
      ),
      okText: formatMessage('common.delete'),
      danger: true,
      cancelText: formatMessage('common.cancel'),
      onOk: async () => {
        try {
          await deleteFleetNode(node.fleetNodeID);
          message.success(formatMessage('fleet.management.deleteSucceeded'));
          await loadNodes();
        } catch {
          message.error(formatMessage('fleet.management.deleteFailed'));
        }
      },
    });
  };

  const confirmRestore = (node: FleetNode) => {
    openConfirmModal(modal, {
      title: formatMessage('fleet.management.restoreConfirmTitle'),
      content: formatMessage('fleet.management.restoreConfirmContent', { name: node.name }),
      okText: formatMessage('fleet.management.restore'),
      cancelText: formatMessage('common.cancel'),
      onOk: async () => {
        try {
          setCreated(await restoreFleetNode(node.fleetNodeID));
          message.success(formatMessage('fleet.management.restoreSucceeded'));
          await loadNodes();
        } catch {
          message.error(formatMessage('fleet.management.restoreFailed'));
        }
      },
    });
  };

  const columns: ColumnsType<FleetNode> = [
    {
      title: formatMessage('fleet.management.nodeName'),
      dataIndex: 'name',
      width: 120,
      ellipsis: true,
      render: (_, record) =>
        record.status === 'deleted' ? (
          <span className={styles.nodeNameText} title={record.name}>
            {record.name}
          </span>
        ) : (
          <button
            type="button"
            className={classNames(styles.nodeNameLink, 'table-nav-link')}
            title={record.name}
            onClick={() => {
              const tabRawName = String(record.name || '').trim();
              const tabName = tabRawName ? `Fleet · ${tabRawName}` : undefined;
              navigate(`/edge-connection/nodes/${record.fleetNodeID}`, {
                state: tabName ? { tabName, tabNameFull: tabName } : undefined,
              });
            }}
          >
            {record.name}
          </button>
        ),
    },
    {
      title: formatMessage('fleet.management.nodeId'),
      dataIndex: 'fleetNodeID',
      width: 180,
      ellipsis: { showTitle: false },
      render: (value: string) => <NodeIdCell value={value} />,
    },
    {
      title: formatMessage('common.status'),
      dataIndex: 'status',
      width: 110,
      render: (value: FleetNodeStatus) => (
        <Tag color={statusColor[value]}>{formatMessage(`fleet.status.${value}`)}</Tag>
      ),
    },
    {
      title: formatMessage('fleet.management.heartbeat'),
      dataIndex: 'lastOnlineAt',
      width: 180,
      render: formatTimestamp,
    },
    {
      title: formatMessage('fleet.management.lastSync'),
      dataIndex: 'lastSyncAt',
      width: 180,
      render: formatTimestamp,
    },
    {
      title: (
        <span className={styles.containerHeader}>
          {formatMessage('fleet.management.containerHealth')}
          <span className={styles.containerHeaderHint}>{formatMessage('fleet.management.containerHealthHint')}</span>
        </span>
      ),
      key: 'containerHealth',
      width: 130,
      onHeaderCell: () => ({ className: styles.containerHeaderCell }),
      render: (_: unknown, record: FleetNode) => {
        if (record.status !== 'online' || record.containerCount == null) {
          return <span className={styles.containerTag}>-/-</span>;
        }
        const unhealthy = record.unhealthyContainerCount ?? 0;
        const healthy = Math.max(record.containerCount - unhealthy, 0);
        return (
          <span className={styles.containerTag}>
            {healthy}/<span className={unhealthy > 0 ? styles.unhealthyCount : undefined}>{unhealthy}</span>
          </span>
        );
      },
    },
    {
      title: formatMessage('common.operation'),
      key: 'operation',
      width: 64,
      render: (_, record) => {
        const items: MenuProps['items'] = [];
        if (record.status === 'deleted') {
          if (hasPermission(ButtonPermission['FleetNode.restore'])) {
            items.push({
              key: 'restore',
              label: formatMessage('fleet.management.restore'),
              onClick: () => confirmRestore(record),
            });
          }
        } else {
          if (hasPermission(ButtonPermission['FleetNode.edit'])) {
            items.push({
              key: 'edit',
              label: formatMessage('common.edit'),
              icon: <Edit size={16} />,
              onClick: () => openEdit(record),
            });
          }
          if (hasPermission(ButtonPermission['FleetNode.delete'])) {
            items.push({
              key: 'delete',
              label: formatMessage('common.delete'),
              icon: <TrashCan size={16} />,
              danger: true,
              onClick: () => confirmDelete(record),
            });
          }
        }
        if (!items.length) return null;
        return (
          <div className="custom-operation">
            <Dropdown
              overlayClassName="pro-table-operation-menu"
              menu={{ items }}
              trigger={['click']}
              placement="bottomRight"
            >
              <Button type="text" icon={<OverflowMenuHorizontal size={16} />} />
            </Dropdown>
          </div>
        );
      },
    },
  ];

  const filterItems: Array<{ key: StatusFilter; label: string; count: number }> = [
    {
      key: 'all',
      label: formatMessage('fleet.management.filterAll'),
      count: summary?.active ?? 0,
    },
    {
      key: 'online',
      label: formatMessage('fleet.management.filterOnline'),
      count: summary?.online ?? 0,
    },
    {
      key: 'offline',
      label: formatMessage('fleet.management.filterOffline'),
      count: summary?.offline ?? 0,
    },
    {
      key: 'unjoined',
      label: formatMessage('fleet.management.filterUnjoined'),
      count: summary?.unjoined ?? 0,
    },
    {
      key: 'deleted',
      label: formatMessage('fleet.management.filterDeleted'),
      count: summary?.deleted ?? 0,
    },
  ];

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current || 1);
    setPageSize(pagination.pageSize || 20);
  };

  return (
    <div className={styles.management}>
      <div className={styles.summaryRow}>
        <div className={styles.summaryCard}>
          <div className={styles.summaryLabelRow}>
            <span className={styles.summaryLabel}>{formatMessage('fleet.management.activeQuota')}</span>
            <HelpTooltip
              title={formatMessage('fleet.management.quotaSummary', {
                active: activeNodes,
                available: nodeLimitLabel,
              })}
            />
          </div>
          <strong className={styles.summaryValue}>{activeQuota}</strong>
        </div>
        <div className={styles.summaryCardGroup}>
          <div className={styles.summaryMetric}>
            <span className={styles.summaryLabel}>{formatMessage('fleet.management.offlineCount')}</span>
            <strong className={styles.summaryValue}>{summary?.offline ?? 0}</strong>
          </div>
          <div className={styles.summaryMetric}>
            <span className={styles.summaryLabel}>{formatMessage('fleet.management.abnormalContainerNodes')}</span>
            <strong className={styles.summaryValue}>{summary?.onlineWithAbnormalContainers ?? 0}</strong>
          </div>
        </div>
      </div>

      <div className={styles.filterToolbar}>
        <Select<StatusFilter>
          className={styles.statusFilter}
          aria-label={formatMessage('common.status')}
          value={statusFilter}
          options={filterItems.map((item) => ({
            value: item.key,
            label: `${item.label} (${item.count})`,
          }))}
          onChange={(value) => {
            setPage(1);
            setStatusFilter(value);
          }}
        />
        <div className={styles.toolbarActions}>
          <Input.Search
            allowClear
            className={styles.search}
            placeholder={formatMessage('fleet.management.searchPlaceholder')}
            onSearch={(value) => {
              setPage(1);
              setKeyword(value);
            }}
          />
          <Tooltip title={formatMessage('common.refresh')}>
            <Button
              className={styles.refreshBtn}
              aria-label={formatMessage('common.refresh')}
              icon={<Renew size={16} />}
              onClick={() => void loadNodes()}
            />
          </Tooltip>
          <Tooltip title={quotaReached ? formatMessage('fleet.management.quotaExceeded') : undefined}>
            <span className={styles.addNodeButtonWrap}>
              <AuthButton
                auth={ButtonPermission['FleetNode.add']}
                type="primary"
                icon={<Add size={16} />}
                disabled={quotaReached}
                onClick={openCreate}
              >
                {formatMessage('fleet.management.addNode')}
              </AuthButton>
            </span>
          </Tooltip>
        </div>
      </div>

      <Table<FleetNode>
        rowKey="fleetNodeID"
        loading={loading}
        columns={columns}
        dataSource={data?.list || []}
        locale={{ emptyText: <ComEmpty description={formatMessage('fleet.management.empty')} /> }}
        pagination={{ current: page, pageSize, total: data?.total || 0, showSizeChanger: true }}
        onChange={onTableChange}
        tableLayout="fixed"
      />

      <Modal
        open={modalOpen}
        title={formatMessage(editing ? 'fleet.management.editNode' : 'fleet.management.addNode')}
        confirmLoading={saving}
        onOk={() => void saveNode()}
        onCancel={() => {
          setModalOpen(false);
          setQuotaExceeded(false);
        }}
        destroyOnHidden
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label={formatMessage('fleet.management.nodeName')}
            rules={[
              { required: true, message: formatMessage('fleet.management.nodeNameRequired') },
              {
                max: MAX_LENGTHS.connectionName,
                message: formatMessage('fleet.management.nameTooLong', {
                  length: MAX_LENGTHS.connectionName,
                }),
              },
            ]}
          >
            <Input maxLength={MAX_LENGTHS.connectionName} showCount />
          </Form.Item>
          <Form.Item
            name="description"
            label={formatMessage('common.description')}
            rules={[
              {
                max: MAX_LENGTHS.description,
                message: formatMessage('fleet.management.descTooLong', {
                  length: MAX_LENGTHS.description,
                }),
              },
            ]}
          >
            <Input.TextArea maxLength={MAX_LENGTHS.description} rows={4} showCount />
          </Form.Item>
          {!editing && quotaExceeded ? (
            <p className={styles.quotaError} role="alert">
              {formatMessage('fleet.management.quotaExceeded')}
            </p>
          ) : null}
        </Form>
      </Modal>

      <Modal
        open={Boolean(created)}
        title={formatMessage('fleet.management.tokenTitle')}
        width={640}
        footer={
          <Button type="primary" onClick={() => setCreated(undefined)}>
            {formatMessage('common.confirm')}
          </Button>
        }
        closable={false}
        maskClosable={false}
        destroyOnHidden
      >
        <p className={styles.tokenHint}>{formatMessage('fleet.management.tokenHint')}</p>
        <Space.Compact className={styles.tokenCompact}>
          <Input.TextArea
            className={styles.tokenArea}
            value={created?.accessToken || ''}
            readOnly
            autoSize={{ minRows: 4, maxRows: 8 }}
          />
          <Button
            className={styles.tokenCopyBtn}
            icon={<Copy size={16} />}
            aria-label={formatMessage('common.copy')}
            onClick={() => {
              copyToClipboard(created?.accessToken || '', (success) => {
                if (success) {
                  message.success(formatMessage('common.copySuccess'));
                } else {
                  message.error(formatMessage('common.copyFail'));
                }
              });
            }}
          />
        </Space.Compact>
      </Modal>
    </div>
  );
};

export default FleetManagementPanel;
