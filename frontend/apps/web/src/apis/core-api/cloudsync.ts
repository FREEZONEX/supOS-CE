import { coreApi } from './core-adapter';
import { CustomAxiosConfigEnum } from '@/utils/request';

const silentErrorConfig = {
  [CustomAxiosConfigEnum.NoMessage]: true,
};

export interface CloudSyncConfigResp {
  configured: boolean;
  workspaceID?: number;
  mqttAuthID?: string;
  connectClientKey?: string;
  edgeNodeID?: string;
  edgeNodeName?: string;
  hasToken?: boolean;
  tokenLength?: number;
  clientID?: string;
  username?: string;
  selectedRootNodeIDs?: string[];
  selectedRootPaths?: string[];
  flattenedRootNames?: string[];
  scopeRevision?: number;
  mountMode?: string;
  syncMode?: string;
  syncMetricEnabled?: boolean;
  desiredConnected?: boolean;
  connectStatus?: string;
  httpEndpoint?: string;
  mqttBrokers?: string;
  topicRoot?: string;
  lastConnectTime?: number;
  lastDisconnectTime?: number;
  lastSyncTime?: number;
  lastError?: string;
  createAt?: number;
  updateAt?: number;
}

export type ClusterNode = {
  id: number;
  nodeKey?: string;
  nodeName?: string;
  nodeID?: string;
  name?: string;
  role?: string;
  status?: string | number;
  description?: string;
  endpoint?: string;
  lastSeen?: number;
  lastSeenTime?: number;
  connectClientID?: number;
  lastError?: string;
};

export type ClusterScope = {
  id: number;
  nodeId: number;
  scopeType: string;
  scopeKey: string;
  direction: string;
  enabled: boolean;
};

export type ClusterOutbox = {
  id: number;
  eventId: string;
  eventType: string;
  aggregateType: string;
  aggregateId: string;
  status: string;
  attempts: number;
  lastError?: string;
};

export type CloudSyncLogItem = {
  id: number;
  syncType?: string;
  direction?: string;
  status?: string;
  syncMode?: string;
  connectClientKey?: string;
  nodeID?: string;
  nodeName?: string;
  snapshotVersion?: number;
  snapshotHash?: string;
  resultCode?: string;
  summary?: string;
  errorMessage?: string;
  environment?: string;
  details?: string;
  createdTime?: number;
  updatedTime?: number;
};

export type ClusterSyncLogListParams = {
  page?: number;
  size?: number;
  syncType?: string;
  direction?: string;
  status?: string;
  connectClientKey?: string;
  nodeID?: string;
};

export type ClusterTokenCreateReq = {
  tokenName: string;
  description?: string;
};

export type ClusterTokenCreateResp = {
  token?: string;
  record?: ClusterNode;
  client?: Record<string, unknown>;
};

export type ClusterHandshakeReq = {
  token: string;
  nodeKey: string;
  nodeName?: string;
  role?: string;
  endpoint?: string;
  capabilitiesJson?: string;
};

export type ClusterHandshakeResp = {
  protocolVersion?: string;
  [key: string]: unknown;
};

type ClusterListResp<T> = {
  list?: T[];
  total?: number;
};

export interface CloudSyncConnectReq {
  token?: string;
  selectedRootNodeIDs?: string[];
  expectedScopeRevision: number;
  syncMetricEnabled?: boolean;
}

export interface CloudSyncTestConnectReq {
  token?: string;
}

export interface CloudSyncTestConnectResp {
  ok: boolean;
  workspaceID?: number;
  mqttAuthID?: string;
  connectClientKey?: string;
  clientID?: string;
  username?: string;
  httpEndpoint?: string;
  mqttBrokers?: string;
  topicRoot?: string;
  syncMode?: string;
  message?: string;
  error?: string;
  debug?: string;
}

export const getCloudSyncConfig = async (): Promise<CloudSyncConfigResp> => {
  return coreApi.get('/cluster/config');
};

export const connectCloudSync = async (data: CloudSyncConnectReq): Promise<CloudSyncConfigResp> => {
  return coreApi.post('/cluster/connect', data, silentErrorConfig);
};

export const disconnectCloudSync = async (): Promise<CloudSyncConfigResp> => coreApi.post('/cluster/disconnect', {});

export const testCloudSyncConnect = async (data: CloudSyncTestConnectReq): Promise<CloudSyncTestConnectResp> =>
  coreApi.post('/cluster/test-connect', data, silentErrorConfig);

export const forceSyncCloudSync = async (): Promise<CloudSyncConfigResp> => coreApi.post('/cluster/force-sync', {});

export const userSyncCloudSync = async (): Promise<CloudSyncLogItem> => coreApi.post('/cluster/user-sync', {});

export const getClusterNodes = async (): Promise<ClusterListResp<ClusterNode>> => coreApi.get('/cluster/nodes');

export const getClusterScopes = async (): Promise<ClusterListResp<ClusterScope>> => coreApi.get('/cluster/scopes');

export const getClusterOutbox = async (): Promise<ClusterListResp<ClusterOutbox>> => coreApi.get('/cluster/outbox');

export const getClusterSyncLogs = async (
  params: ClusterSyncLogListParams
): Promise<ClusterListResp<CloudSyncLogItem> & { page?: number; size?: number }> =>
  coreApi.get('/cluster/sync-logs', { params });

export const createClusterToken = async (data: ClusterTokenCreateReq): Promise<ClusterTokenCreateResp> =>
  coreApi.post('/cluster/tokens', data);

export const handshakeCluster = async (data: ClusterHandshakeReq): Promise<ClusterHandshakeResp> =>
  coreApi.post('/cluster/handshake', data);
