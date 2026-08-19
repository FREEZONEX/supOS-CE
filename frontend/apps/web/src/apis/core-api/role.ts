import { iamApi } from './core-adapter';

// 新增角色
export const addRole = async (data?: Record<string, unknown>) => iamApi.post('/roles', data);

// 更新角色
export const putRole = async (data?: Record<string, unknown>) => iamApi.put(`/roles/${data?.id}`, data);

// 删除角色
export const deleteRole = async (roleId?: string) => iamApi.delete(`/roles/${roleId}`);
