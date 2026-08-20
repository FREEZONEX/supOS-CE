import type { CloudSyncConfigResp, CloudSyncLogItem } from '@/apis/core-api/cloudsync.ts';

type PreviewClusterNode = {
  id: number;
  nodeKey?: string;
  nodeName?: string;
  nodeID?: string;
  name?: string;
  status?: string | number;
  description?: string;
  lastSeenTime?: number;
  connectClientID?: number;
  lastError?: string;
  token?: string;
};

type PreviewCloudSyncConfig = CloudSyncConfigResp & {
  mqttBrokers?: string;
  mqttAuthID?: string;
  httpEndpoint?: string;
  topicRoot?: string;
};

export const EDGE_NODE_PREVIEW_HOST = 'host';
export const EDGE_NODE_PREVIEW_BRANCH = 'branch';
export const EDGE_NODE_PREVIEW_BRANCH_OFFLINE = 'branch-offline';
export const EDGE_NODE_PREVIEW_DETAIL = 'detail';

export type EdgeNodePreviewMode =
  | typeof EDGE_NODE_PREVIEW_HOST
  | typeof EDGE_NODE_PREVIEW_BRANCH
  | typeof EDGE_NODE_PREVIEW_BRANCH_OFFLINE
  | typeof EDGE_NODE_PREVIEW_DETAIL;

export const parseEdgeNodePreviewMode = (value: string | null): EdgeNodePreviewMode | null => {
  if (
    value === EDGE_NODE_PREVIEW_HOST ||
    value === EDGE_NODE_PREVIEW_BRANCH ||
    value === EDGE_NODE_PREVIEW_BRANCH_OFFLINE ||
    value === EDGE_NODE_PREVIEW_DETAIL
  ) {
    return value;
  }
  return null;
};

export const PREVIEW_EDGE_NODES: PreviewClusterNode[] = [
  {
    id: 1,
    nodeName: 'SMT-Line-Edge',
    nodeID: 'edge-smt-01',
    status: 'online',
    lastSeenTime: Date.now() - 60_000,
    connectClientID: 1001,
    lastError: '',
    token:
      'preview-token-smt#eyJzeW5jTW9kZSI6ImVkZ2VFZGdlIiwiaHR0cEJhc2VVUkwiOiJodHRwOi8vMTkyLjE2OC4zMS4yMDU6ODE4OCIsImNsaWVudElEIjoiMTAwMSZob3N0LXN5bmMtMDEiLCJub2RlTmFtZSI6IlNNVC1MaW5lLUVkZ2UifQ==',
  },
  {
    id: 2,
    nodeName: 'Warehouse-Edge',
    nodeID: 'edge-wh-02',
    status: 'offline',
    lastSeenTime: Date.now() - 86_400_000,
    connectClientID: 1002,
    lastError: 'MQTT heartbeat timeout after 90s',
    token:
      'preview-token-wh#eyJzeW5jTW9kZSI6ImVkZ2VFZGdlIiwiaHR0cEJhc2VVUkwiOiJodHRwOi8vMTkyLjE2OC4zMS4yMDU6ODE4OCIsImNsaWVudElEIjoiMTAwMiZob3N0LXN5bmMtMDEiLCJub2RlTmFtZSI6IldhcmVob3VzZS1FZGdlIn0=',
  },
  {
    id: 3,
    nodeName: 'Quality-Lab-Edge',
    nodeID: 'edge-qa-03',
    status: 'connecting',
    lastSeenTime: Date.now() - 300_000,
    connectClientID: 1003,
    lastError: 'Cloud metadata sync failed: HTTP 503',
    token:
      'preview-token-qa#eyJzeW5jTW9kZSI6ImNsb3VkRWRnZSIsImh0dHBCYXNlVVJMIjoiaHR0cHM6Ly9jbG91ZC50aWVyMC5pbyIsImNsaWVudElEIjoiNDImY2xvdWQtc3luYy1rZXkiLCJub2RlTmFtZSI6IlF1YWxpdHktTGFiLUVkZ2UifQ==',
  },
];

