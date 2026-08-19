import { api } from './client';
import type { ModelStatus } from './models';

// bindingJson 结构与源端 InstanceBindingPayload 对齐（tier0-frontend anchor/types.ts）
export type SignedAxis = '+x' | '-x' | '+y' | '-y' | '+z' | '-z';

export interface SelectedObjectNode {
  nodeID: string;
  name: string;
  path: string;
}

export interface MotionMappingNode {
  nodeID: string;
  name: string;
  path: string;
  type: 'position' | 'rotation';
  value: string;
  axis: SignedAxis;
  unit: string;
  factor: number;
}

export interface InstanceDataTag {
  id: string;
  nodeID: string;
  name: string;
  path: string;
  x?: number;
  y?: number;
  z?: number;
  payload: string;
  unit?: string;
  createdAt?: string;
}

// 运动标签在 3D 上的可见性设置（对齐源端 MotionVisibilitySettings，随 bindingJson 持久化）
export interface MotionVisibilitySettings {
  mode: 'showAll' | 'hideAll' | 'custom';
  custom: Record<string, boolean>;
}

export interface InstanceBindingPayload {
  selectedObjects: SelectedObjectNode[];
  motionMappings: MotionMappingNode[];
  dataTags: InstanceDataTag[];
  motionTagVisibility?: MotionVisibilitySettings;
  [key: string]: unknown;
}

export const emptyBinding = (): InstanceBindingPayload => ({
  selectedObjects: [],
  motionMappings: [],
  dataTags: [],
});

export function parseBinding(raw: string | undefined): InstanceBindingPayload {
  if (!raw) return emptyBinding();
  try {
    const parsed = JSON.parse(raw) as Partial<InstanceBindingPayload>;
    return {
      ...parsed,
      selectedObjects: parsed.selectedObjects ?? [],
      motionMappings: parsed.motionMappings ?? [],
      dataTags: parsed.dataTags ?? [],
    };
  } catch {
    return emptyBinding();
  }
}

export interface InstanceInfo {
  id: number;
  modelId: number;
  modelName: string;
  modelOriginFile?: string;
  modelFileUrl: string;
  modelThumbnailUrl: string;
  modelNodesJson: string;
  name: string;
  unsNodeId: string;
  topic: string;
  bindingJson: string;
  height: number;
  status: ModelStatus;
  errorMsg: string;
  createdTime: string;
  updatedTime: string;
}

export interface InstanceListResult {
  list: InstanceInfo[];
  total: number;
  page: number;
  size: number;
}

const BASE = '/api/core/anchor';

export const listInstances = (params: { page?: number; size?: number; keyword?: string; modelId?: number }) =>
  api.get<InstanceListResult>(`${BASE}/instances`, params);

export const getInstance = (id: number | string) => api.get<InstanceInfo>(`${BASE}/instances/${id}`);

export const createInstance = (body: {
  modelId: number;
  name: string;
  unsNodeId?: string;
  topic?: string;
  bindingJson?: string;
  height?: number;
}) => api.post<InstanceInfo>(`${BASE}/instances`, body);

export const updateInstance = (
  id: number,
  body: { name?: string; unsNodeId?: string; topic?: string; bindingJson?: string; height?: number }
) => api.put<InstanceInfo>(`${BASE}/instances/${id}`, body);

export const deleteInstance = (id: number) => api.delete<{ id: number }>(`${BASE}/instances/${id}`);

export const getMqttCredentials = () =>
  api.post<{ username: string; password: string; clientId: string }>(`${BASE}/mqtt-credentials`);

// 二维码分享配置（免登录接口，供 /viewer 扫码页使用）
export interface QrConfig {
  instance: { id: number; name: string };
  model: { id: number; name: string; originFile: string; fileUrl: string; height: number };
  bindingJson: string;
  mqtt: {
    username: string;
    password: string;
    clientId: string;
    wsPort: number;
    wssPort: number;
    path: string;
    topic: string;
  };
}

export const qrConfigUrl = (id: number | string) => `${BASE}/instances/${id}/qr-config`;

export const fetchQrConfig = (configUrl: string) => api.get<QrConfig>(configUrl);
