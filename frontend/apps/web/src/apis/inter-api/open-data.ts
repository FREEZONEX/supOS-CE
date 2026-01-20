import { ApiWrapper } from '@/utils';

const baseUrl = '/inter-api/supos';
const api = new ApiWrapper(baseUrl);

// 获取密钥列表
export const querySecretKeyList = async () => api.get('/app/secretKeyList');

// 新增密钥
export const createSecretKey = async () => api.post('/app/secretKey');

// 删除密钥
export const deleteSecretKey = async (id: number) => api.delete(`/app/secretKey/${id}`);

// 更新密钥
export const updateSecretKey = async (data: any) => api.put('/app/secretKey', data);

// 获取第三方
export const queryMenuList = async () => api.get('/app/menu/list');
