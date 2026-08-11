/* eslint-disable @typescript-eslint/no-unused-vars */
import { iamApi, toLegacyResources } from './core-adapter';

// 获取路由资源
export const getRoutesResourceApi = async (): Promise<any> => {
  const resources = await iamApi.get('/resources');
  return toLegacyResources(resources || []);
};

export const getMenuResourceApi = async (): Promise<any> => {
  const resources = await iamApi.get('/menus');
  return toLegacyResources(resources || [], { includeActionButtons: false, includeHiddenMenus: true });
};

export const getResourceActionsApi = async (resourceId: string) =>
  iamApi.get(`/resources/${encodeURIComponent(resourceId)}/actions`);

export const postResourceActionApi = async (resourceId: string, data: any) =>
  iamApi.post(`/resources/${encodeURIComponent(resourceId)}/actions`, data);

export const putResourceActionApi = async (resourceId: string, actionId: string, data: any) =>
  iamApi.put(`/resources/${encodeURIComponent(resourceId)}/actions/${encodeURIComponent(actionId)}`, data);

export const deleteResourceActionApi = async (resourceId: string, actionId: string) =>
  iamApi.delete(`/resources/${encodeURIComponent(resourceId)}/actions/${encodeURIComponent(actionId)}`);

const normalizeResourceId = (id?: string | number) => String(id || '').replace(/^resource:/, '');

const toResourcePayload = (data: any) => ({
  parentId: normalizeResourceId(data?.parentId) ? Number(normalizeResourceId(data?.parentId)) : 0,
  resourceKey: data?.resourceKey || data?.code,
  name: data?.name,
  type: Number(data?.type || 2),
  routePath: data?.routePath || data?.url || '',
  urlType: Number(data?.urlType || 1),
  openType: Number(data?.openType || 0),
  icon: data?.icon || '',
  sort: Number(data?.sort || 0),
  enabled: data?.enabled ?? (data?.enable === false ? 0 : 1),
});

// 新增修改 单个菜单内容
export const saveResourceApi = async (data: any) => {
  const id = normalizeResourceId(data?.coreResourceId || data?.resourceId || data?.id);
  const payload = toResourcePayload(data);
  if (id) {
    return iamApi.put(`/resources/${encodeURIComponent(id)}`, payload);
  }
  return iamApi.post('/resources', payload);
};

export const postResourceApi = saveResourceApi;

// 删除资源
export const deleteResourceApi = async (_id: string) =>
  Promise.reject(new Error('resource delete is not supported yet'));

// 批量删除资源
export const batchDeleteResourceApi = async (_data: string[]) =>
  Promise.reject(new Error('resource batch delete is not supported yet'));

// 批量修改资源
export const batchEditResourceApi = async (_data: any[]) =>
  Promise.reject(new Error('resource batch edit is not supported yet'));
