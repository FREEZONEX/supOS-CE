import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos/userManage';

const api = new ApiWrapper(baseUrl);

// 获取用户信息
export const getUserManageList = async (data?: Record<string, unknown>) =>
  api.post('/pageList', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

// 获取用户信息 - select使用
export const searchUserManageList = async (data?: Record<string, unknown>) =>
  api.post('/pageList', data).then((data: any) => {
    return data?.map?.((item: any) => ({
      label: item.preferredUsername,
      value: item.id,
    }));
  });

// 更新用户
export const updateUser = async (data?: Record<string, unknown>) => api.put('/updateUser', data);

// 更新手机号
export const updatePhone = async (data?: Record<string, unknown>) => api.put(`/phone`, undefined, { params: data });

// 更新邮箱
export const updateEmail = async (data?: Record<string, unknown>) => api.put(`/email`, undefined, { params: data });

// 删除用户
export const deleteUser = async (id: string) => api.delete(`/deleteById/${id}`);

// 重置密码
export const resetPwd = async (data?: Record<string, unknown>) => api.put('/resetPwd', data);
// 用户重置密码
export const userResetPwd = async (data?: Record<string, unknown>) => api.put('/userResetPwd', data);

// 更新用户tips开关启用状态
export const updateTipsEnable = async (enable: number, data?: Record<string, unknown>) =>
  api.put(`/tipsEnable?tipsEnable=${enable}`, data);

// 创建用户
export const createUser = async (data?: Record<string, unknown>) => api.post('/createUser', data);

// 获取角色列表
export const getRoleList = async () => api.get('/roleList');

// 用户重置密码
export const setHomePageApi = async (data?: Record<string, unknown>) => api.put(`/homePage`, {}, { params: data });
