import { ApiWrapper } from '@/utils';

const baseUrl = '/inter-api/supos/app';
const api = new ApiWrapper(baseUrl);

// 获取密钥列表
export const querySecretKeyList = async () => api.get('/secretKey/list');

// 新增密钥
export const createSecretKey = async () => api.post('/secretKey');

// 删除密钥
export const deleteSecretKey = async (id: number) => api.delete(`/secretKey/${id}`);

// 更新密钥
export const updateSecretKey = async (data: any) => api.put('/secretKey', data);
