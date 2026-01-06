import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/kong';

const api = new ApiWrapper(baseUrl);

// 获取kong所有路由
export const getKongRoutesApi = async (data: any) => api.get('/routeList', data);
