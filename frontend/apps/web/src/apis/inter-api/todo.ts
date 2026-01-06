import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos';

const api = new ApiWrapper(baseUrl);

// 获取系统模块列表
export const getModuleList = async (params?: Record<string, unknown>) => api.get('/todo/moduleList', { params });
// 待办已办列表
export const todoPageList = async (data: any) =>
  api.post('/todo/pageList', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });
