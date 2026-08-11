import { CustomAxiosConfigEnum } from '@/utils/request';
import { coreApi } from './core-adapter';
import type { CloudSyncConfigResp } from './cloudsync';

export type HomeClusterConfig = CloudSyncConfigResp & {
  enabled?: boolean;
  role?: string;
  nodeKey?: string;
  nodeName?: string;
  isHost?: boolean;
};

export type HomeClusterNode = {
  id: number;
  nodeKey: string;
  nodeName: string;
  role: string;
  status: string;
  endpoint?: string;
  lastSeenTime?: number;
  syncUserName?: string;
  syncUserEmail?: string;
  dispatchUserName?: string;
  dispatchUserEmail?: string;
};

export type HomeRecentItem = {
  key: string;
  resourceKey: string;
  resourceName: string;
  resourceType: string;
  businessType?: string;
  routePath: string;
  icon?: string;
  operatedAt?: number;
};

const silentConfig = {
  [CustomAxiosConfigEnum.NoMessage]: true,
};

export const getHomeClusterConfig = async (): Promise<HomeClusterConfig | undefined> => {
  try {
    return await coreApi.get('/cluster/config', silentConfig);
  } catch {
    return undefined;
  }
};

export const getHomeClusterNodes = async (): Promise<HomeClusterNode[]> => {
  try {
    const resp = await coreApi.get('/cluster/nodes', silentConfig);
    return resp?.list || [];
  } catch {
    return [];
  }
};

export const getHomeRecents = async (): Promise<HomeRecentItem[]> => {
  try {
    const resp = await coreApi.get('/home/recents', {
      params: { limit: 5, isShowInRecent: 1 },
      ...silentConfig,
    });
    return Array.isArray(resp?.list) ? resp.list : [];
  } catch {
    return [];
  }
};
