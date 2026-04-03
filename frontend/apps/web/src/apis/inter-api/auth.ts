import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/auth';

const api = new ApiWrapper(baseUrl);

export const loginApi = async (data?: Record<string, unknown>) => api.post('/login', data);
// 获取用户信息
export const getUserInfo = async (params?: Record<string, unknown>) => api.get('/user', { params });
export const logoutApi = async () => api.delete('/logout');
