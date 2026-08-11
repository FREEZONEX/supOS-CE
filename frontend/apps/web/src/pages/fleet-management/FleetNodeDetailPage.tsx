import { CircleHelp, User, Users } from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  Alert,
  App,
  Button,
  Descriptions,
  Drawer,
  Modal,
  Select,
  Space,
  Table,
  Tabs,
  Tag,
  Tooltip,
  Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { useNavigate, useParams } from 'react-router';
import {
  deleteFleetNode,
  getFleetAuditLogs,
  getFleetNode,
  getFleetNodeAccounts,
  getFleetNodeContainer,
  getFleetNodeContainerLogs,
  getFleetNodeContainers,
  getFleetNodeUNSScopes,
  getFleetOperation,
  replaceFleetNodeAccounts,
  restartFleetNodeContainer,
  type FleetAssignedAccount,
  type FleetContainerDetail,
  type FleetContainerSummary,
  type FleetNode,
  type FleetNodeAccountAssignments,
  type FleetNodeStatus,
  type FleetRemoteContainerSnapshot,
  type FleetOperationStatus,
  type FleetUNSScopeState,
} from '@/apis/core-api/fleet';
import { getAuditLogDetail, type AuditLogItem } from '@/apis/core-api/audit-log';
import { getUserManageList } from '@/apis/core-api/user-manage';
import { getTreeData } from '@/apis/core-api/uns';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton } from '@/components/auth';
import ComBackButton from '@/components/com-back-button';
import ComEmpty from '@/components/com-empty';
import ComLayout from '@/components/com-layout';
import HelpTooltip from '@/components/help-tooltip';
import { openConfirmModal } from '@/components/confirm-modal';
import { Close, ContainerServices, Renew, ServerDns, TrashCan } from '@/components/lucide-icon/carbon';
import { useTabName, useTranslate } from '@/hooks';
import { AuditLogDetailDrawer } from '@/pages/audit-log/AuditLogDetailDrawer';
import { FleetContainerLogsDrawer } from '@/pages/fleet-shared/FleetContainerLogsDrawer';
import { FleetContainerTable, fleetContainerSorters } from '@/pages/fleet-shared/FleetContainerTable';
import { useFleetContainerLogs } from '@/pages/fleet-shared/useFleetContainerLogs';
import type { UnsTreeNode } from '@/pages/uns/types';
import { hasPermission } from '@/utils/auth';
import FleetUNSSyncBrowser from './components/FleetUNSSyncBrowser';
import styles from './FleetNodeDetailPage.module.scss';

const NODE_TAB_PREFIX = 'Fleet · ';
const NODE_TAB_MAX_NAME_LENGTH = 18;
const truncateNodeTabName = (value: string) =>
  value.length > NODE_TAB_MAX_NAME_LENGTH ? `${value.slice(0, NODE_TAB_MAX_NAME_LENGTH)}...` : value;

const formatTimestamp = (value?: number | string) => {
  if (!value) return '-';
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss');
};

const statusColor: Record<FleetNodeStatus, string> = {
  online: 'success',
  offline: 'error',
  unjoined: 'warning',
  deleted: 'error',
};

const formatUptime = (startedAt?: number) => {
  if (!startedAt) return '-';
  const seconds = Math.max(0, dayjs().diff(dayjs(startedAt), 'second'));
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return [days ? `${days}d` : '', hours ? `${hours}h` : '', `${minutes}m`].filter(Boolean).join(' ');
};

const shortID = (value?: string) => {
  const normalized = String(value || '').trim();
  if (normalized.length <= 16) return normalized || '-';
  return `${normalized.slice(0, 8)}…${normalized.slice(-6)}`;
};

const parseImageRef = (image?: string) => {
  const raw = String(image || '').trim();
  if (!raw) return { name: '-', tag: '' };

  const digestAt = raw.lastIndexOf('@');
  if (digestAt > 0) {
    return { name: raw.slice(0, digestAt), tag: raw.slice(digestAt + 1) };
  }

  const slash = raw.lastIndexOf('/');
  const colon = raw.lastIndexOf(':');
  if (colon > slash) {
    return { name: raw.slice(0, colon), tag: raw.slice(colon + 1) };
  }

  return { name: raw, tag: '' };
};

