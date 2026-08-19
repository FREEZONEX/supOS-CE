import { coreApi } from './core-adapter';

export interface AuditLogPageParams {
  pageNo?: number;
  pageSize?: number;
  scopeType?: number;
  scopeId?: number;
  resType?: string;
  businessType?: string;
  operatorKeyword?: string;
  code?: number;
  startTime?: string;
  endTime?: string;
}

export interface AuditLogItem {
  id: number;
  requestId?: string;
  operatorId?: string;
  operatorName?: string;
  operatorEmail?: string;
  scopeType: number;
  scopeId?: number;
  scopeName?: string;
  resType: string;
  resId?: string;
  resName?: string;
  businessType: string;
  code: number;
  resultMsg?: string;
  detailJson?: string;
  ip?: string;
  userAgent?: string;
  createdAt: string;
}

export interface AuditLogPageResp {
  pageNo: number;
  pageSize: number;
  total: number;
  data: AuditLogItem[];
}

export const getAuditLogPage = async (data?: AuditLogPageParams): Promise<AuditLogPageResp> =>
  coreApi.post('/audit-logs/page', data);
export const getAuditLogDetail = async (id: number | string): Promise<AuditLogItem> => coreApi.get(`/audit-logs/${id}`);
