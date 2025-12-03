import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos/uns';

const api = new ApiWrapper(baseUrl);

// 新增rule
export const addRule = async (data: any) => api.post('/alarm/rule', data);
// 编辑rule
export const editRule = async (data: any) => api.put('/alarm/rule', data);
// 删除rule
export const deleteRule = async (params: any) => api.delete(``, { params });
// 分页获取报警列表
export const getAlarmList = async (data: any) =>
  api.post('/alarm/pageList', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });
// 确认报警
export const confirmAlarm = async (data: any) => api.post('/alarm/confirm', data);