export const PREVIEW_CLOUD_SYNC_CONNECTED: PreviewCloudSyncConfig = {
  configured: true,
  connectStatus: 'connected',
  edgeNodeName: 'Tier0 Edge Node',
  edgeNodeID: 'edge-local-01',
  clientID: 'edge-client-001',
  connectClientKey: 'edgetoken4',
  mqttAuthID: 'edgetoken4',
  syncMode: 'bidirectional',
  mountMode: 'root',
  lastConnectTime: Date.now() - 120_000,
  lastDisconnectTime: Date.now() - 3_600_000,
  lastSyncTime: Date.now() - 300_000,
  selectedRootPaths: ['v1', 'Plant', 'SMT-Area-1'],
  mqttBrokers: 'tcp://192.168.31.205:1983',
  topicRoot: 'cloudsync',
};

export const PREVIEW_CLOUD_SYNC_DISCONNECTED: PreviewCloudSyncConfig = {
  ...PREVIEW_CLOUD_SYNC_CONNECTED,
  connectStatus: 'disconnected',
  lastDisconnectTime: Date.now() - 60_000,
  lastError: 'Connection timeout while syncing UNS metadata',
};

export const PREVIEW_LAST_UNS_SYNC_LOG: CloudSyncLogItem = {
  id: 101,
  syncType: 'uns',
  createdTime: Date.now() - 300_000,
  status: 'success',
  direction: 'edge_to_cloud',
  syncMode: 'bidirectional',
  snapshotVersion: 12,
  summary: 'Synced 3 UNS nodes',
  environment: JSON.stringify({ nodeID: 'edge-smt-01', syncMode: 'bidirectional' }, null, 2),
  details: JSON.stringify(
    {
      createdFiles: 3,
      updatedFiles: 0,
      deletedFiles: 0,
      totalTopics: 3,
      createdNodes: [{ topic: 'v1/Plant/SMT-Area-1/Metric/cycle_time' }],
    },
    null,
    2
  ),
};

export const PREVIEW_LAST_USER_SYNC_LOG: CloudSyncLogItem = {
  id: 102,
  syncType: 'user',
  createdTime: Date.now() - 900_000,
  status: 'failed',
  direction: 'cloud_to_edge',
  syncMode: 'bidirectional',
  snapshotVersion: 8,
  summary: 'Created 2/ Skipped 1',
  environment: JSON.stringify(
    {
      source: 'cloud',
      operatorId: 1,
      operatorName: 'admin',
      durationMs: 1840,
    },
    null,
    2
  ),
  details: JSON.stringify(
    {
      createdUsers: [
        {
          username: 'tier0',
          userName: 'tier0',
          displayName: 'tier0',
          email: 'admin@freezonex.io',
          sourceRoleName: 'Operator',
          targetRoleName: 'Operator',
        },
        {
          username: 'Steven01',
          userName: 'Steven01',
          displayName: 'Steven',
          email: 'steven@gmail.com',
          sourceRoleName: 'Builder',
          targetRoleName: 'Operator',
        },
      ],
      skippedUsers: [{ username: 'guest-c', reason: 'role conflict' }],
      stats: { total: 3, created: 2, updated: 0, disabled: 0, deleted: 0, skipped: 1, conflicts: 0 },
      durationMs: 1840,
    },
    null,
    2
  ),
};

export const PREVIEW_USER_SYNC_LOGS: CloudSyncLogItem[] = [
  PREVIEW_LAST_USER_SYNC_LOG,
  {
    id: 105,
    syncType: 'user',
    createdTime: Date.now() - 7_200_000,
    status: 'success',
    direction: 'cloud_to_edge',
    syncMode: 'bidirectional',
    snapshotVersion: 7,
    summary: 'Imported 4 users',
    environment: JSON.stringify(
      {
        source: 'host',
        operatorId: 2,
        operatorName: 'builder',
        durationMs: 920,
      },
      null,
      2
    ),
    details: JSON.stringify({ durationMs: 920, stats: { total: 4, created: 4 } }, null, 2),
  },
];

