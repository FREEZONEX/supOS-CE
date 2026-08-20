import { api } from './client';

// 对齐源端 ModelStatusEnum（tier0-frontend anchor/types.ts）
export enum ModelStatus {
  Unknown = 0,
  Parsing = 1,
  Ready = 2,
  Error = 3,
}

export interface ModelInfo {
  id: number;
  name: string;
  originFile: string;
  fileAssetId: number;
  thumbnailAssetId: number;
  fileSize: number;
  nodesJson: string;
  status: ModelStatus;
  errorMsg: string;
  createdBy: number;
  createdTime: string;
  updatedTime: string;
  fileUrl: string;
  thumbnailUrl: string;
}

export interface ModelListResult {
  list: ModelInfo[];
  total: number;
  page: number;
  size: number;
}

interface UploadedAsset {
  fileId: number;
  originalName?: string;
  size?: number;
  [key: string]: unknown;
}

const BASE = '/api/core/anchor';

export const listModels = (params: { page?: number; size?: number; keyword?: string }) =>
  api.get<ModelListResult>(`${BASE}/models`, params);

export const getModel = (id: number | string) => api.get<ModelInfo>(`${BASE}/models/${id}`);

export const createModel = (body: {
  name: string;
  originFile: string;
  fileAssetId: number;
  thumbnailAssetId?: number;
  fileSize: number;
}) => api.post<ModelInfo>(`${BASE}/models`, body);

export const updateModel = (
  id: number,
  body: {
    name?: string;
    thumbnailAssetId?: number;
    nodesJson?: string;
    // Replace Model
    fileAssetId?: number;
    originFile?: string;
    fileSize?: number;
  }
) => api.put<ModelInfo>(`${BASE}/models/${id}`, body);

export const deleteModel = (id: number) => api.delete<{ id: number }>(`${BASE}/models/${id}`);

// 一键创建示例（服务端完成 UNS 节点 + mock 数据流 + 模型 + 实例，对齐云端 modelCreateDemo）
export const createDemoModel = () =>
  api.post<ModelInfo & { demoInstanceId?: number; demoTopic?: string }>(`${BASE}/models/demo`);

export async function uploadAsset(file: File): Promise<UploadedAsset> {
  const form = new FormData();
  form.append('files', file);
  const data = await api.upload<{ list: UploadedAsset[]; total: number }>('/api/core/assets', form);
  const first = data?.list?.[0];
  if (!first?.fileId) {
    throw new Error('asset upload failed: empty response');
  }
  return first;
}
