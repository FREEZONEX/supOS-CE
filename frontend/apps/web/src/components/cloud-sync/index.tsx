import { type ReactNode, useCallback, useEffect, useState } from 'react';
import { App, Button, Drawer, Flex, Form, Input, Space, Switch, Tag, Typography } from 'antd';
import { Help } from '@carbon/icons-react';
import {
  Add,
  CheckmarkFilled,
  ChevronLeft,
  Copy,
  Download,
  Renew,
  ServerDns,
  Signal,
  Time,
  Upload,
  User,
} from '@/components/lucide-icon/carbon';
import { PageTitleIcon } from '@/components/lucide-icon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { ComLayout, ProTable } from '@/components';
import ComBackButton from '@/components/com-back-button';
import ComButton from '@/components/com-button';
import ComEmpty from '@/components/com-empty';
import HelpTooltip from '@/components/help-tooltip';
import OperationForm from '@/components/operation-form';
import ProModal from '@/components/pro-modal';
import { AuthButton } from '@/components/auth';
import { ButtonPermission } from '@/common-types/button-permission';
import { coreApi } from '@/apis/core-api/core-adapter';
import {
  createMqttAuthClient,
  deleteMqttAuthClient,
  disableMqttAuthClient,
  enableMqttAuthClient,
  listMqttAuthClients,
  resetMqttAuthPassword,
  type MqttAuthClient,
} from '@/apis/core-api/mqtt-auth';
import { getHomeClusterConfig, type HomeClusterConfig } from '@/apis/core-api/home';
import {
  getCloudSyncConfig,
  getClusterSyncLogs,
  userSyncCloudSync,
  type CloudSyncConfigResp,
  type CloudSyncLogItem,
} from '@/apis/core-api/cloudsync';
import CloudSyncConnectModal from '@/components/cloud-sync/CloudSyncConnectModal';
import { copyToClipboard } from '@/utils/common';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import { formatTimestamp } from '@/utils';
import { useTranslate } from '@/hooks';
import { useBaseStore } from '@/stores/base';
import { MAX_LENGTHS } from '@/utils/limits';
import classNames from 'classnames';
import styles from './index.module.scss';
import {
  EDGE_NODE_PREVIEW_BRANCH,
  EDGE_NODE_PREVIEW_BRANCH_OFFLINE,
  EDGE_NODE_PREVIEW_DETAIL,
  EDGE_NODE_PREVIEW_HOST,
  parseEdgeNodePreviewMode,
  PREVIEW_CLOUD_SYNC_CONNECTED,
  PREVIEW_CLOUD_SYNC_DISCONNECTED,
  PREVIEW_EDGE_NODES,
  PREVIEW_LAST_UNS_SYNC_LOG,
  PREVIEW_LAST_USER_SYNC_LOG,
  PREVIEW_SYNC_LOGS,
  PREVIEW_USER_SYNC_LOGS,
} from '@/pages/mqtt-auth/edge-node-preview';
import { useLocation, useNavigate, useParams, useSearchParams } from 'react-router';

type FormValues = {
  name: string;
  description?: string;
  clientIdRandomSuffixEnabled?: boolean;
};

type TokenFormValues = {
  tokenName: string;
  description?: string;
};

const displayBaseURL = (value?: string) => {
  const raw = String(value || '').trim();
  if (!raw) return '';
  try {
    const url = new URL(raw);
    url.pathname = '';
    url.search = '';
    url.hash = '';
    return url.toString().replace(/\/$/, '');
  } catch {
    return raw.replace(/\/api\/.*$/, '').replace(/\/$/, '');
  }
};

const userSyncUserSectionByStat: Record<string, { dataKey: string; titleKey: string }> = {
  created: { dataKey: 'createdUsers', titleKey: 'cluster.userSyncCreatedUsers' },
  updated: { dataKey: 'updatedUsers', titleKey: 'cluster.userSyncUpdatedUsers' },
  disabled: { dataKey: 'disabledUsers', titleKey: 'cluster.userSyncDisabledUsers' },
  deleted: { dataKey: 'deletedUsers', titleKey: 'cluster.userSyncDeletedUsers' },
  skipped: { dataKey: 'skippedUsers', titleKey: 'cluster.userSyncSkippedUsers' },
  conflicts: { dataKey: 'conflicts', titleKey: 'cluster.userSyncConflicts' },
};

const unsSyncStatKeys = ['total', 'created', 'updated', 'deleted', 'failed'];
const unsSyncNodeStatKeys = ['created', 'updated', 'deleted', 'failed'];
const defaultUnsSyncNodeDetailLimit = 100;
const defaultClusterSyncTopicRoot = 'cloudsync';
const cloudSyncMetadataAPIPath = '/api/core/uns/cloud-edge-sync/meta';

type UnsSyncNodeRow = {
  key: string;
  label: string;
  topic: unknown;
  detail: Record<string, unknown>;
  danger: boolean;
};

type ClusterNode = {
  id: number;
  nodeKey?: string;
  nodeName?: string;
  nodeID?: string;
  name?: string;
  status?: string | number;
  description?: string;
  lastSeenTime?: number;
  lastSeen?: number;
  connectClientID?: number;
  lastError?: string;
  token?: string;
  syncUserName?: string;
  syncUserEmail?: string;
  dispatchUserName?: string;
  dispatchUserEmail?: string;
};

const hiddenTokenMetadataKeys = new Set(['password', 'signedPassword']);

const decodeBase64Url = (value: string) => {
  const normalized = value.replace(/-/g, '+').replace(/_/g, '/');
  const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=');
  const binary = window.atob(padded);
  const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));
  return new TextDecoder().decode(bytes);
};

const parseTokenMetadata = (token?: string) => {
  const payload = String(token || '')
    .split('#')[1]
    ?.trim();
  if (!payload) return {};
  try {
    const parsed = JSON.parse(decodeBase64Url(payload));
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return {};
    return Object.fromEntries(
      Object.entries(parsed as Record<string, unknown>).filter(([key, value]) => {
        if (hiddenTokenMetadataKeys.has(key)) return false;
        return value !== undefined && value !== null && value !== '';
      })
    );
  } catch {
    return {};
  }
};

const visibleTokenMetadataEntries = (metadata: Record<string, unknown>) =>
  [
    ['username', metadata.username],
    ['clientID', metadata.clientID],
    ['mqttURL', metadata.mqttURL || metadata.mqttBrokers],
    ['httpBaseURL', metadata.httpBaseURL],
  ] as const;

const MQTT_ACCESS_TAB = 'mqttAccess';
const EDGE_NODE_TAB = 'edgeNode';

const canManageEdgeNodes = (isHost?: boolean) => Boolean(isHost);

const statusMessageKey = (status?: string | number) => {
  const normalized = String(status || '').toLowerCase();
  if (normalized === '1') return 'mqttAuth.connected';
  if (normalized === 'online' || normalized === 'connected') return 'mqttAuth.connected';
  if (normalized === 'connecting') return 'home.statusConnecting';
  if (normalized === 'error') return 'home.statusError';
  return 'mqttAuth.disconnected';
};

const nodeKeyOf = (record?: ClusterNode) => record?.nodeKey || record?.nodeID || String(record?.id || '');
const nodeNameOf = (record?: ClusterNode) => record?.nodeName || record?.name || nodeKeyOf(record);
const nodeStatusText = (status?: string | number) => {
  const normalized = String(status || '').toLowerCase();
  if (normalized === '1' || normalized === 'online' || normalized === 'connected') return 'online';
  if (normalized === 'connecting') return 'connecting';
  if (normalized === 'error') return 'error';
  return 'offline';
};
const nodeLastSeenOf = (record?: ClusterNode) => record?.lastSeenTime || record?.lastSeen || 0;

const parseLogObject = (value?: string) => {
  const text = String(value || '').trim();
  if (!text) return {};
  try {
    const parsed = JSON.parse(text);
    return parsed && typeof parsed === 'object' && !Array.isArray(parsed) ? parsed : { value: parsed };
  } catch {
    return { value: text };
  }
};

const logFieldFallbacks: Record<string, string> = {
  role: '节点角色',
  nodeKey: '节点标识',
  nodeName: '节点名称',
  workspaceID: '工作空间',
  mountMode: '挂载模式',
  syncMetricEnabled: '在中心节点存储数据',
  transport: '传输方式',
  version: '版本',
  name: '名称',
  nodeID: '节点 ID',
  clientID: 'Client ID',
  username: 'Username',
  httpBaseURL: 'HTTP Base URL',
  mqttURL: 'MQTT URL',
  mqttBrokers: 'MQTT Brokers',
  syncMode: '同步模式',
  topicRoot: 'Topic Root',
  force: '强制同步',
  endpoint: 'HTTP 端点',
  source: '来源',
  operatorId: '操作人',
  durationMs: '耗时',
  stats: '统计',
  createdUsers: '创建用户',
  updatedUsers: '更新用户',
  disabledUsers: '停用用户',
  deletedUsers: '删除用户',
  skippedUsers: '跳过用户',
  roleMappings: '角色映射',
  passwordPolicy: '密码策略',
  userId: '用户 ID',
  userName: '用户名',
  email: '邮箱',
  status: '状态',
  detail: '详情',
  targetRoleName: '目标角色',
  sourceRoleName: '来源角色',
  fallback: '回退方式',
  httpError: 'HTTP 错误',
  topic: 'Topic',
  sourceID: '源节点 ID',
  errmsg: '错误信息',
  sourcePath: '源路径',
  targetPath: '目标路径',
  nodeKind: '节点类型',
  alias: '别名',
  bytes: '数据大小',
  nodeCount: '节点数',
  rootCount: '根节点数',
  rootTargetPaths: '目标根路径',
  reason: '原因',
  receivedNodes: '接收节点数',
  desiredNodes: '目标节点数',
  importFolders: '导入文件夹',
  importFiles: '导入文件',
  existingNodes: '已有节点',
  missingFolders: '新增文件夹',
  missingFiles: '新增文件',
  failedNodes: '失败 Topic',
  conflictNodes: '冲突 Topic',
  skippedNodeItems: '跳过 Topic',
  blockedTargets: '阻塞 Topic',
  conflictCount: '冲突数',
  skippedNodes: '跳过节点',
  affectedTopics: '受影响 Topics',
  totalTopics: '总数',
  failedCount: '失败',
  createdFiles: '创建文件',
  createdNodes: '创建 Topic',
  updatedFiles: '更新文件',
  updatedNodes: '更新 Topic',
  deletedFiles: '删除文件',
  deletedNodes: '删除 Topic',
  createdFolders: '创建目录',
  snapshotVersion: '快照版本',
  conflicts: '冲突 Topic',
  filteredDesiredNodes: '过滤后节点数',
  filteredExpectedNamespaces: '过滤后命名空间数',
  snapshotWorkspaceID: '快照工作空间',
  snapshotAuthID: '快照连接 ID',
  value: '原始信息',
};

