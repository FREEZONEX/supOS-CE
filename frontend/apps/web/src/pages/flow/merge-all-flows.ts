import { getFlowAndGroupList, sortFlowListItems } from '@/apis/core-api/flow';
import { getEventFlowAndGroupList } from '@/apis/core-api/event-flow';

export const flowListRowKey = (item: { id?: number | string; flowKind?: string }) =>
  `${item.flowKind || 'source'}-${item.id}`;

const withoutPagination = (params?: Record<string, unknown>) => {
  if (!params) return undefined;
  const next = { ...params };
  delete next.page;
  delete next.pageNo;
  delete next.pageSize;
  return next;
};

const paginateMergedFlows = (list: any[], params?: Record<string, unknown>) => {
  const hasPagination = params?.pageNo !== undefined || params?.page !== undefined || params?.pageSize !== undefined;
  if (!hasPagination) {
    return {
      data: list,
      total: list.length,
      pageNo: 1,
      pageSize: Math.max(list.length, 20),
    };
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

export const fetchMergedAllFlows = async (params?: Record<string, unknown>) => {
  const listParams = withoutPagination(params);
  const [sourceResp, eventResp] = await Promise.all([
    getFlowAndGroupList(listParams),
    getEventFlowAndGroupList(listParams),
  ]);
  const merged = [
    ...(sourceResp?.data || []).map((item: any) => ({ ...item, flowKind: 'source' as const })),
    ...(eventResp?.data || []).map((item: any) => ({ ...item, flowKind: 'event' as const })),
  ];
  const sorted = sortFlowListItems(merged, params);
  return paginateMergedFlows(sorted, params);
};