type FleetAuditDetail = {
  action?: string;
  fleetNodeID?: string;
  containerID?: string;
  accountCount?: number;
};

const parseFleetAuditDetail = (value?: string): FleetAuditDetail => {
  try {
    return value ? JSON.parse(value) : {};
  } catch {
    return {};
  }
};

const containerStateColor = (state?: string) => {
  switch (String(state || '').toLowerCase()) {
    case 'running':
      return 'success';
    case 'restarting':
      return 'processing';
    case 'stopped':
    case 'exited':
      return 'default';
    default:
      return 'warning';
  }
};

const terminalOperationStatuses = new Set<FleetOperationStatus>(['canceled', 'failed', 'succeeded', 'timedOut']);

const wait = (milliseconds: number) => new Promise((resolve) => window.setTimeout(resolve, milliseconds));

const FleetNodeDetailPage = () => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
  const { nodeKey = '' } = useParams();
  const nodeID = decodeURIComponent(nodeKey);
  const [loading, setLoading] = useState(false);
  const [node, setNode] = useState<FleetNode>();
  const [scopes, setScopes] = useState<FleetUNSScopeState[]>([]);
  const [unsTree, setUnsTree] = useState<UnsTreeNode[]>([]);
  const [accounts, setAccounts] = useState<FleetNodeAccountAssignments>();
  const [snapshot, setSnapshot] = useState<FleetRemoteContainerSnapshot>();
  const [auditLogs, setAuditLogs] = useState<AuditLogItem[]>([]);
  const [auditDetail, setAuditDetail] = useState<AuditLogItem>();
  const [auditDetailOpen, setAuditDetailOpen] = useState(false);
  const [auditDetailLoading, setAuditDetailLoading] = useState(false);
  const [containerDetail, setContainerDetail] = useState<FleetContainerDetail>();
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [restartingID, setRestartingID] = useState('');
  const [logsDrawerHeight, setLogsDrawerHeight] = useState(0);
  const [accountsOpen, setAccountsOpen] = useState(false);
  const [accountOptions, setAccountOptions] = useState<Array<{ label: string; value: string }>>([]);
  const [selectedAccountIDs, setSelectedAccountIDs] = useState<string[]>([]);
  const [savingAccounts, setSavingAccounts] = useState(false);

  const tabRawName = String(node?.name || '').trim();
  useTabName(tabRawName, {
    displayName: tabRawName ? `${NODE_TAB_PREFIX}${truncateNodeTabName(tabRawName)}` : undefined,
    fullName: tabRawName ? `${NODE_TAB_PREFIX}${tabRawName}` : undefined,
  });

  const loadDetail = useCallback(
    async (refreshContainers = false) => {
      if (!nodeID) return;
      setLoading(true);
      try {
        const nodeData = await getFleetNode(nodeID);
        const [scopeData, accountData, treeData, containerData, auditData] = await Promise.all([
          getFleetNodeUNSScopes(nodeID),
          getFleetNodeAccounts(nodeID),
          getTreeData(),
          nodeData.status === 'unjoined'
            ? Promise.resolve(undefined)
            : getFleetNodeContainers(nodeID, refreshContainers).catch(() => undefined),
          getFleetAuditLogs(nodeID, { page: 1, size: 100 }),
        ]);
        setNode(nodeData);
        setScopes(scopeData.list || []);
        setUnsTree(treeData || []);
        setAccounts(accountData);
        setSnapshot(containerData);
        setAuditLogs(auditData.list || []);
      } catch {
        message.error(formatMessage('fleet.detail.loadFailed'));
      } finally {
        setLoading(false);
      }
    },
    [formatMessage, message, nodeID]
  );

  useEffect(() => {
    const timer = window.setTimeout(() => void loadDetail(false), 0);
    return () => window.clearTimeout(timer);
  }, [loadDetail]);

  const openContainerDetail = useCallback(
    async (container: FleetContainerSummary) => {
      setDetailOpen(true);
      setDetailLoading(true);
      try {
        setContainerDetail(await getFleetNodeContainer(nodeID, container.containerID));
      } catch {
        setContainerDetail(undefined);
        message.error(formatMessage('fleet.detail.containerLoadFailed'));
      } finally {
        setDetailLoading(false);
      }
    },
    [formatMessage, message, nodeID]
  );

  const fetchRemoteContainerLogs = useCallback(
    (container: FleetContainerSummary, request: Parameters<typeof getFleetNodeContainerLogs>[2]) => {
      return getFleetNodeContainerLogs(nodeID, container.containerID, request);
    },
    [nodeID]
  );
  const handleLogsError = useCallback(
    () => message.error(formatMessage('fleet.runtime.logsFailed')),
    [formatMessage, message]
  );
  const {
    open: logsOpen,
    loading: logsLoading,
    downloading: logsDownloading,
    logs,
    container: logsContainer,
    live: logsLive,
    follow: logsFollow,
    filterText: logsFilterText,
    tabs: logTabs,
    activeContainerID,
    logsRef,
    openLogs,
    closeLogs,
    refreshLogs,
    downloadLogs,
    clearLogs,
    closeLogTab,
    reorderTabs,
    setActiveTab,
    setLive: setLogsLive,
    setFollow: setLogsFollow,
    setFilterText: setLogsFilterText,
    handleScroll: handleLogsScroll,
  } = useFleetContainerLogs({ fetchLogs: fetchRemoteContainerLogs, onError: handleLogsError });

  const confirmRestart = useCallback(
    (container: FleetContainerSummary) => {
      openConfirmModal(modal, {
        title: formatMessage('fleet.runtime.restartConfirmTitle'),
        content: formatMessage('fleet.runtime.restartConfirmContent', {
          name: container.displayName || container.serviceName,
        }),
        okText: formatMessage('fleet.runtime.restart'),
        cancelText: formatMessage('common.cancel'),
        onOk: async () => {
          setRestartingID(container.containerID);
          try {
            let operation = await restartFleetNodeContainer(nodeID, container.containerID);
            const deadline = Date.now() + 135_000;
            while (!terminalOperationStatuses.has(operation.status) && Date.now() < deadline) {
              await wait(1000);
              operation = await getFleetOperation(nodeID, operation.operationID);
            }
            if (operation.status !== 'succeeded') {
              throw new Error(operation.errorCode || operation.status);
            }
            message.success(formatMessage('fleet.runtime.restartSucceeded'));
            await loadDetail(true);
          } catch {
            message.error(formatMessage('fleet.runtime.restartFailed'));
          } finally {
            setRestartingID('');
          }
        },
      });
    },
    [formatMessage, loadDetail, message, modal, nodeID]
  );

  const confirmDelete = () => {
    if (!node) return;
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
          navigate('/edge-connection?tab=edgeNode');
        } catch {
          message.error(formatMessage('fleet.management.deleteFailed'));
        }
      },
    });
  };

  const openAccountSettings = async () => {
    try {
      const users = await getUserManageList({ pageNo: 1, pageSize: 1000 });
      setAccountOptions(
        (users.data || [])
          .filter((user: any) => {
            const userID = Number(user.userId || user.id);
            const userName = String(user.preferredUsername || '')
              .trim()
              .toLowerCase();
            const roleCodes = (user.roleList || []).map((role: any) =>
              String(role.roleCode || role.code || '')
                .trim()
                .toLowerCase()
            );
            return userID >= 1_000_000_000 && userName !== 'tier0' && !roleCodes.includes('owner');
          })
          .map((user: any) => ({
            value: String(user.userId || user.id),
            label:
              user.preferredUsername || user.firstName || user.email || formatMessage('fleet.detail.unnamedAccount'),
          }))
      );
      setSelectedAccountIDs(accounts?.centerUserIDs || []);
      setAccountsOpen(true);
    } catch {
      message.error(formatMessage('fleet.detail.accountOptionsFailed'));
    }
  };

  const saveAccounts = async () => {
    if (!accounts) return;
    setSavingAccounts(true);
    try {
      const result = await replaceFleetNodeAccounts(nodeID, {
        centerUserIDs: selectedAccountIDs,
        expectedRevision: accounts.revision,
      });
      setAccounts(result);
      setAccountsOpen(false);
      message.success(formatMessage('fleet.detail.accountsSaved'));
    } catch {
      message.error(formatMessage('fleet.detail.accountsSaveFailed'));
      await loadDetail();
    } finally {
      setSavingAccounts(false);
    }
  };

  const accountColumns = useMemo<ColumnsType<FleetAssignedAccount>>(
    () => [
      {
        title: formatMessage('fleet.detail.accountName'),
        dataIndex: 'nickName',
        render: (value, row) => value || row.userName,
      },
      { title: formatMessage('fleet.detail.userName'), dataIndex: 'userName' },
      {
        title: formatMessage('fleet.detail.role'),
        dataIndex: 'roleCode',
        render: (value?: string) => {
          const role = String(value || '').trim();
          if (!role) return '-';
          const normalized = role.toLowerCase();
          const isAdmin = normalized === 'admin';
          const label = role.charAt(0).toUpperCase() + role.slice(1).toLowerCase();
          return (
            <span className={isAdmin ? styles.roleTagAdmin : styles.roleTagOperator}>
              <User size={14} strokeWidth={1.75} className={styles.roleTagIcon} />
              <span>{label}</span>
            </span>
          );
        },
      },
      {
        title: formatMessage('common.status'),
        dataIndex: 'userStatus',
        render: (value: number) => (
          <span className={value === 0 ? styles.statusTagMuted : styles.statusTagActive}>
            {formatMessage(value === 0 ? 'fleet.account.disabled' : 'fleet.account.active')}
          </span>
        ),
      },
      { title: formatMessage('fleet.detail.lastSync'), dataIndex: 'assignedAt', render: formatTimestamp },
    ],
    [formatMessage]
  );

  const containerColumns = useMemo<ColumnsType<FleetContainerSummary>>(
    () => [
      {
        title: <span className={styles.tableHeaderSingle}>{formatMessage('fleet.runtime.service')}</span>,
        dataIndex: 'displayName',
        width: 168,
        ellipsis: true,
        sorter: fleetContainerSorters.service,
        onHeaderCell: () => ({ className: styles.tableHeaderCellSingle }),
        render: (_, record) => (
          <Button
            type="link"
            className={`table-link-button table-link-button-neutral ${styles.serviceLink}`}
            onClick={() => void openContainerDetail(record)}
          >
            {record.displayName || record.serviceName || shortID(record.containerID)}
          </Button>
        ),
      },
      {
        title: <span className={styles.tableHeaderSingle}>{formatMessage('common.status')}</span>,
        dataIndex: 'state',
        width: 156,
        sorter: fleetContainerSorters.state,
        onHeaderCell: () => ({ className: styles.tableHeaderCellSingle }),
        render: (value: string) => <Tag color={containerStateColor(value)}>{String(value || '-').toLowerCase()}</Tag>,
      },
      {
        title: <span className={styles.tableHeaderSingle}>{formatMessage('fleet.runtime.uptime')}</span>,
        dataIndex: 'startedAt',
        width: 148,
        sorter: fleetContainerSorters.startedAt,
        onHeaderCell: () => ({ className: styles.tableHeaderCellSingle }),
        render: formatUptime,
      },
      {
        title: <span className={styles.tableHeaderSingle}>{formatMessage('fleet.runtime.restartCount')}</span>,
        dataIndex: 'restartCount',
        width: 148,
        sorter: fleetContainerSorters.restartCount,
        onHeaderCell: () => ({ className: styles.tableHeaderCellSingle }),
      },
      {
        title: <span className={styles.tableHeaderSingle}>{formatMessage('fleet.runtime.image')}</span>,
        dataIndex: 'image',
        width: 240,
        ellipsis: true,
        sorter: fleetContainerSorters.image,
        onHeaderCell: () => ({ className: styles.tableHeaderCellSingle }),
        render: (value: string) => {
          const { name, tag } = parseImageRef(value);
          return (
            <Tooltip title={value || undefined}>
              <div className={styles.imageCell}>
                <span className={styles.imageName}>{name}</span>
                {tag ? <span className={styles.imageTag}>{tag}</span> : null}
              </div>
            </Tooltip>
          );
        },
      },
      {
        title: <span className={styles.operationHeader}>{formatMessage('common.operation')}</span>,
        key: 'operation',
        width: 148,
        fixed: 'right',
        onHeaderCell: () => ({ className: styles.operationHeaderCell }),
        render: (_, record) => {
          const canLogs = hasPermission(ButtonPermission['FleetContainer.logs']);
          const canRestart = hasPermission(ButtonPermission['FleetContainer.restart']);
          if (!canLogs && !canRestart) return null;
          return (
            <div className={styles.containerOps}>
              {canLogs ? (
                <Button
                  type="link"
                  className={styles.opLink}
                  disabled={node?.status !== 'online'}
                  onClick={() => void openLogs(record)}
                >
                  {formatMessage('fleet.runtime.logs')}
                </Button>
              ) : null}
              {canRestart ? (
                <Tooltip title={!record.restartable ? formatMessage('fleet.runtime.restartProtected') : undefined}>
                  <Button
                    type="link"
                    danger
                    className={styles.opLink}
                    disabled={node?.status !== 'online' || !record.restartable || Boolean(restartingID)}
                    loading={restartingID === record.containerID}
                    onClick={() => confirmRestart(record)}
                  >
                    {formatMessage('fleet.runtime.restart')}
                  </Button>
                </Tooltip>
              ) : null}
            </div>
          );
        },
      },
    ],
    [confirmRestart, formatMessage, node?.status, openContainerDetail, openLogs, restartingID]
  );

  const containerNames = useMemo(
    () =>
      new Map(
        (snapshot?.containers || []).map((container) => [
          container.containerID,
          container.displayName || container.serviceName || shortID(container.containerID),
        ])
      ),
    [snapshot?.containers]
  );

  const openAuditDetail = useCallback(
    async (record: AuditLogItem) => {
      setAuditDetailOpen(true);
      setAuditDetailLoading(true);
      try {
        setAuditDetail(await getAuditLogDetail(record.id));
      } catch {
        setAuditDetail(undefined);
        message.error(formatMessage('auditLog.detailFailed'));
      } finally {
        setAuditDetailLoading(false);
      }
    },
    [formatMessage, message]
  );

  const operationColumns = useMemo<ColumnsType<AuditLogItem>>(
    () => [
      {
        title: formatMessage('fleet.detail.operationType'),
        key: 'action',
        width: 220,
        render: (_, record) => {
          const detail = parseFleetAuditDetail(record.detailJson);
          const action =
            detail.action ||
            (detail.containerID && record.businessType === 'Start'
              ? 'container.restart'
              : detail.containerID && record.businessType === 'Export'
                ? 'container.logs'
                : 'node.update');
          return formatMessage(`fleet.audit.action.${action}`);
        },
      },
      {
        title: formatMessage('fleet.detail.operationTarget'),
        key: 'target',
        width: 160,
        ellipsis: { showTitle: false },
        render: (_, record) => {
          const detail = parseFleetAuditDetail(record.detailJson);
          const target = detail.containerID
            ? containerNames.get(detail.containerID) || shortID(detail.containerID)
            : node?.name || '-';
          return <Tooltip title={target}>{target}</Tooltip>;
        },
      },
      {
        title: formatMessage('common.status'),
        dataIndex: 'code',
        width: 130,
        render: (value: number) => (
          <Tag color={value >= 200 && value < 400 ? 'success' : 'error'}>
            {formatMessage(value >= 200 && value < 400 ? 'fleet.operation.accepted' : 'fleet.operation.failed')}
          </Tag>
        ),
      },
      {
        title: formatMessage('fleet.detail.operator'),
        key: 'operator',
        width: 180,
        render: (_, record) => record.operatorName || record.operatorEmail || '-',
      },
      {
        title: formatMessage('fleet.detail.occurredAt'),
        dataIndex: 'createdAt',
        width: 180,
        render: formatTimestamp,
      },
    ],
    [containerNames, formatMessage, node?.name]
  );

  const containerHealthy = snapshot
    ? Math.max((snapshot.summary.total ?? 0) - (snapshot.summary.unhealthy ?? 0), 0)
    : null;
  const containerUnhealthy = snapshot?.summary.unhealthy ?? null;

  const overview = (
    <div className={styles.detailGrid}>
      <div className={styles.detailGridTop}>
        <section className={styles.detailCard}>
          <div className={styles.detailCardTitle}>
            <ServerDns size={18} />
            <span>{formatMessage('fleet.detail.hardware')}</span>
          </div>
          <div className={styles.detailFields}>
            <div className={styles.detailField}>
              <span className={styles.detailFieldLabel}>{formatMessage('fleet.detail.osArch')}</span>
              <span className={styles.detailFieldValue}>{node ? `${node.os || '-'} / ${node.arch || '-'}` : '-'}</span>
            </div>
            <div className={styles.detailField}>
              <span className={styles.detailFieldLabel}>{formatMessage('fleet.detail.lastReported')}</span>
              <span className={styles.detailFieldValue}>{formatTimestamp(node?.lastOnlineAt)}</span>
            </div>
          </div>
          {node?.status === 'offline' ? (
            <p className={styles.staleHint}>{formatMessage('fleet.detail.offlineLastValueHint')}</p>
          ) : null}
        </section>
        <section className={styles.detailCard}>
          <div className={styles.detailCardTitle}>
            <ContainerServices size={18} />
            <span>{formatMessage('fleet.detail.hostDocker')}</span>
          </div>
          <div className={styles.detailFields}>
            <div className={styles.detailField}>
              <span className={styles.detailFieldLabel}>{formatMessage('fleet.detail.agentVersion')}</span>
              <span className={styles.detailFieldValue}>{node?.agentVersion || '-'}</span>
            </div>
            <div className={styles.detailField}>
              <span className={styles.detailFieldLabel}>{formatMessage('fleet.detail.protocolVersion')}</span>
              <span className={styles.detailFieldValue}>{node?.protocolVersion || '-'}</span>
            </div>
            <div className={styles.detailField}>
              <span className={styles.detailFieldLabel}>{formatMessage('fleet.detail.dockerVersion')}</span>
              <span className={styles.detailFieldValue}>{snapshot?.engine?.version || '-'}</span>
            </div>
          </div>
        </section>
        <section className={`${styles.detailCard} ${styles.detailCardMetric}`}>
          <div className={styles.detailCardTitle}>
            <ContainerServices size={18} />
            <span>{formatMessage('fleet.detail.containers')}</span>
          </div>
          <span className={styles.detailMetricLabel}>{formatMessage('fleet.detail.containersHealthyUnhealthy')}</span>
          <div className={styles.containersMetric}>
            {containerHealthy == null ? (
              '--'
            ) : (
              <>
                <span>{containerHealthy}</span>
                <span className={styles.metricDivider}>/</span>
                <span className={styles.unhealthyCount}>{containerUnhealthy ?? 0}</span>
              </>
            )}
          </div>
          {!snapshot ? (
            <p className={styles.staleHint}>{formatMessage('fleet.detail.remoteContainersUnavailable')}</p>
          ) : null}
        </section>
      </div>
    </div>
  );

  const tabs = [
    { key: 'overview', label: formatMessage('fleet.detail.overview'), children: overview },
    {
      key: 'containers',
      label: formatMessage('fleet.detail.containers'),
      children: (
        <div className={styles.containersPanel}>
          {snapshot?.stale ? (
            <Alert type="error" showIcon message={formatMessage('fleet.detail.snapshotStale')} />
          ) : null}
          <div className={styles.containersSummary}>
            <div className={styles.containersSummaryItem}>
              <span className={styles.containersSummaryLabel}>{formatMessage('fleet.detail.containers')}</span>
              <strong className={styles.containersSummaryValue}>{snapshot?.summary.total ?? '-'}</strong>
            </div>
            <div className={styles.containersSummaryItem}>
              <span className={styles.containersSummaryLabel}>{formatMessage('fleet.runtime.stateRunning')}</span>
              <strong className={styles.containersSummaryValue}>{snapshot?.summary.running ?? '-'}</strong>
            </div>
            <div className={styles.containersSummaryItem}>
              <span className={styles.containersSummaryLabel}>{formatMessage('fleet.detail.unhealthyContainers')}</span>
              <strong className={`${styles.containersSummaryValue} ${styles.unhealthyCount}`}>
                {snapshot?.summary.unhealthy ?? '-'}
              </strong>
            </div>
          </div>
          <FleetContainerTable
            columns={containerColumns}
            containers={snapshot?.containers}
            emptyText={<ComEmpty description={formatMessage('fleet.detail.remoteContainersUnavailable')} />}
          />
        </div>
      ),
    },
    {
      key: 'uns',
      label: formatMessage('cluster.unsUplink'),
      children: <FleetUNSSyncBrowser scopes={scopes} treeData={unsTree} />,
    },
    {
      key: 'accounts',
      label: formatMessage('cluster.userDownlink'),
      children: (
        <section className={styles.accountCard}>
          <div className={styles.sectionHead}>
            <div className={styles.sectionHeader}>
              <div className={styles.sectionTitle}>
                <Users size={16} strokeWidth={1.75} className={styles.sectionTitleIcon} />
                <span>{formatMessage('fleet.detail.accountsSection')}</span>
                <HelpTooltip
                  title={formatMessage('fleet.detail.accountSyncHint')}
                  icon={<CircleHelp size={16} strokeWidth={1.75} className={styles.sectionHelpIcon} />}
                />
              </div>
              <div className={styles.sectionActions}>
                <AuthButton
                  auth={ButtonPermission['FleetNode.accounts']}
                  type="primary"
                  onClick={() => void openAccountSettings()}
                >
                  {formatMessage('fleet.detail.configureAccounts')}
                </AuthButton>
              </div>
            </div>
            <div className={styles.sectionDividerLine} aria-hidden />
          </div>
          <Table<FleetAssignedAccount>
            className={styles.accountTable}
            rowKey="centerUserID"
            columns={accountColumns}
            dataSource={accounts?.accounts || []}
            pagination={false}
          />
        </section>
      ),
    },
    {
      key: 'operations',
      label: formatMessage('fleet.detail.operations'),
      children: (
        <Table<AuditLogItem>
          className={styles.operationsTable}
          rowKey="id"
          columns={operationColumns}
          dataSource={auditLogs}
          pagination={false}
          locale={{ emptyText: <ComEmpty description={formatMessage('fleet.detail.operationsEmpty')} /> }}
          scroll={{ x: 880 }}
          onRow={(record) => ({
            onClick: () => void openAuditDetail(record),
          })}
        />
      ),
    },
  ];

  return (
    <ComLayout loading={loading} className={styles.detailPage}>
      <div className={styles.content} style={{ paddingBottom: logsDrawerHeight ? logsDrawerHeight + 16 : 0 }}>
        <div className={styles.detailHeader}>
          <div className={styles.detailHeaderLeft}>
            <ComBackButton onClick={() => navigate('/edge-connection?tab=edgeNode')} />
            <div className={styles.detailHeaderInfo}>
              <div className={styles.detailTitleRow}>
                <Typography.Title level={3} className={styles.detailTitle}>
                  {node?.name || nodeID}
                </Typography.Title>
                {node ? (
                  <Tag className={styles.statusTag} color={statusColor[node.status]}>
                    {formatMessage(`fleet.status.${node.status}`)}
                  </Tag>
                ) : null}
              </div>
              <div className={styles.meta}>
                <Tooltip title={node?.fleetNodeID || nodeID}>
                  <span className={styles.metaChip}>
                    {formatMessage('fleet.detail.nodeID', { id: node?.fleetNodeID || nodeID })}
                  </span>
                </Tooltip>
                <span className={styles.metaChip}>
                  {formatMessage('fleet.detail.nodeType', { type: node?.centerType || '-' })}
                </span>
                <span className={styles.metaChip}>
                  {formatMessage('fleet.detail.heartbeatAt', { time: formatTimestamp(node?.lastOnlineAt) })}
                </span>
              </div>
            </div>
          </div>
          <Space size={8} className={styles.detailHeaderActions}>
            <Tooltip title={formatMessage('common.refresh')}>
              <Button
                className={styles.refreshBtn}
                aria-label={formatMessage('common.refresh')}
                icon={<Renew size={16} />}
                onClick={() => void loadDetail(true)}
              />
            </Tooltip>
            <Tooltip title={formatMessage('common.delete')}>
              <AuthButton
                auth={ButtonPermission['FleetNode.delete']}
                className={styles.refreshBtn}
                danger
                disabled={!node}
                aria-label={formatMessage('common.delete')}
                icon={<TrashCan size={16} />}
                onClick={confirmDelete}
              />
            </Tooltip>
          </Space>
        </div>

        <Tabs items={tabs} tabBarGutter={24} />
      </div>

      <AuditLogDetailDrawer
        open={auditDetailOpen}
        loading={auditDetailLoading}
        detail={auditDetail}
        onClose={() => setAuditDetailOpen(false)}
      />

      <Modal
        open={accountsOpen}
        title={formatMessage('fleet.detail.configureAccounts')}
        confirmLoading={savingAccounts}
        onOk={() => void saveAccounts()}
        onCancel={() => setAccountsOpen(false)}
      >
        <Typography.Paragraph>{formatMessage('fleet.detail.configureAccountsHint')}</Typography.Paragraph>
        <Select
          mode="multiple"
          showSearch
          optionFilterProp="label"
          value={selectedAccountIDs}
          options={accountOptions}
          onChange={setSelectedAccountIDs}
          className={styles.fullWidth}
          placeholder={formatMessage('fleet.detail.selectAccounts')}
        />
      </Modal>

      <Drawer
        rootClassName={styles.containerDetailDrawer}
        open={detailOpen}
        width={720}
        loading={detailLoading}
        title={formatMessage('fleet.detail.containerDetails')}
        closable={false}
        destroyOnClose
        onClose={() => setDetailOpen(false)}
        extra={
          <Tooltip title={formatMessage('common.close')}>
            <Button color="default" variant="text" onClick={() => setDetailOpen(false)} icon={<Close size={20} />} />
          </Tooltip>
        }
        classNames={{
          body: styles.containerDetailDrawerBody,
        }}
      >
        <Descriptions
          className={styles.containerDetailDescriptions}
          column={1}
          size="small"
          bordered
          labelStyle={{ whiteSpace: 'nowrap' }}
        >
          <Descriptions.Item label={formatMessage('fleet.runtime.service')}>
            {containerDetail?.displayName || containerDetail?.serviceName || '-'}
          </Descriptions.Item>
          <Descriptions.Item label={formatMessage('fleet.runtime.image')}>
            {containerDetail?.image || '-'}
          </Descriptions.Item>
          <Descriptions.Item label={formatMessage('fleet.detail.imageID')}>
            {containerDetail?.imageID || '-'}
          </Descriptions.Item>
          <Descriptions.Item label={formatMessage('common.status')}>
            {containerDetail ? (
              <Tag color={containerStateColor(containerDetail.state)}>{containerDetail.state}</Tag>
            ) : (
              '-'
            )}
          </Descriptions.Item>
          <Descriptions.Item label={formatMessage('fleet.detail.health')}>
            {containerDetail?.healthStatus || containerDetail?.health?.status || '-'}
          </Descriptions.Item>
          <Descriptions.Item label={formatMessage('fleet.detail.createdAt')}>
            {formatTimestamp(containerDetail?.createdAt)}
          </Descriptions.Item>
          <Descriptions.Item label={formatMessage('fleet.runtime.restartCount')}>
            {containerDetail?.restartCount ?? '-'}
          </Descriptions.Item>
        </Descriptions>
      </Drawer>

      <FleetContainerLogsDrawer
        open={logsOpen}
        loading={logsLoading}
        downloading={logsDownloading}
        logs={logs}
        container={logsContainer}
        live={logsLive}
        follow={logsFollow}
        filterText={logsFilterText}
        tabs={logTabs}
        activeContainerID={activeContainerID}
        logsRef={logsRef}
        closeLogs={closeLogs}
        refreshLogs={refreshLogs}
        downloadLogs={downloadLogs}
        clearLogs={clearLogs}
        closeLogTab={closeLogTab}
        reorderTabs={reorderTabs}
        setActiveTab={setActiveTab}
        setLive={setLogsLive}
        setFollow={setLogsFollow}
        setFilterText={setLogsFilterText}
        handleScroll={handleLogsScroll}
        onHeightChange={setLogsDrawerHeight}
      />
    </ComLayout>
  );
};

export default FleetNodeDetailPage;