const hiddenErrorItemKeys = new Set(['sourcePath', 'targetPath', 'namespace', 'subtree']);
const errorListLogKeys = ['failedNodes', 'conflictNodes', 'skippedNodeItems'];
const errorScalarLogKeys = new Set(['errmsg', 'httpError']);
const nestedLogFieldPriority = ['topic', 'errmsg', 'alias', 'reason'];

const formatBytes = (value: unknown) => {
  const bytes = Number(value);
  if (!Number.isFinite(bytes)) return String(value ?? '-');
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
};

const syncLogKeyOf = (record: any) => String(record.id || `${record.createdTime}-${record.snapshotVersion}`);

const MqttAuthPage = () => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const location = useLocation();
  const { nodeKey } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const systemCluster = useBaseStore((state) => state.systemInfo?.cluster);
  const { modal, message } = App.useApp();
  const [form] = Form.useForm<FormValues>();
  const [tokenForm] = Form.useForm<TokenFormValues>();
  const isEdgeConnectionRoute = location.pathname.startsWith('/edge-connection');
  const [activeTab, setActiveTab] = useState(() =>
    searchParams.get('tab') === EDGE_NODE_TAB ? EDGE_NODE_TAB : MQTT_ACCESS_TAB
  );
  const [dataSource, setDataSource] = useState<MqttAuthClient[]>([]);
  const [edgeNodes, setEdgeNodes] = useState<ClusterNode[]>([]);
  const [loading, setLoading] = useState(false);
  const [edgeLoading, setEdgeLoading] = useState(false);
  const [cloudSyncConfig, setCloudSyncConfig] = useState<CloudSyncConfigResp | null>(null);
  const [cloudSyncConfigLoading, setCloudSyncConfigLoading] = useState(false);
  const [cloudSyncUserSyncing, setCloudSyncUserSyncing] = useState(false);
  const [cloudSyncConnectOpen, setCloudSyncConnectOpen] = useState(false);
  const [syncLogDrawerOpen, setSyncLogDrawerOpen] = useState(false);
  const [syncLogNode, setSyncLogNode] = useState<ClusterNode | null>(null);
  const [syncLogs, setSyncLogs] = useState<CloudSyncLogItem[]>([]);
  const [syncLogLoading, setSyncLogLoading] = useState(false);
  const [syncLogPagination, setSyncLogPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [focusedSyncLog, setFocusedSyncLog] = useState<CloudSyncLogItem | null>(null);
  const [lastUnsSyncLog, setLastUnsSyncLog] = useState<CloudSyncLogItem | null>(null);
  const [lastUserSyncLog, setLastUserSyncLog] = useState<CloudSyncLogItem | null>(null);
  const [userSyncDrawerOpen, setUserSyncDrawerOpen] = useState(false);
  const [userSyncLogs, setUserSyncLogs] = useState<CloudSyncLogItem[]>([]);
  const [userSyncLoading, setUserSyncLoading] = useState(false);
  const [userSyncPagination, setUserSyncPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [focusedUserSyncLog, setFocusedUserSyncLog] = useState<CloudSyncLogItem | null>(null);
  const [syncLogDetailStat, setSyncLogDetailStat] = useState<Record<string, string>>({});
  const [userSyncDetailStat, setUserSyncDetailStat] = useState<Record<string, string>>({});
  const [clusterConfig, setClusterConfig] = useState<HomeClusterConfig>();
  const [clusterConfigLoaded, setClusterConfigLoaded] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);
  const decodedNodeKey = nodeKey ? decodeURIComponent(nodeKey) : '';
  const edgeNodePreviewMode = import.meta.env.DEV
    ? parseEdgeNodePreviewMode(searchParams.get('edgeNodePreview'))
    : null;
  const isEdgeNodePreview = Boolean(edgeNodePreviewMode);
  const disablePreviewMutations = isEdgeNodePreview;
  const previewHost = edgeNodePreviewMode === EDGE_NODE_PREVIEW_HOST;
  const previewBranch = edgeNodePreviewMode === EDGE_NODE_PREVIEW_BRANCH;
  const previewBranchOffline = edgeNodePreviewMode === EDGE_NODE_PREVIEW_BRANCH_OFFLINE;
  const previewDetail = edgeNodePreviewMode === EDGE_NODE_PREVIEW_DETAIL;
  const activeCloudSyncConfig = previewBranch
    ? PREVIEW_CLOUD_SYNC_CONNECTED
    : previewBranchOffline
      ? PREVIEW_CLOUD_SYNC_DISCONNECTED
      : cloudSyncConfig;
  const activeLastUnsSyncLog = previewBranch || previewBranchOffline ? PREVIEW_LAST_UNS_SYNC_LOG : lastUnsSyncLog;
  const activeLastUserSyncLog = previewBranch || previewBranchOffline ? PREVIEW_LAST_USER_SYNC_LOG : lastUserSyncLog;
  const activeEdgeNodeRows = previewHost || previewDetail ? PREVIEW_EDGE_NODES : edgeNodes;
  const edgeNodeRows = activeEdgeNodeRows;
  const edgeNodeRequested = activeTab === EDGE_NODE_TAB || Boolean(decodedNodeKey);
  const clusterHostFlag = clusterConfig?.isHost ?? systemCluster?.isHost;
  const clusterConfigResolved = clusterConfigLoaded && clusterHostFlag !== undefined;
  const clusterIsHost = Boolean(clusterHostFlag);
  const edgeNodeManageAllowed = clusterConfigResolved && canManageEdgeNodes(clusterIsHost);
  const showEdgeNodeTab = true;
  const canConnectCloudSync = clusterConfigResolved && !clusterIsHost;
  const showEdgeNodeHost = previewHost || (!isEdgeNodePreview && edgeNodeManageAllowed);
  const showEdgeNodeBranch = previewBranch || previewBranchOffline || (!isEdgeNodePreview && canConnectCloudSync);

  useEffect(() => {
    let cancelled = false;
    void getHomeClusterConfig().then((config) => {
      if (!cancelled) {
        setClusterConfig(config);
        setClusterConfigLoaded(true);
      }
    });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    const tab = searchParams.get('tab');
    if (tab === EDGE_NODE_TAB || tab === MQTT_ACCESS_TAB) {
      setActiveTab(tab);
    }
  }, [searchParams]);

  const setTab = useCallback(
    (nextTab: string) => {
      setActiveTab(nextTab);
      const next = new URLSearchParams(searchParams);
      next.set('tab', nextTab);
      setSearchParams(next, { replace: true });
    },
    [searchParams, setSearchParams]
  );

  useEffect(() => {
    const tab = searchParams.get('tab');
    if (tab === MQTT_ACCESS_TAB || tab === EDGE_NODE_TAB || decodedNodeKey) return;
    if (isEdgeConnectionRoute) return;
    if (edgeNodeManageAllowed) {
      setActiveTab(EDGE_NODE_TAB);
    }
  }, [decodedNodeKey, edgeNodeManageAllowed, isEdgeConnectionRoute, searchParams]);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const resp = await listMqttAuthClients();
      setDataSource(resp.data || []);
    } finally {
      setLoading(false);
    }
  }, []);

  const loadEdgeNodes = useCallback(async () => {
    if (!edgeNodeManageAllowed) {
      setEdgeNodes([]);
      return;
    }
    setEdgeLoading(true);
    try {
      const resp = await coreApi.get('/cluster/nodes');
      setEdgeNodes(resp?.list || []);
    } finally {
      setEdgeLoading(false);
    }
  }, [edgeNodeManageAllowed]);

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    if (isEdgeNodePreview) return;
    if (!edgeNodeRequested) return;
    if (!clusterConfigResolved) return;
    if (edgeNodeManageAllowed) {
      void loadEdgeNodes();
      return;
    }
    setEdgeNodes([]);
  }, [clusterConfigResolved, edgeNodeManageAllowed, edgeNodeRequested, isEdgeNodePreview, loadEdgeNodes]);

  const copyValue = (value?: string) => {
    copyToClipboard(value || '', (success) => {
      if (success) {
        message.success(formatMessage('common.copySuccess'));
      } else {
        message.error(formatMessage('common.copyFail'));
      }
    });
  };

  const refreshCloudSyncConfig = useCallback(async () => {
    setCloudSyncConfigLoading(true);
    try {
      const config = await getCloudSyncConfig();
      setCloudSyncConfig(config);
      return config;
    } catch {
      setCloudSyncConfig(null);
      return null;
    } finally {
      setCloudSyncConfigLoading(false);
    }
  }, []);

  useEffect(() => {
    if (isEdgeNodePreview) return;
    if (!edgeNodeRequested || !canConnectCloudSync) return;
    void refreshCloudSyncConfig();
  }, [canConnectCloudSync, edgeNodeRequested, isEdgeNodePreview, refreshCloudSyncConfig]);

  const loadBranchLatestSyncLogs = useCallback(async () => {
    const [unsResp, userResp] = await Promise.all([
      getClusterSyncLogs({ syncType: 'uns', page: 1, size: 1 }),
      getClusterSyncLogs({ syncType: 'user', page: 1, size: 1 }),
    ]);
    setLastUnsSyncLog(unsResp.list?.[0] || null);
    setLastUserSyncLog(userResp.list?.[0] || null);
  }, []);

  const loadUserSyncLogs = useCallback(
    async (current = 1, pageSize = 10) => {
      setUserSyncLoading(true);
      try {
        if (isEdgeNodePreview) {
          setUserSyncLogs(PREVIEW_USER_SYNC_LOGS);
          setUserSyncPagination({
            current,
            pageSize,
            total: PREVIEW_USER_SYNC_LOGS.length,
          });
          return;
        }
        const resp = await getClusterSyncLogs({
          syncType: 'user',
          page: current,
          size: pageSize,
        });
        setUserSyncLogs(resp.list || []);
        setUserSyncPagination({
          current,
          pageSize,
          total: Number(resp.total || 0),
        });
        if (current === 1) {
          setLastUserSyncLog(resp.list?.[0] || null);
        }
      } finally {
        setUserSyncLoading(false);
      }
    },
    [isEdgeNodePreview]
  );

  const loadSyncLogs = useCallback(
    async (node: ClusterNode | null, current = 1, pageSize = 10) => {
      setSyncLogLoading(true);
      try {
        if (isEdgeNodePreview) {
          setSyncLogs(PREVIEW_SYNC_LOGS);
          setSyncLogPagination({
            current,
            pageSize,
            total: PREVIEW_SYNC_LOGS.length,
          });
          return;
        }
        const params: Record<string, string | number | undefined> = {
          syncType: 'uns',
          page: current,
          size: pageSize,
        };
        if (node) {
          if (node.connectClientID) {
            params.connectClientKey = String(node.connectClientID);
          } else {
            params.nodeID = nodeKeyOf(node);
          }
        }
        const resp = await getClusterSyncLogs(params);
        setSyncLogs(resp.list || []);
        setSyncLogPagination({
          current,
          pageSize,
          total: Number(resp.total || 0),
        });
      } finally {
        setSyncLogLoading(false);
      }
    },
    [isEdgeNodePreview]
  );

  useEffect(() => {
    if (isEdgeNodePreview) return;
    if (!edgeNodeRequested || !canConnectCloudSync) return;
    void loadBranchLatestSyncLogs();
  }, [canConnectCloudSync, edgeNodeRequested, isEdgeNodePreview, loadBranchLatestSyncLogs]);

  const credentialModal = (record: MqttAuthClient) => {
    modal.info({
      title: formatMessage('mqttAuth.credentials'),
      width: 620,
      icon: null,
      content: (
        <Space direction="vertical" style={{ width: '100%', marginTop: 12 }}>
          <Space.Compact style={{ width: '100%' }}>
            <Input value={record.clientId} readOnly addonBefore="Client ID" />
            <Button icon={<Copy size={16} />} onClick={() => copyValue(record.clientId)} />
          </Space.Compact>
          {record.clientIdRandomSuffixEnabled ? (
            <Typography.Text type="secondary" style={{ fontSize: 12 }}>
              {formatMessage('mqttAuth.randomSuffixHint', { clientID: '{clientID}' })}
            </Typography.Text>
          ) : null}
          {[
            [formatMessage('mqttAuth.username'), record.username],
            [formatMessage('mqttAuth.password'), record.password],
          ].map(([label, value]) => (
            <Space.Compact key={label} style={{ width: '100%' }}>
              <Input value={value} readOnly addonBefore={label} />
              <Button icon={<Copy size={16} />} onClick={() => copyValue(value)} />
            </Space.Compact>
          ))}
        </Space>
      ),
    });
  };

  const tokenModal = (token: string) => {
    modal.info({
      title: formatMessage('mqttAuth.edgeTokenCreated'),
      width: 620,
      icon: <CheckmarkFilled size={22} fill="var(--supos-theme-color)" />,
      content: (
        <Space.Compact style={{ width: '100%', marginTop: 12 }}>
          <Input value={token} readOnly />
          <Button icon={<Copy size={16} />} onClick={() => copyValue(token)} />
        </Space.Compact>
      ),
    });
  };

  const closeCreate = () => {
    form.resetFields();
    setCreateOpen(false);
  };

  const onCreateSave = async () => {
    const values = await form.validateFields();
    const record = await createMqttAuthClient({
      name: values.name?.trim(),
      description: values.description,
      clientIdRandomSuffixEnabled: !!values.clientIdRandomSuffixEnabled,
    });
    message.success(formatMessage('mqttAuth.generated'));
    await load();
    closeCreate();
    credentialModal(record);
  };

  const openCreate = () => {
    form.setFieldsValue({ name: '', description: '', clientIdRandomSuffixEnabled: false });
    setCreateOpen(true);
  };

  const openCreateEdgeToken = () => {
    tokenForm.setFieldsValue({ tokenName: '', description: '' });
    modal.confirm({
      title: formatMessage('common.create'),
      width: 620,
      icon: null,
      className: styles.confirmModal,
      content: (
        <Form form={tokenForm} layout="vertical" style={{ marginTop: 16 }} preserve={false}>
          <Form.Item
            name="tokenName"
            label={formatMessage('mqttAuth.name')}
            rules={[
              { required: true, whitespace: true, message: formatMessage('mqttAuth.nameRequired') },
              { max: MAX_LENGTHS.connectionName, message: formatMessage('mqttAuth.nameMaxLength') },
            ]}
          >
            <Input maxLength={MAX_LENGTHS.connectionName} placeholder="edge-token" />
          </Form.Item>
          <Form.Item name="description" label={formatMessage('common.description')}>
            <Input.TextArea
              rows={3}
              maxLength={MAX_LENGTHS.description}
              placeholder={formatMessage('route.optional')}
            />
          </Form.Item>
        </Form>
      ),
      onOk: async () => {
        const values = await tokenForm.validateFields();
        const resp = await coreApi.post('/cluster/tokens', {
          tokenName: values.tokenName?.trim(),
          description: values.description || '',
        });
        message.success(formatMessage('mqttAuth.edgeTokenCreated'));
        tokenModal(resp?.token || '');
        await loadEdgeNodes();
      },
    });
  };

  const openCloudSyncConnect = useCallback(async () => {
    setActiveTab(EDGE_NODE_TAB);
    setCloudSyncConnectOpen(true);
    await refreshCloudSyncConfig();
  }, [refreshCloudSyncConfig]);

  const openUserSyncHistory = () => {
    setFocusedUserSyncLog(null);
    setUserSyncDrawerOpen(true);
    void loadUserSyncLogs(1, userSyncPagination.pageSize);
  };

  const openSyncedAccountList = () => {
    if (!activeLastUserSyncLog) {
      message.info(formatMessage('common.noData', {}, 'No data'));
      return;
    }
    setFocusedUserSyncLog(activeLastUserSyncLog);
    setUserSyncDrawerOpen(true);
  };

  const handleCloudSyncUserSync = async () => {
    if (!cloudSyncConfig?.configured) {
      message.warning(
        formatMessage('cluster.userSyncNotConfigured', {}, 'Configure Edge Connection before pulling users.')
      );
      return;
    }
    setCloudSyncUserSyncing(true);
    try {
      const log = await userSyncCloudSync();
      setLastUserSyncLog(log);
      if (userSyncDrawerOpen) {
        await loadUserSyncLogs(1, userSyncPagination.pageSize);
      }
      const status = String(log?.status || '').toLowerCase();
      if (status === 'failed') {
        message.error(log?.errorMessage || formatMessage('cluster.userSyncFailed', {}, 'User sync failed.'));
      } else if (status === 'partial') {
        message.warning(formatMessage('cluster.userSyncPartial', {}, 'User sync completed with skips or conflicts.'));
      } else {
        message.success(formatMessage('cluster.userSyncSuccess', {}, 'User sync completed.'));
      }
    } catch {
      message.error(formatMessage('cluster.userSyncFailed', {}, 'User sync failed.'));
    } finally {
      setCloudSyncUserSyncing(false);
    }
  };

  const resetPassword = (record: MqttAuthClient) => {
    modal.confirm({
      title: formatMessage('mqttAuth.resetPasswordTitle'),
      icon: null,
      className: styles.confirmModal,
      content: formatMessage('mqttAuth.resetPasswordContent'),
      onOk: async () => {
        const updated = await resetMqttAuthPassword(record.id);
        message.success(formatMessage('mqttAuth.passwordReset'));
        await load();
        credentialModal(updated);
      },
    });
  };

  const toggleAuth = async (record: MqttAuthClient) => {
    if (record.authStatus === 1) {
      await disableMqttAuthClient(record.id);
      message.success(formatMessage('OpenData.disabledSuccessfully'));
    } else {
      await enableMqttAuthClient(record.id);
      message.success(formatMessage('OpenData.startSuccessfully'));
    }
    await load();
  };

  const remove = (record: MqttAuthClient) => {
    modal.confirm({
      ...createDeleteConfirmOptions({
        title: formatMessage('mqttAuth.deleteTitle'),
        name: record.name || record.clientId,
        formatMessage,
        okText: formatMessage('common.confirm'),
        cancelText: formatMessage('common.cancel'),
      }),
      onOk: async () => {
        await deleteMqttAuthClient(record.id);
        message.success(formatMessage('OpenData.deleteSuccessfully'));
        await load();
      },
    });
  };

  const renderConnectionStatus = (record: MqttAuthClient) => {
    const connectedTime = record.connectedTime || 0;
    const connected = record.connectStatus === 1 && connectedTime > 0;
    return (
      <Space direction="vertical" size={2}>
        <Tag className={connected ? styles.successTag : undefined}>
          {connected ? formatMessage('mqttAuth.connected') : formatMessage('mqttAuth.disconnected')}
        </Tag>
        {connected ? (
          <span className={styles.connectionMeta}>
            {[record.ip, formatTimestamp(connectedTime)].filter(Boolean).join(' / ')}
          </span>
        ) : null}
      </Space>
    );
  };

  const credentialOperationItems = (record: MqttAuthClient): any[] => {
    const togglePermission =
      record.authStatus === 1 ? ButtonPermission['MqttAuth.disable'] : ButtonPermission['MqttAuth.enable'];
    return [
      {
        key: 'reset',
        label: formatMessage('common.reset'),
        auth: ButtonPermission['MqttAuth.reset'],
        icon: <Renew size={16} />,
        onClick: () => resetPassword(record),
      },
      {
        key: 'toggle',
        label: record.authStatus === 1 ? formatMessage('common.disable') : formatMessage('common.enable'),
        auth: togglePermission,
        onClick: () => toggleAuth(record),
      },
      { key: 'deleteDivider', type: 'divider' },
      {
        key: 'delete',
        label: formatMessage('common.delete'),
        auth: ButtonPermission['MqttAuth.delete'],
        danger: true,
        onClick: () => remove(record),
      },
    ];
  };

  const columns: any = [
    {
      title: formatMessage('mqttAuth.name'),
      dataIndex: 'name',
      width: 180,
      ellipsis: true,
      render: (value: string) => value || '-',
    },
    {
      title: formatMessage('common.description'),
      dataIndex: 'description',
      width: 280,
      maxWidth: 400,
      ellipsis: true,
      render: (value: string) => value || '-',
    },
    {
      title: formatMessage('mqttAuth.clientId', {}, 'Client ID'),
      dataIndex: 'clientId',
      width: 230,
      render: (value: string) => (
        <Space>
          <span className={styles.credentialValue} title={value}>
            {value}
          </span>
          <Button type="text" size="small" icon={<Copy size={14} />} onClick={() => copyValue(value)} />
        </Space>
      ),
    },
    {
      title: formatMessage('mqttAuth.connectionStatus'),
      dataIndex: 'connectedTime',
      width: 140,
      maxWidth: 180,
      render: (_: number, record: MqttAuthClient) => renderConnectionStatus(record),
    },
    {
      title: formatMessage('common.status'),
      dataIndex: 'authStatus',
      width: 100,
      render: (value: number) =>
        value === 1 ? (
          <Tag className={styles.successTag}>{formatMessage('common.enable')}</Tag>
        ) : (
          <Tag>{formatMessage('common.disable')}</Tag>
        ),
    },
    {
      title: formatMessage('OpenData.creationTime'),
      dataIndex: 'createdTime',
      width: 180,
      maxWidth: 200,
      ellipsis: true,
      render: (value: number) => {
        const text = value ? formatTimestamp(value) : '-';
        return (
          <span className={styles.cellEllipsis} title={text !== '-' ? text : undefined}>
            {text}
          </span>
        );
      },
    },
  ];

  const gotoNodeDetail = useCallback(
    (record: ClusterNode) => {
      const key = nodeKeyOf(record) || String(record.id || '');
      if (!key) return;
      const next = new URLSearchParams(searchParams);
      next.set('tab', EDGE_NODE_TAB);
      if (isEdgeNodePreview) {
        next.set('edgeNodePreview', EDGE_NODE_PREVIEW_DETAIL);
      } else {
        next.delete('edgeNodePreview');
      }
      navigate(`/edge-connection/nodes/${encodeURIComponent(key)}?${next.toString()}`);
    },
    [isEdgeNodePreview, navigate, searchParams]
  );

  const syncStatusTag = (status?: string) => {
    const normalized = String(status || '').toLowerCase();
    const tagClass = styles.syncStatusTag;
    if (normalized === 'success')
      return <Tag className={classNames(tagClass, styles.successTag)}>{formatMessage('mqttAuth.syncSuccess')}</Tag>;
    if (normalized === 'partial')
      return <Tag className={classNames(tagClass, styles.warningTag)}>{formatMessage('mqttAuth.syncPartial')}</Tag>;
    if (normalized === 'failed')
      return <Tag className={classNames(tagClass, styles.errorTag)}>{formatMessage('mqttAuth.syncFailed')}</Tag>;
    if (normalized === 'skipped')
      return <Tag className={classNames(tagClass, styles.neutralTag)}>{formatMessage('mqttAuth.syncSkipped')}</Tag>;
    return <Tag className={classNames(tagClass, styles.neutralTag)}>{formatMessage('mqttAuth.syncUnknown')}</Tag>;
  };

  const syncModeText = (syncMode?: string) => {
    if (syncMode === 'cloudEdge') return formatMessage('mqttAuth.syncModeCloudEdge');
    if (syncMode === 'edgeEdge') return formatMessage('mqttAuth.syncModeEdgeEdge');
    return syncMode || '-';
  };

  const logFieldLabel = (key: string) =>
    formatMessage(`mqttAuth.syncLogField.${key}`, undefined, logFieldFallbacks[key] || key);

  const logNumberValue = (value: unknown) => {
    if (Array.isArray(value)) return value.length;
    const numeric = Number(value);
    return Number.isFinite(numeric) ? numeric : 0;
  };

  const syncLogDetailsData = (record: CloudSyncLogItem) => {
    const data = parseLogObject(record.details) as Record<string, unknown>;
    if (record.snapshotVersion) {
      data.snapshotVersion = record.snapshotVersion;
    }
    return data;
  };

  const dataHasField = (data: Record<string, unknown>, key: string) => Object.prototype.hasOwnProperty.call(data, key);

  const cudCountsFromData = (data: Record<string, unknown>) => ({
    created: logNumberValue(data.createdFiles),
    updated: logNumberValue(data.updatedFiles),
    deleted: logNumberValue(data.deletedFiles),
  });

  const cudTotalFromData = (data: Record<string, unknown>) => {
    const counts = cudCountsFromData(data);
    return counts.created + counts.updated + counts.deleted;
  };

  const dataHasCudCounts = (data: Record<string, unknown>) =>
    dataHasField(data, 'createdFiles') || dataHasField(data, 'updatedFiles') || dataHasField(data, 'deletedFiles');

  const affectedTopicsFromData = (data: Record<string, unknown>) => {
    const cudCount = cudTotalFromData(data);
    if (dataHasCudCounts(data)) return cudCount;
    const affected = logNumberValue(data.affectedTopics);
    if (affected > 0 || dataHasField(data, 'affectedTopics')) return affected;
    if (cudCount > 0) return cudCount;
    const importFiles = logNumberValue(data.importFiles);
    if (importFiles > 0) return importFiles;
    return 0;
  };

  const logValueText = (key: string, value: unknown): ReactNode => {
    if (value === null || value === undefined || value === '') {
      return <span className={styles.logMuted}>-</span>;
    }
    if (typeof value === 'boolean') {
      return formatMessage(value ? 'uns.true' : 'uns.false', undefined, value ? '是' : '否');
    }
    if (key === 'bytes') {
      return formatBytes(value);
    }
    const text = String(value);
    if (key === 'syncMode') return syncModeText(text);
    if ((key === 'endpoint' || key === 'httpEndpoint' || key === 'httpBaseURL') && /^https?:\/\//i.test(text)) {
      return displayBaseURL(text);
    }
    if (key === 'role' && text === 'branch')
      return formatMessage('mqttAuth.syncLogValue.branch', undefined, '分支节点');
    if (key === 'role' && text === 'host') return formatMessage('mqttAuth.syncLogValue.host', undefined, 'Host 节点');
    if (key === 'mountMode' && text === 'root') return formatMessage('mqttAuth.syncLogValue.root', undefined, '根目录');
    if (key === 'mountMode' && text === 'singleTree')
      return formatMessage('mqttAuth.syncLogValue.singleTree', undefined, '单目录');
    if (key === 'mountMode' && text === 'multiRoots')
      return formatMessage('mqttAuth.syncLogValue.multiRoots', undefined, '多目录');
    if (key === 'transport' && text === 'http') return 'HTTP';
    if (key === 'transport' && text === 'mqtt') return 'MQTT';
    if (key === 'transport' && text === 'mqtt_fallback')
      return formatMessage('mqttAuth.syncLogValue.mqttFallback', undefined, 'MQTT 回退');
    if (key === 'fallback' && text === 'http')
      return formatMessage('mqttAuth.syncLogValue.httpFallback', undefined, 'HTTP 失败后回退');
    if (key === 'fallback' && text === 'mqtt_unavailable')
      return formatMessage('mqttAuth.syncLogValue.mqttUnavailable', undefined, 'MQTT 不可用');
    if (key === 'fallback' && text === 'skipped_auth_error')
      return formatMessage('mqttAuth.syncLogValue.authFallbackSkipped', undefined, '鉴权失败，未回退');
    return text;
  };

  const orderedNestedEntries = (value: Record<string, unknown>) => {
    const ordered = nestedLogFieldPriority.filter((key) => Object.prototype.hasOwnProperty.call(value, key));
    const rest = Object.keys(value)
      .filter((key) => !ordered.includes(key) && !hiddenErrorItemKeys.has(key))
      .sort();
    return [...ordered, ...rest].map((key) => [key, value[key]] as const);
  };

  const normalizeErrorRows = (value: unknown[]) =>
    value.map((item) => (item && typeof item === 'object' ? (item as Record<string, unknown>) : { topic: item }));

  const errorTopicOf = (row: Record<string, unknown>) =>
    row.topic || row.namespace || row.targetPath || row.sourcePath || row.alias || '-';

  const isConflictErrorRow = (row: Record<string, unknown>) => {
    const text = ['errmsg', 'reason', 'message', 'error']
      .map((key) => row[key])
      .filter((value) => value !== undefined && value !== null)
      .join(' ')
      .toLowerCase();
    return text.includes('conflict') || text.includes('冲突');
  };

  const errorTypeLabel = (key: string, row?: Record<string, unknown>) => {
    if (key === 'conflictNodes' || (row && isConflictErrorRow(row))) {
      return formatMessage('mqttAuth.syncLogError.conflict', undefined, '冲突');
    }
    if (key === 'failedNodes') return formatMessage('mqttAuth.syncLogError.failed', undefined, '失败');
    if (key === 'skippedNodeItems') return formatMessage('mqttAuth.syncLogError.skipped', undefined, '跳过');
    if (key === 'httpError') return 'HTTP';
    return formatMessage('mqttAuth.syncLogError.error', undefined, '错误');
  };

  const collectErrorRows = (data: Record<string, unknown>) => {
    const rows: Array<{ key: string; label: string; topic: unknown; detail: Record<string, unknown> }> = [];
    errorListLogKeys.forEach((key) => {
      const value = data[key];
      if (!Array.isArray(value) || !value.length) return;
      normalizeErrorRows(value).forEach((row) => {
        rows.push({
          key,
          label: errorTypeLabel(key, row),
          topic: errorTopicOf(row),
          detail: row,
        });
      });
    });
    errorScalarLogKeys.forEach((key) => {
      const value = data[key];
      if (value === undefined || value === null || value === '') return;
      rows.push({
        key,
        label: errorTypeLabel(key, { [key]: value }),
        topic: data.topic || '-',
        detail: {
          topic: data.topic || '-',
          [key]: value,
        },
      });
    });
    return rows;
  };

  const abnormalTopicsFromData = (data: Record<string, unknown>, errorRows: ReturnType<typeof collectErrorRows>) => {
    if (dataHasField(data, 'abnormalCount')) return logNumberValue(data.abnormalCount);
    const failed = logNumberValue(data.failedCount);
    const conflicts = logNumberValue(data.conflictCount);
    const skipped = dataHasField(data, 'skippedCount') ? logNumberValue(data.skippedCount) : 0;
    const explicit = failed + conflicts + skipped;
    if (explicit > 0 || dataHasField(data, 'failedCount') || dataHasField(data, 'conflictCount')) {
      return explicit;
    }
    const listCount =
      logNumberValue(data.failedNodes) + logNumberValue(data.conflictNodes) + logNumberValue(data.skippedNodeItems);
    return Math.max(errorRows.length, listCount);
  };

  const logSummaryStats = (data: Record<string, unknown>) => {
    const cud = cudCountsFromData(data);
    const errorRows = collectErrorRows(data);
    const affected = affectedTopicsFromData(data);
    const abnormal = abnormalTopicsFromData(data, errorRows);
    const total = dataHasField(data, 'totalTopics') ? logNumberValue(data.totalTopics) : affected + abnormal;
    return {
      total,
      affected,
      cud,
      abnormal,
    };
  };

  const unsSyncStatLabel = (key: string) => {
    if (key === 'total') return formatMessage('mqttAuth.syncLogStat.total', undefined, '总数');
    if (key === 'created') return logFieldLabel('createdFiles');
    if (key === 'updated') return logFieldLabel('updatedFiles');
    if (key === 'deleted') return logFieldLabel('deletedFiles');
    if (key === 'failed') return formatMessage('mqttAuth.syncLogStat.failed', undefined, '失败');
    return key;
  };

  const unsSyncStatValue = (stats: ReturnType<typeof logSummaryStats>, key: string) => {
    if (key === 'total') return stats.total;
    if (key === 'created') return stats.cud.created;
    if (key === 'updated') return stats.cud.updated;
    if (key === 'deleted') return stats.cud.deleted;
    if (key === 'failed') return stats.abnormal;
    return 0;
  };

  const unsSyncNodeDetailLimit = (data: Record<string, unknown>) => {
    const value = logNumberValue(data.nodeDetailLimit);
    return value > 0 ? value : defaultUnsSyncNodeDetailLimit;
  };

  const unsSyncNodeTotalByStat = (data: Record<string, unknown>, key: string, rowsLength: number) => {
    if (key === 'created') return Math.max(logNumberValue(data.createdFiles), rowsLength);
    if (key === 'updated') return Math.max(logNumberValue(data.updatedFiles), rowsLength);
    if (key === 'deleted') return Math.max(logNumberValue(data.deletedFiles), rowsLength);
    if (key === 'failed') return Math.max(abnormalTopicsFromData(data, collectErrorRows(data)), rowsLength);
    return rowsLength;
  };

  const activeSyncLogStatKey = (record: CloudSyncLogItem) => {
    const selected = syncLogDetailStat[syncLogKeyOf(record)];
    if (selected && selected !== 'total' && unsSyncStatKeys.includes(selected)) return selected;
    const stats = logSummaryStats(syncLogDetailsData(record));
    return unsSyncStatKeys.find((key) => key !== 'total' && unsSyncStatValue(stats, key) > 0) || '';
  };

  const collectUnsSyncNodeRows = (
    data: Record<string, unknown>,
    dataKey: string,
    label: string,
    danger = false
  ): UnsSyncNodeRow[] => {
    const value = data[dataKey];
    if (!Array.isArray(value) || !value.length) return [];
    return normalizeErrorRows(value).map((row, index) => ({
      key: `${dataKey}-${index}`,
      label,
      topic: errorTopicOf(row),
      detail: row,
      danger,
    }));
  };

  const collectUnsSyncFailedRows = (data: Record<string, unknown>): UnsSyncNodeRow[] =>
    collectErrorRows(data).map((row, index) => ({
      key: `${row.key}-${index}`,
      label: row.label,
      topic: row.topic,
      detail: row.detail,
      danger: true,
    }));

  const unsSyncRowsByStat = (data: Record<string, unknown>, key: string) => {
    if (key === 'created') return collectUnsSyncNodeRows(data, 'createdNodes', unsSyncStatLabel('created'));
    if (key === 'updated') return collectUnsSyncNodeRows(data, 'updatedNodes', unsSyncStatLabel('updated'));
    if (key === 'deleted') return collectUnsSyncNodeRows(data, 'deletedNodes', unsSyncStatLabel('deleted'));
    if (key === 'failed') return collectUnsSyncFailedRows(data);
    return [];
  };

  const renderUnsSyncNodeRow = (row: UnsSyncNodeRow) => {
    const detailText = [row.detail.errmsg, row.detail.reason]
      .filter((value) => value !== undefined && value !== null && value !== '')
      .map((value) => String(value))
      .join(' / ');
    return (
      <div className={styles.logNodeRow}>
        <span className={classNames(styles.logNodeBadge, row.danger && styles.logNodeBadgeDanger)}>{row.label}</span>
        <span className={styles.logNodeTopic}>{logValueText('topic', row.topic)}</span>
        {detailText ? <span className={styles.logNodeReason}>{detailText}</span> : null}
      </div>
    );
  };

  const renderUnsSyncNodeSection = (data: Record<string, unknown>, key: string, forceVisible = false) => {
    const rows = unsSyncRowsByStat(data, key);
    const total = unsSyncNodeTotalByStat(data, key, rows.length);
    const limit = unsSyncNodeDetailLimit(data);
    const truncated = total > rows.length;
    if (!rows.length && !forceVisible) return null;
    return (
      <div className={styles.logNodeSection}>
        <div className={styles.logNodeSectionTitle}>
          <strong>{unsSyncStatLabel(key)}</strong>
          <Tag>{truncated ? `${rows.length}/${total}` : rows.length}</Tag>
        </div>
        {truncated ? (
          <div className={styles.logNodeLimitHint}>
            {formatMessage('mqttAuth.syncLogDetailLimitHint', { limit }, `当前记录仅保留前 ${limit} 条明细。`)}
          </div>
        ) : null}
        <div className={styles.logNodeList}>
          {rows.length ? (
            rows.map((row) => <div key={row.key}>{renderUnsSyncNodeRow(row)}</div>)
          ) : (
            <div className={styles.logEmpty}>{formatMessage('common.noData', {}, 'No data')}</div>
          )}
        </div>
      </div>
    );
  };

  const renderSelectedUnsSyncSection = (data: Record<string, unknown>, activeKey: string) => {
    if (!activeKey) {
      return <div className={styles.logEmpty}>{formatMessage('common.noData')}</div>;
    }
    if (activeKey === 'total') {
      const sections = unsSyncNodeStatKeys.map((key) => renderUnsSyncNodeSection(data, key)).filter(Boolean);
      return sections.length ? sections : <div className={styles.logEmpty}>{formatMessage('common.noData')}</div>;
    }
    return renderUnsSyncNodeSection(data, activeKey, true);
  };

  const renderLogStats = (record: CloudSyncLogItem, stats: ReturnType<typeof logSummaryStats>, activeKey: string) => {
    const logKey = syncLogKeyOf(record);
    return (
      <div className={styles.logStats}>
        {unsSyncStatKeys.map((key) => {
          const clickable = key !== 'total';
          return (
            <div
              key={key}
              role={clickable ? 'button' : undefined}
              tabIndex={clickable ? 0 : undefined}
              className={classNames(
                clickable ? styles.logStatItem : styles.logStatItemDisabled,
                clickable && activeKey === key && styles.logStatItemActive
              )}
              onClick={clickable ? () => setSyncLogDetailStat((prev) => ({ ...prev, [logKey]: key })) : undefined}
              onKeyDown={
                clickable
                  ? (event) => {
                      if (event.key === 'Enter' || event.key === ' ') {
                        event.preventDefault();
                        setSyncLogDetailStat((prev) => ({ ...prev, [logKey]: key }));
                      }
                    }
                  : undefined
              }
            >
              <span>{unsSyncStatLabel(key)}</span>
              <strong>{unsSyncStatValue(stats, key)}</strong>
            </div>
          );
        })}
      </div>
    );
  };

  const renderNestedLogObject = (value: Record<string, unknown>) => (
    <div className={styles.logNested}>
      {orderedNestedEntries(value).map(([key, child]) => (
        <div key={key}>
          <span>{logFieldLabel(key)}</span>
          <strong>
            {Array.isArray(child) || (child && typeof child === 'object')
              ? JSON.stringify(child)
              : logValueText(key, child)}
          </strong>
        </div>
      ))}
    </div>
  );

  const renderLogValue = (key: string, value: unknown): ReactNode => {
    if (Array.isArray(value)) {
      if (value.length === 0)
        return <span className={styles.logMuted}>{formatMessage('common.none', undefined, '无')}</span>;
      return (
        <div className={styles.logListValue}>
          {value.map((item, index) => (
            <div key={`${key}-${index}`} className={styles.logListItem}>
              {item && typeof item === 'object'
                ? renderNestedLogObject(item as Record<string, unknown>)
                : logValueText(key, item)}
            </div>
          ))}
        </div>
      );
    }
    if (value && typeof value === 'object') {
      return renderNestedLogObject(value as Record<string, unknown>);
    }
    return logValueText(key, value);
  };

  const renderLogSection = (record: CloudSyncLogItem) => {
    const data = syncLogDetailsData(record);
    const stats = logSummaryStats(data);
    if (stats.total <= 0) {
      return <div className={styles.logEmpty}>-</div>;
    }
    const activeKey = activeSyncLogStatKey(record);
    return (
      <div className={styles.logDetails}>
        {renderLogStats(record, stats, activeKey)}
        {renderSelectedUnsSyncSection(data, activeKey)}
      </div>
    );
  };

  const renderSyncLogFocusedDetail = (record: CloudSyncLogItem) => {
    const summary = record.errorMessage || record.summary || '-';
    return (
      <div className={styles.syncLogFocused}>
        <div className={styles.syncLogFocusToolbar}>
          <Button
            variant="outlined"
            color="default"
            className={styles.detailBackBtn}
            onClick={() => setFocusedSyncLog(null)}
          >
            <Flex align="center" gap={8}>
              <ChevronLeft size={16} />
              {formatMessage('mqttAuth.syncLogBackToList', {}, '返回同步历史')}
            </Flex>
          </Button>
        </div>
        <div className={styles.syncLogFocusMeta}>
          <div>
            <span>{formatMessage('mqttAuth.syncTime')}</span>
            <strong>{record.createdTime ? formatTimestamp(record.createdTime) : '-'}</strong>
          </div>
          <div>
            <span>{formatMessage('mqttAuth.syncStatus')}</span>
            <div className={styles.metaTagValue}>{syncStatusTag(record.status)}</div>
          </div>
          <div>
            <span>{formatMessage('mqttAuth.syncSummary')}</span>
            <strong title={summary !== '-' ? summary : undefined}>{summary}</strong>
          </div>
        </div>
        <div className={styles.syncLogFocusBody}>{renderLogSection(record)}</div>
      </div>
    );
  };

  const renderTokenMetadata = (node: ClusterNode | null) => {
    const metadata = parseTokenMetadata(node?.token);
    const entries = visibleTokenMetadataEntries(metadata);
    return (
      <div className={styles.drawerMeta}>
        <span className={styles.drawerMetaTitle}>
          {formatMessage('mqttAuth.tokenMetadata', undefined, 'Token Metadata')}
        </span>
        <div className={styles.drawerMetaGrid}>
          {entries.map(([key, value]) => (
            <div key={key}>
              <span>{logFieldLabel(key)}</span>
              <strong>{renderLogValue(key, value)}</strong>
            </div>
          ))}
        </div>
      </div>
    );
  };

  const syncLogColumns: any = [
    {
      title: formatMessage('mqttAuth.syncTime'),
      dataIndex: 'createdTime',
      width: 155,
      ellipsis: true,
      render: (value: number) => (value ? formatTimestamp(value) : '-'),
    },
    {
      title: formatMessage('mqttAuth.syncStatus'),
      dataIndex: 'status',
      width: 100,
      render: (value: string) => syncStatusTag(value),
    },
    {
      title: formatMessage('mqttAuth.syncSummary'),
      dataIndex: 'summary',
      minWidth: 220,
      ellipsis: true,
      render: (value: string, record: CloudSyncLogItem) => {
        const text = record.errorMessage || value || '-';
        return (
          <span className={styles.cellEllipsis} title={text !== '-' ? text : undefined}>
            {text}
          </span>
        );
      },
    },
  ];

  const syncLogEnvironmentData = (record: CloudSyncLogItem) =>
    parseLogObject(record.environment) as Record<string, unknown>;

  const userSyncDetailsData = (record: CloudSyncLogItem) => parseLogObject(record.details) as Record<string, unknown>;

  const userSyncSourceText = (record: CloudSyncLogItem) => {
    const env = syncLogEnvironmentData(record);
    const source = String(env.source || '').toLowerCase();
    if (source === 'cloud') return formatMessage('cluster.userSyncSourceCloud', {}, 'Cloud');
    if (source === 'host') return formatMessage('cluster.userSyncSourceHost', {}, 'Host');
    if (record.syncMode === 'cloudEdge') return formatMessage('cluster.userSyncSourceCloud', {}, 'Cloud');
    if (record.syncMode === 'edgeEdge') return formatMessage('cluster.userSyncSourceHost', {}, 'Host');
    return source || '-';
  };

  const userSyncStatKeys = ['total', 'created', 'updated', 'disabled', 'deleted', 'skipped', 'conflicts'];

  const userSyncStatLabel = (key: string) =>
    formatMessage(`cluster.userSyncStat.${key}`, undefined, logFieldFallbacks[key] || key);

  const userSyncStats = (record: CloudSyncLogItem) => {
    const data = userSyncDetailsData(record);
    const stats = data.stats && typeof data.stats === 'object' && !Array.isArray(data.stats) ? data.stats : {};
    return stats as Record<string, unknown>;
  };

  const activeUserSyncStatKey = (record: CloudSyncLogItem) => {
    const selected = userSyncDetailStat[syncLogKeyOf(record)];
    if (selected && selected !== 'total' && userSyncStatKeys.includes(selected)) return selected;
    const stats = userSyncStats(record);
    return userSyncStatKeys.find((key) => key !== 'total' && Number(stats[key] || 0) > 0) || '';
  };

  const userSyncSummaryText = (record: CloudSyncLogItem) => {
    if (record.errorMessage) return record.errorMessage;
    const stats = userSyncStats(record);
    const pieces = userSyncStatKeys
      .filter((key) => key !== 'total')
      .map((key) => [key, Number(stats[key] || 0)] as const)
      .filter(([, value]) => value > 0)
      .map(([key, value]) => `${userSyncStatLabel(key)} ${value}`);
    if (pieces.length) return pieces.join(' / ');
    return record.summary || '-';
  };

  const userSyncRowsOf = (record: CloudSyncLogItem, key: string) => {
    const data = userSyncDetailsData(record);
    return Array.isArray(data[key]) ? (data[key] as Record<string, unknown>[]) : [];
  };

  const isUserSyncLogExpandable = (record: CloudSyncLogItem) => {
    const data = userSyncDetailsData(record);
    if (data.stats && typeof data.stats === 'object') return true;
    return [
      'createdUsers',
      'updatedUsers',
      'disabledUsers',
      'deletedUsers',
      'skippedUsers',
      'conflicts',
      'roleMappings',
    ].some((key) => Array.isArray(data[key]) && (data[key] as unknown[]).length > 0);
  };

  const renderUserSyncStats = (record: CloudSyncLogItem, activeKey: string) => {
    const stats = userSyncStats(record);
    const logKey = syncLogKeyOf(record);
    return (
      <div className={styles.userSyncStats}>
        {userSyncStatKeys.map((key) => {
          const clickable = key !== 'total';
          return (
            <div
              key={key}
              role={clickable ? 'button' : undefined}
              tabIndex={clickable ? 0 : undefined}
              className={classNames(
                styles.userSyncStatItem,
                !clickable && styles.userSyncStatItemDisabled,
                clickable && activeKey === key && styles.userSyncStatItemActive
              )}
              onClick={clickable ? () => setUserSyncDetailStat((prev) => ({ ...prev, [logKey]: key })) : undefined}
              onKeyDown={
                clickable
                  ? (event) => {
                      if (event.key === 'Enter' || event.key === ' ') {
                        event.preventDefault();
                        setUserSyncDetailStat((prev) => ({ ...prev, [logKey]: key }));
                      }
                    }
                  : undefined
              }
            >
              <span>{userSyncStatLabel(key)}</span>
              <strong>{Number(stats[key] || 0)}</strong>
            </div>
          );
        })}
      </div>
    );
  };

  const renderUserSyncUserRow = (row: Record<string, unknown>, extra?: ReactNode) => {
    const title = String(row.userName || row.email || row.userId || '-');
    const meta = [row.email ? String(row.email) : '', row.userId ? `#${row.userId}` : ''].filter(Boolean).join(' / ');
    return (
      <div className={styles.userSyncUserRow}>
        <div>
          <strong title={title}>{title}</strong>
          <span title={meta}>{meta || '-'}</span>
        </div>
        {row.reason || row.detail ? (
          <div className={styles.userSyncReason}>
            {[row.reason, row.detail]
              .filter(Boolean)
              .map((value) => String(value))
              .join(' / ')}
          </div>
        ) : null}
        {extra}
      </div>
    );
  };

  const renderUserSyncUserSection = (
    record: CloudSyncLogItem,
    key: string,
    titleKey: string,
    extra?: ReactNode,
    forceVisible = false
  ) => {
    const rows = userSyncRowsOf(record, key);
    if (!rows.length && !forceVisible) return null;
    return (
      <div className={styles.userSyncDetailSection}>
        <div className={styles.userSyncDetailTitle}>
          <strong>{formatMessage(titleKey)}</strong>
          <Tag>{rows.length}</Tag>
        </div>
        <div className={styles.userSyncUserList}>
          {rows.length ? (
            rows.map((row, index) => <div key={`${key}-${index}`}>{renderUserSyncUserRow(row, extra)}</div>)
          ) : (
            <div className={styles.userSyncEmpty}>{formatMessage('common.noData', {}, 'No data')}</div>
          )}
        </div>
      </div>
    );
  };

  const renderSelectedUserSyncSection = (record: CloudSyncLogItem, activeKey: string) => {
    if (!activeKey) {
      return <div className={styles.userSyncEmpty}>{formatMessage('common.noData', {}, 'No data')}</div>;
    }
    if (activeKey === 'total') {
      return (
        <>
          {Object.values(userSyncUserSectionByStat).map(({ dataKey, titleKey }) =>
            renderUserSyncUserSection(record, dataKey, titleKey)
          )}
        </>
      );
    }
    const section = userSyncUserSectionByStat[activeKey];
    if (!section) return null;
    return renderUserSyncUserSection(record, section.dataKey, section.titleKey, undefined, true);
  };

  const renderRoleMappings = (record: CloudSyncLogItem) => {
    const rows = userSyncRowsOf(record, 'roleMappings');
    if (!rows.length) return null;
    const sourceLabel = formatMessage(
      'cluster.userSyncRoleSource',
      { source: userSyncSourceText(record) },
      'Source role'
    );
    const targetLabel = formatMessage('cluster.userSyncRoleTarget', {}, 'Branch role');
    return (
      <div className={styles.userSyncDetailSection}>
        <div className={styles.userSyncDetailTitle}>
          <strong>{formatMessage('cluster.userSyncRoleMappings', {}, 'Role mappings')}</strong>
          <Tag>{rows.length}</Tag>
        </div>
        <div className={styles.userSyncRoleList}>
          <div className={styles.userSyncRoleHeader}>
            <span>{sourceLabel}</span>
            <span>{targetLabel}</span>
            <span>{formatMessage('common.status')}</span>
          </div>
          {rows.map((row, index) => {
            const source = String(row.sourceRoleName || row.sourceRoleCode || row.sourceRoleId || '-');
            const target = String(row.targetRoleName || row.targetRoleCode || row.targetRoleId || '-');
            const mapped = row.status === 'mapped';
            return (
              <div key={`role-${index}`} className={styles.userSyncRoleRow}>
                <div className={styles.userSyncRoleCell}>
                  <em>{sourceLabel}</em>
                  <strong title={source}>{source}</strong>
                </div>
                <div className={styles.userSyncRoleCell}>
                  <em>{targetLabel}</em>
                  <strong title={target}>{target}</strong>
                </div>
                <Tag className={mapped ? styles.successTag : undefined}>
                  {mapped
                    ? formatMessage('cluster.userSyncRoleMapped', {}, 'Mapped')
                    : formatMessage('cluster.userSyncRoleMissing', {}, 'Missing')}
                </Tag>
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  const renderUserSyncDetails = (record: CloudSyncLogItem) => {
    const activeKey = activeUserSyncStatKey(record);
    return (
      <div className={styles.userSyncDetails}>
        {renderUserSyncStats(record, activeKey)}
        {renderSelectedUserSyncSection(record, activeKey)}
        {renderRoleMappings(record)}
      </div>
    );
  };

  const renderUserSyncFocusedDetail = (record: CloudSyncLogItem) => {
    const summary = userSyncSummaryText(record);
    return (
      <div className={styles.syncLogFocused}>
        <div className={styles.syncLogFocusToolbar}>
          <Button
            variant="outlined"
            color="default"
            className={styles.detailBackBtn}
            onClick={() => setFocusedUserSyncLog(null)}
          >
            <Flex align="center" gap={8}>
              <ChevronLeft size={16} />
              {formatMessage('mqttAuth.syncLogBackToList', {}, '返回同步历史')}
            </Flex>
          </Button>
        </div>
        <div className={styles.syncLogFocusMeta}>
          <div>
            <span>{formatMessage('mqttAuth.syncTime')}</span>
            <strong>{record.createdTime ? formatTimestamp(record.createdTime) : '-'}</strong>
          </div>
          <div>
            <span>{formatMessage('mqttAuth.syncStatus')}</span>
            <div className={styles.metaTagValue}>{syncStatusTag(record.status)}</div>
          </div>
          <div>
            <span>{formatMessage('mqttAuth.syncSummary')}</span>
            <strong title={summary !== '-' ? summary : undefined}>{summary}</strong>
          </div>
        </div>
        <div className={styles.syncLogFocusBody}>{renderUserSyncDetails(record)}</div>
      </div>
    );
  };

  const userSyncColumns: any = [
    {
      title: formatMessage('mqttAuth.syncTime'),
      dataIndex: 'createdTime',
      width: 160,
      render: (value: number) => (value ? formatTimestamp(value) : '-'),
    },
    {
      title: formatMessage('mqttAuth.syncStatus'),
      dataIndex: 'status',
      width: 112,
      render: (value: string) => syncStatusTag(value),
    },
    {
      title: formatMessage('cluster.userSyncSource', {}, 'Source'),
      dataIndex: 'source',
      width: 110,
      render: (_: unknown, record: CloudSyncLogItem) => userSyncSourceText(record),
    },
    {
      title: formatMessage('mqttAuth.syncSummary'),
      dataIndex: 'summary',
      minWidth: 260,
      ellipsis: true,
      render: (_: string, record: CloudSyncLogItem) => {
        const text = userSyncSummaryText(record);
        return (
          <span className={styles.cellEllipsis} title={text !== '-' ? text : undefined}>
            {text}
          </span>
        );
      },
    },
  ];

  const edgeColumns: any = [
    {
      title: formatMessage('mqttAuth.nodeName'),
      dataIndex: 'name',
      width: 220,
      render: (_: string, record: ClusterNode) => nodeNameOf(record) || '-',
    },
    {
      title: formatMessage('mqttAuth.nodeID'),
      dataIndex: 'nodeID',
      width: 220,
      render: (_: string, record: ClusterNode) => nodeKeyOf(record) || '-',
    },
    {
      title: formatMessage('common.description'),
      dataIndex: 'description',
      width: 220,
      ellipsis: true,
      render: (value: string) => value || '-',
    },
    {
      title: formatMessage('common.status'),
      dataIndex: 'status',
      width: 120,
      render: (value: string | number) => (
        <Tag className={nodeStatusText(value) === 'online' ? styles.successTag : undefined}>
          {formatMessage(statusMessageKey(value))}
        </Tag>
      ),
    },
    {
      title: formatMessage('mqttAuth.lastSeen'),
      dataIndex: 'lastSeen',
      width: 180,
      render: (_: number, record: ClusterNode) => {
        const value = nodeLastSeenOf(record);
        return value ? formatTimestamp(value) : '-';
      },
    },
    {
      title: formatMessage('mqttAuth.authToken', undefined, 'Token'),
      dataIndex: 'token',
      width: 150,
      render: (value: string) =>
        value ? (
          <span className={styles.tokenCell}>
            <span className={styles.tokenMask}>••••••</span>
            <Button
              type="text"
              size="small"
              className={styles.tokenCopyButton}
              aria-label={formatMessage('mqttAuth.copyToken')}
              icon={<Copy size={14} />}
              onClick={(event) => {
                event.stopPropagation();
                copyValue(value);
              }}
            />
          </span>
        ) : (
          '-'
        ),
    },
  ];

  const syncLogPaginationProps = (node: ClusterNode | null) => ({
    total: syncLogPagination.total,
    current: syncLogPagination.current,
    pageSize: syncLogPagination.pageSize,
    showTotal: (total: number) =>
      formatMessage(
        'mqttAuth.syncLogPaginationTotal',
        { total, limit: defaultUnsSyncNodeDetailLimit },
        `共 ${total} 条记录，每条最多保留 ${defaultUnsSyncNodeDetailLimit} 条明细。`
      ),
    onChange: (page: number, pageSize: number) => {
      void loadSyncLogs(node, page, pageSize);
    },
    onShowSizeChange: (page: number, pageSize: number) => {
      void loadSyncLogs(node, page, pageSize);
    },
  });

  const userSyncPaginationProps = {
    total: userSyncPagination.total,
    current: userSyncPagination.current,
    pageSize: userSyncPagination.pageSize,
    showTotal: (total: number) => formatMessage('cluster.userSyncHistoryTotal', { total }, `共 ${total} 条记录`),
    onChange: (page: number, pageSize: number) => {
      void loadUserSyncLogs(page, pageSize);
    },
    onShowSizeChange: (page: number, pageSize: number) => {
      void loadUserSyncLogs(page, pageSize);
    },
  };

  const renderSyncLogTable = (node: ClusterNode | null, options?: { sectionTitle?: string }) => (
    <div className={options?.sectionTitle ? styles.detailSection : styles.syncLogTable}>
      {options?.sectionTitle ? <span className={styles.detailSectionTitle}>{options.sectionTitle}</span> : null}
      <div className={styles.syncLogTable}>
        <ProTable
          columns={syncLogColumns}
          dataSource={syncLogs}
          rowKey={syncLogKeyOf}
          loading={syncLogLoading}
          pagination={syncLogPaginationProps(node)}
          rowClassName={() => styles.syncLogExpandableRow}
          scroll={{ x: 480 }}
          onRow={(record: CloudSyncLogItem) => ({
            onClick: () => {
              setFocusedSyncLog(record);
            },
          })}
        />
      </div>
    </div>
  );

  const renderUserSyncLogTable = () => (
    <div className={styles.syncLogTable}>
      <ProTable
        columns={userSyncColumns}
        dataSource={userSyncLogs}
        rowKey={syncLogKeyOf}
        loading={userSyncLoading}
        pagination={userSyncPaginationProps}
        rowClassName={(record: CloudSyncLogItem) =>
          isUserSyncLogExpandable(record) ? styles.syncLogExpandableRow : ''
        }
        scroll={{ x: 680 }}
        onRow={(record: CloudSyncLogItem) => ({
          onClick: () => {
            if (isUserSyncLogExpandable(record)) setFocusedUserSyncLog(record);
          },
        })}
      />
    </div>
  );

  const renderBranchEdgeNodeContent = () => {
    const connected = activeCloudSyncConfig?.connectStatus === 'connected';
    const isCloudEdgeMode = activeCloudSyncConfig?.syncMode === 'cloudEdge';
    const mqttAuthID = activeCloudSyncConfig?.mqttAuthID || activeCloudSyncConfig?.connectClientKey || '';
    const topicRoot = String(activeCloudSyncConfig?.topicRoot || defaultClusterSyncTopicRoot).replace(/^\/+|\/+$/g, '');
    const buildSyncTopic = (kind: string) => (mqttAuthID ? `${topicRoot}/${mqttAuthID}/${kind}` : '-');
    const connectionEndpoint = isCloudEdgeMode
      ? displayBaseURL(activeCloudSyncConfig?.httpEndpoint) || '-'
      : activeCloudSyncConfig?.mqttBrokers ||
        formatMessage('cluster.mqttEndpointPending', {}, 'Waiting for broker info');
    const connectionTarget =
      activeCloudSyncConfig?.mqttAuthID ||
      activeCloudSyncConfig?.connectClientKey ||
      activeCloudSyncConfig?.edgeNodeName ||
      '-';
    const unsSyncRoute = isCloudEdgeMode ? cloudSyncMetadataAPIPath : buildSyncTopic('meta');
    const lastDisconnectTime = activeCloudSyncConfig?.lastDisconnectTime || 0;
    const unsReportedAt = activeLastUnsSyncLog?.createdTime || activeCloudSyncConfig?.lastSyncTime || 0;
    const userReportedAt = activeLastUserSyncLog?.createdTime || 0;
    const lastError = activeCloudSyncConfig?.lastError?.trim() || '-';

    const renderBranchInfoField = (label: ReactNode, value: ReactNode) => (
      <div className={styles.branchInfoField}>
        <span className={styles.branchInfoLabel}>{label}</span>
        <div className={styles.branchInfoValue}>{value}</div>
      </div>
    );

    const renderConnectionStatusTag = () => (
      <Tag className={classNames(styles.syncStatusTag, connected ? styles.successTag : styles.neutralTag)}>
        {formatMessage(statusMessageKey(activeCloudSyncConfig?.connectStatus))}
      </Tag>
    );

    return (
      <div className={styles.branchPanel}>
        <section className={styles.branchInfoSection}>
          <div className={styles.branchInfoTitle}>
            {formatMessage('mqttAuth.basicInformation', {}, 'Basic Information')}
          </div>
          <div className={styles.branchInfoRow}>
            {renderBranchInfoField(
              formatMessage('mqttAuth.connectionTarget', {}, 'Connection Target'),
              <span title={connectionTarget !== '-' ? connectionTarget : undefined}>{connectionTarget}</span>
            )}
            {renderBranchInfoField(
              formatMessage('mqttAuth.connectionStatus'),
              <div className={styles.metaTagValue}>{renderConnectionStatusTag()}</div>
            )}
            {renderBranchInfoField(
              formatMessage('cluster.mqttEndpoint', {}, 'MQTT Endpoint'),
              <span title={connectionEndpoint !== '-' ? connectionEndpoint : undefined}>{connectionEndpoint}</span>
            )}
          </div>
          <div className={styles.branchInfoRow}>
            {renderBranchInfoField(
              formatMessage('mqttAuth.lastConnectTime'),
              activeCloudSyncConfig?.lastConnectTime ? formatTimestamp(activeCloudSyncConfig.lastConnectTime) : '-'
            )}
            {renderBranchInfoField(
              formatMessage('mqttAuth.lastDisconnectTime'),
              lastDisconnectTime ? formatTimestamp(lastDisconnectTime) : '-'
            )}
            {renderBranchInfoField(
              formatMessage('mqttAuth.lastError'),
              <span title={lastError !== '-' ? lastError : undefined}>{lastError}</span>
            )}
          </div>
        </section>

        <hr className={styles.branchSectionDivider} aria-hidden />

        <section className={styles.branchBlock}>
          <div className={styles.branchBlockHeader}>
            <div className={styles.branchBlockTitle}>
              <Upload size={16} />
              <span>{formatMessage('cluster.unsUplink', {}, 'UNS Sync')}</span>
            </div>
          </div>
          <div className={styles.branchBlockCard}>
            <div className={styles.branchInfoRow}>
              {renderBranchInfoField(
                formatMessage('mqttAuth.selectedRootNode', {}, 'Selected Root Node'),
                <span title={unsSyncRoute !== '-' ? unsSyncRoute : undefined}>{unsSyncRoute}</span>
              )}
              {renderBranchInfoField(
                formatMessage('mqttAuth.syncStatus'),
                <div className={styles.metaTagValue}>
                  {activeLastUnsSyncLog ? (
                    syncStatusTag(activeLastUnsSyncLog.status)
                  ) : (
                    <Tag className={classNames(styles.syncStatusTag, styles.neutralTag)}>
                      {formatMessage('mqttAuth.syncUnknown')}
                    </Tag>
                  )}
                </div>
              )}
            </div>
            <div className={styles.branchInfoRow}>
              {renderBranchInfoField(
                formatMessage('mqttAuth.lastReportedAt', {}, 'Last Reported At'),
                unsReportedAt ? formatTimestamp(unsReportedAt) : '-'
              )}
              {renderBranchInfoField(
                formatMessage('mqttAuth.lastDisconnectTime'),
                lastDisconnectTime ? formatTimestamp(lastDisconnectTime) : '-'
              )}
            </div>
          </div>
        </section>

        <section className={styles.branchBlock}>
          <div className={styles.branchBlockHeader}>
            <div className={styles.branchBlockTitle}>
              <Download size={16} />
              <span>{formatMessage('cluster.userDownlink', {}, 'User Sync')}</span>
            </div>
            <div className={styles.branchBlockActions}>
              <Button className={styles.branchSecondaryAction} icon={<Time size={16} />} onClick={openUserSyncHistory}>
                {formatMessage('cluster.userSyncHistory', {}, 'History')}
              </Button>
              <Button
                className={styles.branchSecondaryAction}
                icon={<User size={16} />}
                onClick={openSyncedAccountList}
              >
                {formatMessage('cluster.syncedAccountList', {}, 'Synced Account List')}
              </Button>
              <Button
                className={styles.branchPrimaryAction}
                icon={<Download size={16} />}
                loading={cloudSyncUserSyncing}
                disabled={disablePreviewMutations}
                onClick={() => void handleCloudSyncUserSync()}
              >
                {formatMessage('cluster.pullUsers', {}, 'Pull Users')}
              </Button>
            </div>
          </div>
          <div className={styles.branchBlockCard}>
            <div className={styles.branchInfoRow}>
              {renderBranchInfoField(
                formatMessage('mqttAuth.lastReportedAt', {}, 'Last Reported At'),
                userReportedAt ? formatTimestamp(userReportedAt) : '-'
              )}
              {renderBranchInfoField(
                formatMessage('mqttAuth.syncStatus'),
                <div className={styles.metaTagValue}>
                  {activeLastUserSyncLog ? (
                    syncStatusTag(activeLastUserSyncLog.status)
                  ) : (
                    <Tag className={classNames(styles.syncStatusTag, styles.neutralTag)}>
                      {formatMessage('mqttAuth.syncUnknown')}
                    </Tag>
                  )}
                </div>
              )}
            </div>
          </div>
        </section>
      </div>
    );
  };

  const selectedNode = decodedNodeKey
    ? edgeNodeRows.find((node) => nodeKeyOf(node) === decodedNodeKey || String(node.id) === decodedNodeKey) ||
      (previewDetail ? PREVIEW_EDGE_NODES[0] : undefined)
    : undefined;

  useEffect(() => {
    if (!decodedNodeKey) return;
    const targetNode = selectedNode || (previewDetail ? PREVIEW_EDGE_NODES[0] : null);
    setSyncLogNode(targetNode);
    setFocusedSyncLog(null);
    void loadSyncLogs(targetNode, 1, syncLogPagination.pageSize);
  }, [decodedNodeKey, loadSyncLogs, previewDetail, selectedNode, syncLogPagination.pageSize]);

  if (decodedNodeKey) {
    const nodeTitle = nodeNameOf(selectedNode) || decodedNodeKey;
    const backToEdgeNodes = () => {
      const next = new URLSearchParams();
      next.set('tab', EDGE_NODE_TAB);
      if (isEdgeNodePreview) {
        next.set('edgeNodePreview', EDGE_NODE_PREVIEW_HOST);
      }
      navigate(`/edge-connection?${next.toString()}`);
    };

    return (
      <ComLayout loading={previewDetail ? false : edgeLoading} className={styles.edgePage}>
        <div className={styles.pageColumn}>
          <div className={`${styles.subHeader} ${styles.detailSubHeader}`}>
            <div className={styles.detailHeaderMain}>
              <ComBackButton onClick={backToEdgeNodes} />
              <span className={styles.pageTitle}>{nodeTitle}</span>
            </div>
          </div>

          <div className={styles.pageBody}>
            <div className={styles.main}>
              <div className={styles.contentFrame}>
                <div className={styles.detailPage}>
                  <div className={styles.detailSection}>{renderTokenMetadata(selectedNode || null)}</div>
                  <div className={styles.detailSection}>
                    <span className={styles.detailSectionTitle}>{formatMessage('mqttAuth.syncHistory')}</span>
                    {focusedSyncLog
                      ? renderSyncLogFocusedDetail(focusedSyncLog)
                      : renderSyncLogTable(selectedNode || null)}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </ComLayout>
    );
  }

  return (
    <ComLayout loading={loading || edgeLoading} className={styles.edgePage}>
      <div className={styles.pageColumn}>
        <div className={styles.subHeader}>
          <div className={styles.pageTitleWrap}>
            <PageTitleIcon resourceKey="mqtt.auth.manage" />
            <span className={styles.pageTitle}>{formatMessage('mqttAuth.title')}</span>
          </div>
          {activeTab === MQTT_ACCESS_TAB ? (
            <AuthButton
              auth={ButtonPermission['MqttAuth.add']}
              className={styles.headerBtn}
              icon={<Add {...toolbarIconProps} />}
              onClick={openCreate}
            >
              {formatMessage('mqttAuth.generateCredentials')}
            </AuthButton>
          ) : (
            <Space>
              {showEdgeNodeBranch ? (
                <Button
                  className={styles.headerBtn}
                  icon={<Add size={16} strokeWidth={1.75} />}
                  loading={cloudSyncConfigLoading && !isEdgeNodePreview}
                  disabled={disablePreviewMutations}
                  onClick={openCloudSyncConnect}
                >
                  {formatMessage('mqttAuth.setPrimaryNodeToken', {}, 'Set Primary Node Token')}
                </Button>
              ) : null}
              {showEdgeNodeHost ? (
                <Button
                  className={styles.headerBtn}
                  type="primary"
                  icon={<Add {...toolbarIconProps} />}
                  disabled={isEdgeNodePreview}
                  onClick={openCreateEdgeToken}
                >
                  {formatMessage('common.create')}
                </Button>
              ) : null}
            </Space>
          )}
        </div>

        <div className={styles.globalHeader}>
          <button
            type="button"
            className={classNames(
              styles.globalTab,
              activeTab === MQTT_ACCESS_TAB && styles.globalTabActive,
              activeTab !== MQTT_ACCESS_TAB && styles.globalTabInactive
            )}
            onClick={() => setTab(MQTT_ACCESS_TAB)}
          >
            <span className={styles.globalTabInner}>
              <Signal size={16} strokeWidth={1.75} />
              <span>{formatMessage('mqttAuth.mqttAccess')}</span>
            </span>
            {activeTab === MQTT_ACCESS_TAB ? <span className={styles.globalTabInk} aria-hidden /> : null}
          </button>
          {showEdgeNodeTab ? (
            <button
              type="button"
              className={classNames(
                styles.globalTab,
                activeTab === EDGE_NODE_TAB && styles.globalTabActive,
                activeTab !== EDGE_NODE_TAB && styles.globalTabInactive
              )}
              onClick={() => setTab(EDGE_NODE_TAB)}
            >
              <span className={styles.globalTabInner}>
                <ServerDns size={16} strokeWidth={1.75} />
                <span>{formatMessage('mqttAuth.edgeNode')}</span>
              </span>
              {activeTab === EDGE_NODE_TAB ? <span className={styles.globalTabInk} aria-hidden /> : null}
            </button>
          ) : null}
        </div>

        <div className={styles.pageBody}>
          <div className={styles.main}>
            <div className={styles.contentFrame}>
              <div className={styles.tableWrap}>
                {activeTab === MQTT_ACCESS_TAB ? (
                  <ProTable
                    columns={columns}
                    dataSource={dataSource}
                    rowKey="id"
                    pagination={false}
                    operationOptions={{
                      width: 80,
                      render: (record) => credentialOperationItems(record),
                    }}
                  />
                ) : showEdgeNodeHost ? (
                  <ProTable
                    columns={edgeColumns}
                    dataSource={edgeNodeRows}
                    rowKey="id"
                    loading={isEdgeNodePreview ? false : edgeLoading}
                    pagination={false}
                    rowClassName={() => styles.edgeNodeClickableRow}
                    onRow={(record) => ({
                      onClick: () => gotoNodeDetail(record as ClusterNode),
                    })}
                  />
                ) : showEdgeNodeBranch ? (
                  renderBranchEdgeNodeContent()
                ) : (
                  <div className={styles.emptyPanel}>
                    <ComEmpty description={formatMessage('mqttAuth.edgeNodeConnectOnly')} />
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      <ProModal
        open={createOpen}
        onCancel={closeCreate}
        title={formatMessage('mqttAuth.generateCredentials')}
        width={500}
        maskClosable={false}
        keyboard={false}
        destroyOnHidden
        styles={{
          header: {
            marginBottom: 0,
          },
          body: {
            paddingBlockStart: 16,
          },
        }}
      >
        {() => (
          <OperationForm
            form={form}
            onCancel={closeCreate}
            onSave={onCreateSave}
            formConfig={{
              layout: 'vertical',
              labelCol: { span: 24 },
              wrapperCol: { span: 24 },
            }}
            formItemOptions={[
              {
                name: 'name',
                label: formatMessage('mqttAuth.name'),
                rules: [
                  { required: true, whitespace: true, message: formatMessage('mqttAuth.nameRequired') },
                  { max: MAX_LENGTHS.connectionName, message: formatMessage('mqttAuth.nameMaxLength') },
                ],
                properties: {
                  placeholder: formatMessage('mqttAuth.namePlaceholder'),
                  maxLength: MAX_LENGTHS.connectionName,
                  allowClear: true,
                },
              },
              {
                name: 'description',
                label: formatMessage('common.description'),
                type: 'TextArea',
                properties: {
                  rows: 3,
                  maxLength: MAX_LENGTHS.description,
                  placeholder: formatMessage('route.optional'),
                },
              },
              {
                name: 'clientIdRandomSuffixEnabled',
                render: (
                  <Form.Item label={formatMessage('mqttAuth.clientIdRandomSuffix')}>
                    <Space size={8} align="center">
                      <Form.Item name="clientIdRandomSuffixEnabled" valuePropName="checked" noStyle>
                        <Switch />
                      </Form.Item>
                      <HelpTooltip
                        title={formatMessage('mqttAuth.randomSuffixHelp', { clientID: '{clientID}' })}
                        icon={<Help style={{ cursor: 'default' }} />}
                        styles={{ root: { maxWidth: 'none' }, body: { whiteSpace: 'nowrap' } }}
                      />
                    </Space>
                  </Form.Item>
                ),
              },
            ]}
            style={{ padding: 0, overflow: 'visible' }}
            footer={
              <Flex gap="10px" justify="end">
                <ComButton
                  color="default"
                  variant="filled"
                  onClick={closeCreate}
                  title={formatMessage('common.cancel')}
                >
                  {formatMessage('common.cancel')}
                </ComButton>
                <ComButton type="primary" variant="solid" onClick={onCreateSave} title={formatMessage('common.save')}>
                  {formatMessage('common.save')}
                </ComButton>
              </Flex>
            }
          />
        )}
      </ProModal>
      <Drawer
        title={formatMessage('mqttAuth.syncDetailTitle')}
        open={syncLogDrawerOpen}
        width="60vw"
        destroyOnClose
        onClose={() => {
          setSyncLogDrawerOpen(false);
          setFocusedSyncLog(null);
        }}
      >
        {focusedSyncLog ? (
          renderSyncLogFocusedDetail(focusedSyncLog)
        ) : (
          <>
            {renderTokenMetadata(syncLogNode)}
            {renderSyncLogTable(syncLogNode)}
          </>
        )}
      </Drawer>
      <Drawer
        title={formatMessage('cluster.userSyncHistoryTitle', {}, 'User sync history')}
        open={userSyncDrawerOpen}
        width="60vw"
        destroyOnClose
        onClose={() => {
          setUserSyncDrawerOpen(false);
          setFocusedUserSyncLog(null);
        }}
      >
        {focusedUserSyncLog ? renderUserSyncFocusedDetail(focusedUserSyncLog) : renderUserSyncLogTable()}
      </Drawer>
      <CloudSyncConnectModal
        open={cloudSyncConnectOpen}
        initialConfig={cloudSyncConfig}
        titleMessageKey="mqttAuth.setPrimaryNodeToken"
        primaryActionMessageKey="common.save"
        onCancel={() => setCloudSyncConnectOpen(false)}
        onSaved={(config) => {
          setCloudSyncConfig(config);
          void loadBranchLatestSyncLogs();
        }}
      />
    </ComLayout>
  );
};

export default MqttAuthPage;
