import { authApi, iamApi } from './core-adapter';

const defaultProfileEmail = 'tier0@example.com';
const validEmailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const normalizeProfileEmail = (email?: string, userName?: string) => {
  const value = String(email || '').trim();
  if (validEmailPattern.test(value)) return value;
  return String(userName || '').trim() === 'tier0' ? defaultProfileEmail : value;
};

const unsupported = (name: string): Promise<any> =>
  Promise.reject(new Error(`${name} is not implemented in the current backend`));

const unsupportedHandler =
  (name: string) =>
  (...args: unknown[]): Promise<any> => {
    void args;
    return unsupported(name);
  };

const mapRole = (item: any) => ({
  ...item,
  roleId: String(item.id ?? item.roleId),
  roleName: item.name ?? item.roleName ?? item.code,
  roleCode: item.code ?? item.roleCode,
  defaultHomePage: item.defaultHomePage ?? item.homePage ?? '/home',
  resourceList: item.resourceList || [],
  denyResourceList: item.denyResourceList || [],
});

const mapUser = (item: any) => ({
  ...item,
  id: String(item.userId ?? item.id),
  userId: String(item.userId ?? item.id),
  preferredUsername: item.userName ?? item.preferredUsername,
  firstName: item.nickName ?? item.firstName ?? item.userName,
  phone: item.phone ?? '',
  email: normalizeProfileEmail(item.email, item.userName ?? item.preferredUsername),
  enabled: item.enabled ?? item.status !== 0,
  roleList: item.roleList || [],
});

const normalizeRoleList = (value: any) => {
  const roles = Array.isArray(value) ? value : value ? [value] : [];
  return roles
    .map((item: any) => ({
      roleId: Number(item?.roleId ?? item?.id ?? item?.value),
      roleName: item?.roleName ?? item?.name ?? item?.label,
    }))
    .filter((item: any) => Number.isFinite(item.roleId) && item.roleId > 0);
};

const normalizeUserPayload = (data?: Record<string, any>) => ({
  username: data?.username ?? data?.userName ?? data?.preferredUsername,
  password: data?.password,
  firstName: data?.firstName ?? data?.nickName,
  email: data?.email,
  phone: data?.phone,
  enabled: data?.enabled !== false,
  homePage: data?.homePage,
  roleList: normalizeRoleList(data?.roleList),
});

const normalizeCurrentUserProfilePayload = (data?: Record<string, any>) => ({
  firstName: data?.firstName ?? data?.nickName,
  email: data?.email,
  phone: data?.phone,
});

// 获取用户信息
export const getUserManageList = async (data?: Record<string, unknown>) => {
  const resp = await iamApi.get('/users', { params: data });
  const list = Array.isArray(resp?.list) ? resp.list : [];
  return {
    data: list.map(mapUser),
    total: Number(resp?.total ?? list.length),
    pageNo: Number(data?.pageNo ?? 1),
    pageSize: Number(data?.pageSize ?? Math.max(list.length, 20)),
  };
};

// 获取用户信息 - select使用
export const searchUserManageList = async (data?: Record<string, unknown>) =>
  getUserManageList(data).then((resp: any) =>
    (resp?.data || []).map((item: any) => ({
      label: item.preferredUsername,
      value: item.id,
    }))
  );

// 更新用户
export const updateUser = async (data?: Record<string, any>) => {
  const userId = data?.userId ?? data?.id;
  return mapUser(await iamApi.put(`/users/${userId}`, normalizeUserPayload(data)));
};

// 更新当前登录用户基础信息
export const updateCurrentUserProfile = async (data?: Record<string, any>) =>
  mapUser(await authApi.put('/profile', normalizeCurrentUserProfilePayload(data)));

// 更新手机号
export const updatePhone = async (data?: Record<string, any>) => updateCurrentUserProfile(data);

// 更新邮箱
export const updateEmail = async (data?: Record<string, any>) => updateCurrentUserProfile(data);

// 删除用户
export const deleteUser = async (id: string | number) => iamApi.delete(`/users/${id}`);

// 重置密码
export const resetPwd = async (data?: Record<string, any>) => {
  const userId = data?.userId ?? data?.id;
  return iamApi.post(`/users/${userId}/password/reset`, { password: data?.password });
};

// 用户重置密码
export const userResetPwd = async (data?: Record<string, any>) =>
  authApi.put('/password', {
    oldPassword: data?.oldPassword ?? data?.password,
    password: data?.password,
    newPassword: data?.newPassword,
  });

// 更新用户tips开关启用状态
export const updateTipsEnable = unsupportedHandler('updateTipsEnable');

// 创建用户
export const createUser = async (data?: Record<string, any>) =>
  mapUser(await iamApi.post('/users', normalizeUserPayload(data)));

// 获取角色列表
export const getRoleList = async () => {
  const resp = await iamApi.get('/roles');
  const list = Array.isArray(resp?.list) ? resp.list : [];
  return list.map(mapRole);
};

// 更新个人首页
export const setHomePageApi = async (_data?: Record<string, unknown>) => authApi.put('/config', _data || {});