export type PreviewSyncedAccount = {
  id: string;
  userId: string;
  preferredUsername: string;
  firstName: string;
  email: string;
  enabled: boolean;
  roleList: { roleName: string }[];
};

export const PREVIEW_SYNCED_ACCOUNTS: PreviewSyncedAccount[] = [
  {
    id: '1',
    userId: '1',
    preferredUsername: 'user01',
    firstName: 'user01',
    email: 'user01@example.com',
    enabled: true,
    roleList: [{ roleName: 'Operator' }],
  },
  {
    id: '2',
    userId: '2',
    preferredUsername: 'user02',
    firstName: 'user02',
    email: 'user02@example.com',
    enabled: true,
    roleList: [{ roleName: 'Builder' }],
  },
  {
    id: '3',
    userId: '3',
    preferredUsername: 'user03',
    firstName: 'user03',
    email: 'user03@example.com',
    enabled: true,
    roleList: [{ roleName: 'Operator' }],
  },
  {
    id: '4',
    userId: '4',
    preferredUsername: 'user03',
    firstName: 'user03',
    email: 'user03@example.com',
    enabled: false,
    roleList: [{ roleName: 'Operator' }],
  },
  {
    id: '5',
    userId: '5',
    preferredUsername: 'user03',
    firstName: 'user03',
    email: 'user03@example.com',
    enabled: true,
    roleList: [{ roleName: 'Operator' }],
  },
  {
    id: '6',
    userId: '6',
    preferredUsername: 'user03',
    firstName: 'user03',
    email: 'user03@example.com',
    enabled: false,
    roleList: [{ roleName: 'Operator' }],
  },
  {
    id: '7',
    userId: '7',
    preferredUsername: 'user03',
    firstName: 'user03',
    email: 'user03@example.com',
    enabled: false,
    roleList: [{ roleName: 'Operator' }],
  },
];

export const PREVIEW_SYNC_LOGS: CloudSyncLogItem[] = [
  PREVIEW_LAST_UNS_SYNC_LOG,
  {
    id: 103,
    syncType: 'uns',
    createdTime: Date.now() - 1_800_000,
    status: 'failed',
    direction: 'cloud_to_edge',
    syncMode: 'bidirectional',
    snapshotVersion: 11,
    errorMessage: 'Connection timeout',
    summary: 'Failed to pull UNS snapshot',
    environment: JSON.stringify({ nodeID: 'edge-smt-01', syncMode: 'bidirectional' }, null, 2),
    details: JSON.stringify({ failedCount: 2, errmsg: 'timeout' }, null, 2),
  },
  {
    id: 104,
    syncType: 'uns',
    createdTime: Date.now() - 3_600_000,
    status: 'success',
    direction: 'edge_to_cloud',
    syncMode: 'bidirectional',
    snapshotVersion: 10,
    summary: 'Synced 18 UNS nodes',
    environment: JSON.stringify({ nodeID: 'edge-smt-01', syncMode: 'bidirectional' }, null, 2),
    details: JSON.stringify(
      {
        createdFiles: 18,
        updatedFiles: 2,
        deletedFiles: 0,
        totalTopics: 20,
        createdNodes: [{ topic: 'v1/Plant/SMT-Area-1/State/running' }],
      },
      null,
      2
    ),
  },
];

export const edgeNodePreviewUrls = {
  host: '/edge-connection?tab=edgeNode&edgeNodePreview=host',
  branch: '/edge-connection?tab=edgeNode&edgeNodePreview=branch',
  branchOffline: '/edge-connection?tab=edgeNode&edgeNodePreview=branch-offline',
  detail: '/edge-connection/nodes/edge-smt-01?edgeNodePreview=detail',
} as const;
