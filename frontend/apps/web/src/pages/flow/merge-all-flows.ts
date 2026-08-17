import { getAllFlowAndGroupList } from '@/apis/core-api/flow';

export const flowListRowKey = (item: { id?: number | string; flowKind?: string }) =>
  `${item.flowKind || 'source'}-${item.id}`;

export const fetchMergedAllFlows = async (params?: Record<string, unknown>) => getAllFlowAndGroupList(params);
