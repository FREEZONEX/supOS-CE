import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos/uns/dashboard';

const api = new ApiWrapper(baseUrl);

// 获取列表 category: group-分组 file-文件   不传=全部; groupId=分组ID（非必填）; k=搜索关键字;page pageSize
export const getDashboardAndGroupList = async (params?: Record<string, unknown>) =>
  api.get('/getGroupedDashboardList', {
    params,
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

// 新增文件
export const addDashboard = async (data: any) => api.post('', data);

// 编辑文件
export const editDashboard = async (data: any) => api.put('', data);

// 删除dashboard
export const deleteDashboard = async (uid: string) => api.delete(`/${uid}`);

// 置顶
export const markDashboard = async (id: string) => api.post('/mark', { id });
export const unmarkDashboard = async (id: string) => api.delete(`/unmark?id=${id}`);

// 获取dashboard信息
export const getDashboardByUns = async (unsAlias: string, config?: any) =>
  api.get(`/getByUns?unsAlias=${unsAlias}`, config);

// 获取dashboard详情
export const getDashboardDetail = async (id: any) => api.get(`/${id}`);

// 获取dashboard信息
export const bindDashboardForUns = async (data: any) => api.post(`/bindUns`, data);

// 创建dashboard
export const createDashboard = async (alias: string) => await api.post(`/createGrafanaByUns/${alias}`);

// 检查dashboard
export const checkDashboardIsExist = async (params?: Record<string, unknown>) =>
  api.get('/isExist', { params, [CustomAxiosConfigEnum.NoCode]: true });

// 只有dashboard的接口
export const getDashboardList = async (params?: Record<string, unknown>) =>
  api.get('', {
    params,
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  }); // 获取dashboard
