import { api } from './client';

export interface UnsNode {
  id: string;
  parentId?: string;
  name: string;
  path?: string;
  namespace?: string;
  alias?: string;
  type?: string | number;
  hasChildren?: boolean;
  countChildren?: number;
  fields?: { name: string; type?: string; unit?: string }[];
  lastPayload?: string | Record<string, unknown>;
  children?: UnsNode[];
  [key: string]: unknown;
}

// 判定文件型节点（可订阅的 topic）；与 web mapUnsNode 的 pathType 判定一致
export const isTopicNode = (node: UnsNode) =>
  String(node.type ?? '').toLowerCase() === 'file' || Number(node.type) === 2;

export const listUnsNodes = (params: { parentId?: string | number; parentIdSet?: boolean; keyword?: string }) =>
  api.get<UnsNode[] | { list: UnsNode[] }>('/api/core/uns/nodes', {
    parentId: params.parentId !== undefined ? String(params.parentId) : undefined,
    parentIdSet: params.parentIdSet ? 'true' : undefined,
    keyword: params.keyword,
  });

export const getUnsNode = (nodeId: string | number) => api.get<UnsNode>(`/api/core/uns/nodes/${nodeId}`);

export const normalizeNodeList = (data: UnsNode[] | { list: UnsNode[] } | null | undefined): UnsNode[] => {
  if (!data) return [];
  if (Array.isArray(data)) return data;
  return data.list ?? [];
};
