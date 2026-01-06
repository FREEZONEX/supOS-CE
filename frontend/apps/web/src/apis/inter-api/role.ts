import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/role';

const api = new ApiWrapper(baseUrl);

// 新增角色
export const addRole = async (data?: Record<string, unknown>) => api.post('', data);
// 更新角色
export const putRole = async (data?: Record<string, unknown>) => api.put('', data);
// 删除角色
export const deleteRole = async (roleId: string) => api.delete(`/${roleId}`);
