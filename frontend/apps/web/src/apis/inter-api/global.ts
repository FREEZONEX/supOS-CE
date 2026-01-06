import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos/global';
const api = new ApiWrapper(baseUrl);

// 获取导出结果
export const getGlobalExportRecords = async (params?: Record<string, unknown>) =>
  api.get('/user/getExportRecords', params);

// 确认导出记录
export const globalExportRecordConfirm = async (params?: Record<string, unknown>) =>
  api.post('/user/exportRecordConfirm', params);

// 导入
export const importGlobal = async (data: any) =>
  api.upload(`/data/import`, data, {
    method: 'post',
  });

// 导出
export const exportGlobal = async (data: any) => api.post('/data/export', data);

// 文件下载
export const downloadGlobalFile = async (params?: Record<string, unknown>) =>
  api.get('/file/download', { params, responseType: 'blob', [CustomAxiosConfigEnum.NoCode]: true });
