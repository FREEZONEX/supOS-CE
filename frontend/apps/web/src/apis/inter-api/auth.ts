import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/auth';

const api = new ApiWrapper(baseUrl);

// 获取用户信息
export const getUserInfo = async (params?: Record<string, unknown>) => api.get('/user', { params });
export const logoutApi = async () => api.delete('/logout');
