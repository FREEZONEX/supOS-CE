import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/resource';

const api = new ApiWrapper(baseUrl);

// 获取路由资源
export const getRoutesResourceApi = async () => api.get('');

// 新增修改 单个菜单内容
export const postResourceApi = async (data: any) => api.post('', data);

// 删除资源
export const deleteResourceApi = async (id: string) => api.delete(`/${id}`);

// 批量删除资源
export const batchDeleteResourceApi = async (data: string[]) => api.delete(`/batch`, undefined, data);

// 批量修改资源
export const batchEditResourceApi = async (data: any[]) => api.put('/batch', data);
