/* eslint-disable @typescript-eslint/no-unused-vars */
import { coreApi, mapFlow } from './core-adapter';

const toNodeID = (value: any) => {
  const nodeID = Number(value);
  return Number.isFinite(nodeID) && nodeID > 0 ? nodeID : 0;
};

const normalizeUnsNodeIDs = (value: any) => {
  const values = Array.isArray(value) ? value : value ? [value] : [];
  return Array.from(new Set(values.map(toNodeID).filter(Boolean)));
};

const paginateFlowItems = (list: any[], params?: Record<string, unknown>) => {
  const hasPagination = params?.pageNo !== undefined || params?.page !== undefined || params?.pageSize !== undefined;
  if (!hasPagination) {
    return { data: list, total: list.length, pageNo: 1, pageSize: Math.max(list.length, 20) };
  }
  const pageSize = Math.max(Number(params?.pageSize) || 20, 1);
  const requestedPageNo = Math.max(Number(params?.pageNo ?? params?.page) || 1, 1);
  const maxPageNo = Math.max(Math.ceil(list.length / pageSize), 1);
  const pageNo = Math.min(requestedPageNo, maxPageNo);
  const start = (pageNo - 1) * pageSize;
  return {
    data: list.slice(start, start + pageSize),
    total: list.length,
    pageNo,
    pageSize,
  };
};

const toPayload = (data: any) => ({
  parentId: Number(data?.parentId || data?.groupId || 0),
  flowType: 'event',
  nodeType: data?.nodeType || 'flow',
  name: data?.flowName || data?.name || 'untitled-event-flow',
  description: data?.description || '',
  template: data?.template || data?.flowTemplate || 'node-red',
  unsNodeIds: normalizeUnsNodeIDs(data?.unsNodeIds),
});

const listFlows = async (params?: Record<string, unknown>) => {
  const resp = await coreApi.get('/flows', {
    params: {
      flowType: 'event',
      parentId: params?.groupId,
      keyword: params?.k || params?.keyword,
    },
  });
  const category = String(params?.category || '').toLowerCase();
  const list = sortFlowItems(
    (resp?.list || [])
      .map(mapFlow)
      .filter((item: any) => !category || category === 'all' || item.category === category),
    params
  );
  return paginateFlowItems(list, params);
};

const sortValue = (item: any, key: string) => {
  if (key === 'flowName') return item.name || item.flowName || '';
  if (key === 'createAt' || key === 'createdTime' || key === 'createTime') {
    return Number(item.createAt || item.createdTime || item.createTime || 0);
  }
  return item?.[key] ?? '';
};

const isPinnedFlow = (item: any) =>
  item?.pinned === true ||
  Number(item?.sort || 0) === 1 ||
  Number(item?.isFavorite || 0) === 1 ||
  Number(item?.sortKey || 0) > 0;

const comparePinnedFlow = (a: any, b: any) => {
  const aPinned = isPinnedFlow(a);
  const bPinned = isPinnedFlow(b);
  if (aPinned !== bPinned) return aPinned ? -1 : 1;
  if (aPinned && bPinned) return Number(b?.sortKey || 0) - Number(a?.sortKey || 0);
  return 0;
};

const sortFlowItems = (items: any[], params?: Record<string, any>) => {
  const orderCode = String(params?.orderCode || params?.sortData?.[0]?.orderCode || '').trim();
  const ascValue = params?.isAsc ?? params?.sortData?.[0]?.isAsc;
  const isAsc = ascValue === true || ascValue === 1 || ascValue === 'true' || ascValue === '1';
  return [...items].sort((a, b) => {
    const pinnedResult = comparePinnedFlow(a, b);
    if (pinnedResult !== 0) return pinnedResult;
    if (!orderCode) return 0;
    const left = sortValue(a, orderCode);
    const right = sortValue(b, orderCode);
    const result =
      typeof left === 'number' && typeof right === 'number'
        ? left - right
        : String(left).localeCompare(String(right), undefined, { sensitivity: 'base' });
    return isAsc ? result : -result;
  });
};

export const addFlow = async (data: any) => mapFlow(await coreApi.post('/flows', toPayload(data)));
export const getFlowDetail = async (id: string | number) => mapFlow(await coreApi.get(`/flows/${id}`));
export const copyFlow = async (data: any) =>
  addFlow({ ...data, flowName: `${data?.flowName || data?.name || 'event-flow'} copy` });
export const editFlow = async (data: any) => mapFlow(await coreApi.put(`/flows/${data?.id}`, toPayload(data)));
export const deleteFlow = async (id: string) => coreApi.delete(`/flows/${id}`);
const toFlowDataPayload = (data: any) => ({
  flowData: JSON.stringify(data?.flows || data?.flowData || data || []),
});
export const deployFlow = async (data: any) => coreApi.post(`/flows/${data?.id}/deploy`, toFlowDataPayload(data));
export const updateFlowStatus = async (id: string | number, status: 'disabled' | 'deployed') =>
  mapFlow(await coreApi.put(`/flows/${id}/status`, { status }));
export const getVersionFlow = async () => coreApi.get('/flows/version', { params: { flowType: 'event' } });
export const saveFlow = async (data: any) => coreApi.put(`/flows/${data?.id}/data`, toFlowDataPayload(data));
export type FlowVersionItem = {
  id: number | string;
  versionId?: number | string;
  versionID?: number | string;
  versionName: string;
  versionType: number;
  versionTypeLabel?: string;
  description?: string;
  createdTime?: number | string;
  updatedTime?: number | string;
  createdBy?: number | string;
  operatorName?: string;
  isCurrent?: number | string;
};
export const createFlowVersion = async (
  id: string | number,
  data: { versionName?: string; description?: string; flowData?: string; flowsJson?: string }
) => coreApi.post(`/flows/${id}/version`, data);
export const listFlowVersions = async (id: string | number, params: { pageNo?: number; pageSize?: number } = {}) =>
  coreApi.post(`/flows/${id}/version/list`, params);
export const restoreFlowVersion = async (id: string | number, versionId: string | number) =>
  coreApi.post(`/flows/${id}/version/${versionId}/restore`, {});
export const saveAsFlowVersion = async (
  id: string | number,
  versionId: string | number,
  data: { name?: string; flowName?: string; description?: string }
) => coreApi.post(`/flows/${id}/version/${versionId}/save-as`, data);
export const flowPage = async (params?: Record<string, unknown>) => listFlows(params);
export const processList = async () => ({ data: [], total: 0, pageNo: 1, pageSize: 20 });
export const markFlow = async (id: string) => mapFlow(await coreApi.put(`/flows/${id}/mark`, { pinned: true }));
export const unmarkFlow = async (id: string) => mapFlow(await coreApi.put(`/flows/${id}/mark`, { pinned: false }));
export const getEventFlowAndGroupList = async (params?: Record<string, unknown>, _config?: any) => listFlows(params);
