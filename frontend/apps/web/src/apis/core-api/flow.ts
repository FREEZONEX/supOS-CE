/* eslint-disable @typescript-eslint/no-unused-vars */
import { coreApi, mapFlow } from './core-adapter';

const lastPathPart = (value: any) =>
  String(value || '')
    .split('/')
    .filter(Boolean)
    .pop() || '';

const namespaceFromData = (data: any) =>
  String(data?.path || data?.namespace || '')
    .split('/')
    .filter(Boolean)
    .join('/');

const flowNameFromData = (data: any) =>
  data?.flowName ||
  data?.name ||
  namespaceFromData(data) ||
  lastPathPart(data?.unsAlias || data?.alias) ||
  'untitled-flow';

const toNodeID = (value: any) => {
  const nodeID = Number(value);
  return Number.isFinite(nodeID) && nodeID > 0 ? nodeID : 0;
};

const normalizeUnsNodeIDs = (value: any) => {
  const values = Array.isArray(value) ? value : value ? [value] : [];
  return Array.from(new Set(values.map(toNodeID).filter(Boolean)));
};

const normalizeMockFields = (value: any) => {
  const values = Array.isArray(value) ? value : [];
  return values
    .map((item: any) => ({
      name: String(item?.name || item?.fieldName || '').trim(),
      type: String(item?.type || item?.dataType || '').trim(),
    }))
    .filter((item) => item.name);
};

const explicitUnsNodeID = (data: any) => toNodeID(data?.unsId ?? data?.unsNodeId ?? data?.nodeId);

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

const toPayload = (data: any) => {
  const mockData = Boolean(data?.mockData);
  return {
    sourceId: toNodeID(data?.sourceId),
    parentId: Number(data?.parentId || data?.groupId || 0),
    flowType: 'source',
    nodeType: data?.nodeType || 'flow',
    name: flowNameFromData(data),
    description: data?.description || '',
    template: data?.template || data?.flowTemplate || 'node-red',
    unsNodeIds: normalizeUnsNodeIDs(data?.unsNodeIds),
    ...(mockData
      ? {
          mockData,
          mockTopic: data?.mockTopic || namespaceFromData(data) || data?.unsAlias || data?.alias || '',
          mockFields: normalizeMockFields(data?.mockFields || data?.fields),
          mockTriggerMode: data?.mockTriggerMode || data?.triggerMode || data?.injectMode || 'auto',
        }
      : {}),
  };
};

const listFlows = async (params?: Record<string, unknown>) => {
  const resp = await coreApi.get('/flows', {
    params: {
      flowType: 'source',
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
  if (key === 'flowName' || key === 'name') return item.name || item.flowName || '';
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

export const sortFlowListItems = sortFlowItems;

export const createFlowRowSorter = (key: string) => (a: any, b: any) => {
  const pinnedResult = comparePinnedFlow(a, b);
  if (pinnedResult !== 0) return pinnedResult;
  const left = sortValue(a, key);
  const right = sortValue(b, key);
  if (typeof left === 'number' && typeof right === 'number') return left - right;
  return String(left).localeCompare(String(right), undefined, { sensitivity: 'base' });
};

const resolveUnsNode = async (unsAlias?: string) => {
  const alias = String(unsAlias || '').trim();
  if (!alias) return null;
  const resp = await coreApi.get('/uns/nodes', { params: { keyword: alias } });
  const list = resp?.list || resp?.nodes || [];
  return (
    list.find((item: any) => item.alias === alias || item.namespace === alias || item.path === alias) ||
    list.find((item: any) => String(item.namespace || '').endsWith(`/${alias}`) || item.name === alias) ||
    null
  );
};

const sourceFlows = async (keyword?: string) => {
  const result = await listFlows(keyword ? { keyword } : undefined);
  return result.data || [];
};

export const addFlow = async (data: any) => mapFlow(await coreApi.post('/flows', toPayload(data)));
export const createFlow = async (data: any) => {
  const name = flowNameFromData(data);
  const inputNodeIDs = normalizeUnsNodeIDs(data?.unsNodeIds);
  const nodeID = explicitUnsNodeID(data) || inputNodeIDs[0] || 0;
  const unsNode = nodeID ? null : await resolveUnsNode(data?.unsAlias || data?.alias);
  const bindNodeID = nodeID || toNodeID(unsNode?.id);
  if (bindNodeID) {
    const boundFlow = (await sourceFlows()).find((item: any) =>
      item.unsNodeIds?.map(String).includes(String(bindNodeID))
    );
    if (boundFlow) {
      if ((boundFlow.flowName || boundFlow.name) !== name) {
        return editFlow({ ...boundFlow, flowName: name, name, unsNodeIds: boundFlow.unsNodeIds });
      }
      return boundFlow;
    }
  }
  const existing = (await sourceFlows(name)).find((item: any) => item.flowName === name || item.name === name);
  if (existing) {
    if (bindNodeID && !existing.unsNodeIds?.map(String).includes(String(bindNodeID))) {
      return bindFlowForUns({ flowId: existing.id, unsId: bindNodeID, unsAlias: data?.unsAlias || data?.alias });
    }
    return existing;
  }
  return addFlow({ ...data, flowName: name, unsNodeIds: bindNodeID ? [bindNodeID] : inputNodeIDs });
};
export const getFlowDetail = async (id: string | number) => mapFlow(await coreApi.get(`/flows/${id}`));
export const copyFlow = async (data: any) =>
  addFlow({ ...data, flowName: `${data?.flowName || data?.name || 'flow'} copy` });
export const editFlow = async (data: any) => mapFlow(await coreApi.put(`/flows/${data?.id}`, toPayload(data)));
export const deleteFlow = async (id: string) => coreApi.delete(`/flows/${id}`);
const toFlowDataPayload = (data: any) => ({
  flowData: JSON.stringify(data?.flows || data?.flowData || data || []),
});
export const deployFlow = async (data: any) => coreApi.post(`/flows/${data?.id}/deploy`, toFlowDataPayload(data));
export const updateFlowStatus = async (id: string | number, status: 'disabled' | 'deployed') =>
  mapFlow(await coreApi.put(`/flows/${id}/status`, { status }));
export const getVersionFlow = async () => coreApi.get('/flows/version', { params: { flowType: 'source' } });
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
export const goFlow = async (alias?: string, unsId?: string | number) => {
  const nodeId = unsId || (await resolveUnsNode(alias))?.id;
  if (!nodeId) return {};
  return (await sourceFlows()).find((item: any) => item.unsNodeIds?.map(String).includes(String(nodeId))) || {};
};
export const flowPage = async (params?: Record<string, unknown>) => listFlows(params);
export const processList = async () => ({ data: [], total: 0, pageNo: 1, pageSize: 20 });
export const markFlow = async (id: string) => mapFlow(await coreApi.put(`/flows/${id}/mark`, { pinned: true }));
export const unmarkFlow = async (id: string) => mapFlow(await coreApi.put(`/flows/${id}/mark`, { pinned: false }));
export const bindFlowForUns = async (params: any) => {
  const flowId = params?.flowId || params?.id;
  const nodeId = explicitUnsNodeID(params) || (await resolveUnsNode(params?.unsAlias || params?.alias))?.id;
  if (!flowId || !nodeId) return {};
  const detail = await getFlowDetail(flowId);
  const unsNodeIds = Array.from(new Set([...(detail?.unsNodeIds || []).map(Number), Number(nodeId)])).filter(Boolean);
  return editFlow({ ...detail, unsNodeIds });
};
export const getFlowAndGroupList = async (params?: Record<string, unknown>, _config?: any) => listFlows(params);
