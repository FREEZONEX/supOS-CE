/* eslint-disable @typescript-eslint/no-unused-vars */
import type { ProjectMember, ProjectRole } from '@/pages/project/types';
import { CustomAxiosConfigEnum } from '@/utils/request';
import { coreApi, mapFlow } from './core-adapter';

export interface ProjectListParams {
  k?: string;
  pageNo?: number;
  pageSize?: number;
}

export interface ProjectPermissionResp {
  projectId: number;
  isMember: boolean;
  canView: boolean;
  canEdit: boolean;
}

export type ProjectRoleOption = {
  label: string;
  value: string;
};

export interface ProjectRoleListParams {
  pageNo?: number;
  pageSize?: number;
}

const normalizeProject = (item: any) => {
  const source = item?.project || item || {};
  return {
    ...source,
    id: String(source.id || ''),
    displayName: source.displayName || source.name,
    createdAt: source.createdAt || source.createdTime,
    updatedAt: source.updatedAt || source.updatedTime,
  };
};

export const getProjectList = async (params?: ProjectListParams) => {
  const resp = await coreApi.get('/projects', { params: { keyword: params?.k } });
  return {
    data: (resp?.list || []).map(normalizeProject),
    total: resp?.total || 0,
    pageNo: params?.pageNo || 1,
    pageSize: params?.pageSize || 20,
  };
};

export const createProject = async (data?: Record<string, unknown>) =>
  normalizeProject(await coreApi.post('/projects', data));
export const getProject = async (id: string) =>
  normalizeProject(await coreApi.get(`/projects/${id}`, { [CustomAxiosConfigEnum.NoMessage]: true }));
export const updateProject = async (id: string, data?: Record<string, unknown>) =>
  normalizeProject(await coreApi.put(`/projects/${id}`, data));
export const deleteProject = async (id: string) => coreApi.delete(`/projects/${id}`);

export const updateProjectStatus = async (_projectId: string, _data?: Record<string, unknown>) => ({});

export const getProjectPermission = async (projectId: number): Promise<ProjectPermissionResp> => ({
  projectId,
  isMember: true,
  canView: true,
  canEdit: true,
});

export const getProjectApps = async (projectId: string) => coreApi.get(`/projects/${projectId}/apps`);
export const getProjectAppDetail = async (projectId: string, appId: string) =>
  coreApi.get(`/projects/${projectId}/apps/${appId}`);
export const importProjectApp = async (projectId: string, data?: Record<string, unknown>) =>
  coreApi.post(`/projects/${projectId}/apps`, data);
export const createManualProjectApp = async (projectId: string, data?: Record<string, unknown>) =>
  coreApi.post(`/projects/${projectId}/apps/manual`, data);
export const updateProjectApp = async (projectId: string, appId: string, data?: Record<string, unknown>) =>
  coreApi.put(`/projects/${projectId}/apps/${appId}`, data);
export const replaceProjectApp = async (projectId: string, appId: string, data?: Record<string, unknown>) =>
  coreApi.post(`/projects/${projectId}/apps/${appId}/replace`, data);
export const deleteProjectApp = async (projectId: string, appId: string) =>
  coreApi.delete(`/projects/${projectId}/apps/${appId}`);
export const updateProjectAppStatus = async (projectId: string, appId: string, data?: Record<string, unknown>) =>
  coreApi.put(`/projects/${projectId}/apps/${appId}/status`, data);
export const getProjectFlows = async (projectId: string) => {
  const resp = await coreApi.get(`/projects/${projectId}/flows`);
  return (resp?.list || []).map(mapFlow);
};
export const getProjectMembers = async (projectId: number): Promise<ProjectMember[]> => {
  const resp = await coreApi.get(`/projects/${projectId}/members`);
  return (resp?.list || []).map((item: any) => ({
    ...item,
    id: String(item.id || item.memberId || ''),
    memberId: String(item.memberId || item.id || ''),
    userId: String(item.userId || ''),
    updatedAt: item.updatedAt || item.updatedTime,
  }));
};
export const addProjectMember = async (projectId: number | string, data?: Record<string, unknown>) =>
  coreApi.post(`/projects/${projectId}/members`, data);
export const updateProjectMember = async (
  projectId: number | string,
  memberId: number | string,
  data?: Record<string, unknown>
) => coreApi.put(`/projects/${projectId}/members/${memberId}`, data);
export const deleteProjectMember = async (projectId: number | string, memberId: string) =>
  coreApi.delete(`/projects/${projectId}/members/${memberId}`);
export const getProjectRoles = async (projectId: number, _params?: ProjectRoleListParams): Promise<ProjectRole[]> => {
  const resp = await coreApi.get(`/projects/${projectId}/roles`);
  return resp?.list || [];
};
export const getProjectRoleOptions = async (projectId: number): Promise<ProjectRoleOption[]> =>
  (await getProjectRoles(projectId)).map((role) => ({
    label: role.roleName || role.roleKey,
    value: String(role.roleId),
  }));
